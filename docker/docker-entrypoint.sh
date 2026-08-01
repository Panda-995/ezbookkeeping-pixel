#!/bin/sh

set -eu

conf_path_param=""
data_directory="${EZBOOKKEEPING_DATA_DIRECTORY:-/ezbookkeeping/data}"
storage_directory="${EZBOOKKEEPING_STORAGE_DIRECTORY:-/ezbookkeeping/storage}"
log_directory="${EZBOOKKEEPING_LOG_DIRECTORY:-/ezbookkeeping/log}"
runtime_uid="${EZBOOKKEEPING_RUN_UID:-1000}"
runtime_gid="${EZBOOKKEEPING_RUN_GID:-1000}"
container_uid="$(id -u)"
container_gid="$(id -g)"

validate_runtime_identity() {
  identity_name="$1"
  identity_value="$2"

  case "${identity_value}" in
    ''|*[!0-9]*)
      echo "EZBOOKKEEPING_RUN_${identity_name} must be a positive numeric ID, got: ${identity_value}" >&2
      exit 1
      ;;
    0)
      echo "EZBOOKKEEPING_RUN_${identity_name} must not be 0; the application is required to run as non-root." >&2
      exit 1
      ;;
  esac
}

validate_runtime_identity "UID" "${runtime_uid}"
validate_runtime_identity "GID" "${runtime_gid}"

print_directory_diagnostics() {
  directory_path="$1"

  echo "Container identity: uid=${container_uid} gid=${container_gid}" >&2
  echo "Application identity: uid=${runtime_uid} gid=${runtime_gid}" >&2
  echo "Directory details:" >&2
  ls -ldn "${directory_path}" >&2 2>/dev/null || true
  echo "Host bind mounts must be writable by the application identity." >&2

  if [ "${container_uid}" = "0" ]; then
    echo "Automatic ownership repair failed. This usually means the mount is a NAS/NFS/CIFS share that rejects chown." >&2
  else
    echo "The container was explicitly forced to run as a non-root user, so it cannot repair bind-mount ownership." >&2
    echo "Remove the Compose 'user:' override to enable automatic startup repair." >&2
  fi

  echo "Alternatively, grant uid:gid ${runtime_uid}:${runtime_gid} read/write access on the host or NAS." >&2
}

run_as_application() {
  if [ "${container_uid}" = "0" ]; then
    su-exec "${runtime_uid}:${runtime_gid}" "$@"
  else
    "$@"
  fi
}

probe_writable_directory() {
  directory_path="$1"
  probe_path="${directory_path}/.ezbookkeeping-write-check-$$"

  if ! run_as_application touch "${probe_path}" 2>/dev/null; then
    return 1
  fi

  if ! run_as_application rm -f "${probe_path}" 2>/dev/null; then
    return 1
  fi

  return 0
}

ensure_writable_directory() {
  directory_path="$1"
  directory_label="$2"

  if ! mkdir -p "${directory_path}"; then
    echo "ezBookkeeping cannot create the ${directory_label} directory: ${directory_path}" >&2
    print_directory_diagnostics "${directory_path}"
    exit 1
  fi

  if probe_writable_directory "${directory_path}"; then
    return 0
  fi

  if [ "${container_uid}" = "0" ]; then
    echo "Initializing ${directory_label} directory ownership as ${runtime_uid}:${runtime_gid}: ${directory_path}" >&2

    if ! chown -R "${runtime_uid}:${runtime_gid}" "${directory_path}"; then
      echo "ezBookkeeping cannot repair the ${directory_label} directory ownership: ${directory_path}" >&2
      print_directory_diagnostics "${directory_path}"
      exit 1
    fi

    if probe_writable_directory "${directory_path}"; then
      return 0
    fi
  fi

  echo "ezBookkeeping cannot write to the ${directory_label} directory: ${directory_path}" >&2
  print_directory_diagnostics "${directory_path}"
  exit 1
}

if [ "${container_uid}" = "0" ] && ! command -v su-exec >/dev/null 2>&1; then
  echo "ezBookkeeping cannot drop root privileges because su-exec is unavailable." >&2
  exit 1
fi

ensure_writable_directory "${data_directory}" "database"
ensure_writable_directory "${storage_directory}" "storage"
ensure_writable_directory "${log_directory}" "log"

if [ "${EBK_CONF_PATH:-}" != "" ]; then
  conf_path_param="--conf-path=${EBK_CONF_PATH}"
fi

if [ $# -gt 0 ]; then
  set -- "$@"
elif [ "${conf_path_param}" != "" ]; then
  set -- /ezbookkeeping/ezbookkeeping server run "${conf_path_param}"
else
  set -- /ezbookkeeping/ezbookkeeping server run
fi

if [ "${container_uid}" = "0" ]; then
  exec su-exec "${runtime_uid}:${runtime_gid}" "$@"
else
  exec "$@"
fi
