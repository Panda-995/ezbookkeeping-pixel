#!/bin/sh

set -eu

entrypoint_path="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)/docker-entrypoint.sh"
temporary_root="$(mktemp -d)"

cleanup() {
  rm -rf "${temporary_root}"
}

trap cleanup EXIT INT TERM

# A bind mount that Docker created on the host is commonly root:root and 0755.
# Startup must repair it, then execute the application as the unprivileged user.
chmod 0755 "${temporary_root}"
mkdir -p "${temporary_root}/root-owned/data" \
  "${temporary_root}/root-owned/storage" \
  "${temporary_root}/root-owned/log"
chown -R 0:0 "${temporary_root}/root-owned"
chmod -R 0755 "${temporary_root}/root-owned"
printf '%s\n' 'existing database content' > "${temporary_root}/root-owned/data/ezbookkeeping.db"
chmod 0600 "${temporary_root}/root-owned/data/ezbookkeeping.db"

EZBOOKKEEPING_DATA_DIRECTORY="${temporary_root}/root-owned/data" \
EZBOOKKEEPING_STORAGE_DIRECTORY="${temporary_root}/root-owned/storage" \
EZBOOKKEEPING_LOG_DIRECTORY="${temporary_root}/root-owned/log" \
ENTRYPOINT_UNDER_TEST="${entrypoint_path}" \
  sh -c '
    set -- sh -c '\''
      id -u > "${EZBOOKKEEPING_DATA_DIRECTORY}/runtime-uid"
      id -g > "${EZBOOKKEEPING_DATA_DIRECTORY}/runtime-gid"
    '\''
    . "${ENTRYPOINT_UNDER_TEST}"
  '

test "$(cat "${temporary_root}/root-owned/data/runtime-uid")" = "1000"
test "$(cat "${temporary_root}/root-owned/data/runtime-gid")" = "1000"
test "$(stat -c '%u:%g' "${temporary_root}/root-owned/data")" = "1000:1000"
test "$(stat -c '%u:%g' "${temporary_root}/root-owned/storage")" = "1000:1000"
test "$(stat -c '%u:%g' "${temporary_root}/root-owned/log")" = "1000:1000"
test "$(stat -c '%u:%g' "${temporary_root}/root-owned/data/ezbookkeeping.db")" = "1000:1000"
grep -F 'existing database content' "${temporary_root}/root-owned/data/ezbookkeeping.db" >/dev/null

# An explicit non-root override cannot repair a root-owned bind mount and should
# still fail with actionable diagnostics instead of starting as root.
mkdir -p "${temporary_root}/non-root/data" \
  "${temporary_root}/non-root/storage" \
  "${temporary_root}/non-root/log"
chown -R 0:0 "${temporary_root}/non-root"
chmod -R 0755 "${temporary_root}/non-root"

set +e
failure_output="$({
  su-exec 1000:1000 env \
    EZBOOKKEEPING_DATA_DIRECTORY="${temporary_root}/non-root/data" \
    EZBOOKKEEPING_STORAGE_DIRECTORY="${temporary_root}/non-root/storage" \
    EZBOOKKEEPING_LOG_DIRECTORY="${temporary_root}/non-root/log" \
    ENTRYPOINT_UNDER_TEST="${entrypoint_path}" \
      sh -c 'set -- /bin/true; . "${ENTRYPOINT_UNDER_TEST}"'
} 2>&1)"
failure_status=$?
set -e

if [ "${failure_status}" -eq 0 ]; then
  echo "expected an explicit non-root process with an unwritable bind mount to stop startup" >&2
  exit 1
fi

echo "${failure_output}" | grep -F "cannot write to the database directory" >/dev/null
echo "${failure_output}" | grep -F "Container identity: uid=1000 gid=1000" >/dev/null
echo "${failure_output}" | grep -F "Host bind mounts must be writable" >/dev/null

echo "docker entrypoint tests passed"
