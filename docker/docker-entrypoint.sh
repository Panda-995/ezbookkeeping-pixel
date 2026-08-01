#!/bin/sh

set -e;

conf_path_param="";

ensure_writable_directory() {
  directory_path="$1"
  directory_label="$2"

  if ! mkdir -p "${directory_path}"; then
    echo "ezBookkeeping cannot create the ${directory_label} directory: ${directory_path}" >&2
    exit 1
  fi

  probe_path="${directory_path}/.ezbookkeeping-write-check-$$"

  if ! : > "${probe_path}"; then
    echo "ezBookkeeping cannot write to the ${directory_label} directory: ${directory_path}" >&2
    echo "Check the bind-mount path and permissions before restarting the container." >&2
    exit 1
  fi

  rm -f "${probe_path}"
}

ensure_writable_directory "/ezbookkeeping/data" "database"
ensure_writable_directory "/ezbookkeeping/storage" "storage"
ensure_writable_directory "/ezbookkeeping/log" "log"

if [ "${EBK_CONF_PATH}" != "" ]; then
  conf_path_param="--conf-path=${EBK_CONF_PATH}";
fi

if [ $# -gt 0 ]; then
    exec "$@"
else
    exec /ezbookkeeping/ezbookkeeping server run ${conf_path_param};
fi
