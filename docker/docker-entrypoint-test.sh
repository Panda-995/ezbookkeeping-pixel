#!/bin/sh

set -eu

entrypoint_path="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)/docker-entrypoint.sh"
temporary_root="$(mktemp -d)"

cleanup() {
  rm -rf "${temporary_root}"
}

trap cleanup EXIT INT TERM

mkdir -p "${temporary_root}/data" "${temporary_root}/storage" "${temporary_root}/log"

EZBOOKKEEPING_DATA_DIRECTORY="${temporary_root}/data" \
EZBOOKKEEPING_STORAGE_DIRECTORY="${temporary_root}/storage" \
EZBOOKKEEPING_LOG_DIRECTORY="${temporary_root}/log" \
ENTRYPOINT_UNDER_TEST="${entrypoint_path}" \
  sh -c 'set -- /bin/true; . "${ENTRYPOINT_UNDER_TEST}"'

set +e
failure_output="$({
  EZBOOKKEEPING_DATA_DIRECTORY="${temporary_root}/data" \
  EZBOOKKEEPING_STORAGE_DIRECTORY="${temporary_root}/storage" \
  EZBOOKKEEPING_LOG_DIRECTORY="${temporary_root}/log" \
  ENTRYPOINT_UNDER_TEST="${entrypoint_path}" \
    sh -c '
      touch() { return 1; }
      set -- /bin/true
      . "${ENTRYPOINT_UNDER_TEST}"
    '
} 2>&1)"
failure_status=$?
set -e

if [ "${failure_status}" -eq 0 ]; then
  echo "expected an unwritable bind mount to stop startup" >&2
  exit 1
fi

echo "${failure_output}" | grep -F "cannot write to the database directory" >/dev/null
echo "${failure_output}" | grep -F "Container identity: uid=" >/dev/null
echo "${failure_output}" | grep -F "Host bind mounts must be writable by that UID/GID" >/dev/null

echo "docker entrypoint tests passed"
