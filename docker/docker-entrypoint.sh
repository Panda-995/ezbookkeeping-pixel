#!/bin/sh

set -eu

conf_path_param=""
data_directory="${EZBOOKKEEPING_DATA_DIRECTORY:-/ezbookkeeping/data}"
storage_directory="${EZBOOKKEEPING_STORAGE_DIRECTORY:-/ezbookkeeping/storage}"
log_directory="${EZBOOKKEEPING_LOG_DIRECTORY:-/ezbookkeeping/log}"

print_directory_diagnostics() {
  directory_path="$1"

  echo "Container identity: uid=$(id -u) gid=$(id -g)" >&2
  echo "Directory details:" >&2
  ls -ldn "${directory_path}" >&2 2>/dev/null || true
  echo "Host bind mounts must be writable by that UID/GID." >&2
  echo "For the default non-root image user, run this once on the host:" >&2
  echo "  sudo chown -R 1000:1000 ./data ./storage" >&2
  echo "  sudo chmod -R u+rwX ./data ./storage" >&2
}

ensure_writable_directory() {
  directory_path="$1"
  directory_label="$2"

  if ! mkdir -p "${directory_path}"; then
    echo "ezBookkeeping cannot create the ${directory_label} directory: ${directory_path}" >&2
    print_directory_diagnostics "${directory_path}"
    exit 1
  fi

  probe_path="${directory_path}/.ezbookkeeping-write-check-$$"

  if ! touch "${probe_path}"; then
    echo "ezBookkeeping cannot write to the ${directory_label} directory: ${directory_path}" >&2
    print_directory_diagnostics "${directory_path}"
    exit 1
  fi

  rm -f "${probe_path}"
}

ensure_writable_directory "${data_directory}" "database"
ensure_writable_directory "${storage_directory}" "storage"
ensure_writable_directory "${log_directory}" "log"

if [ "${EBK_CONF_PATH:-}" != "" ]; then
  conf_path_param="--conf-path=${EBK_CONF_PATH}"
fi

if [ $# -gt 0 ]; then
  exec "$@"
elif [ "${conf_path_param}" != "" ]; then
  exec /ezbookkeeping/ezbookkeeping server run "${conf_path_param}"
else
  exec /ezbookkeeping/ezbookkeeping server run
fi
