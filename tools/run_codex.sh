#!/usr/bin/env bash
# Launch a Codex task with a deadline, live status, and an audit copy.
set -u

HOME_DIR=${HOME:-}
if [ -z "$HOME_DIR" ]; then
  echo "run_codex.sh: HOME is required" >&2
  exit 2
fi

LOG_ROOT="$HOME_DIR/siduri-wt"
AUDIT_ROOT=${RUN_CODEX_AUDIT_ROOT:-$LOG_ROOT/audit/runs}

label=""
wt=""
log=""
status=""
pid_file=""
requested_label=""
deadline=""
add_dirs=()
add_dir_args=()

usage() {
  echo "usage: run_codex.sh --deadline <s> [--label <l>] <worktree> <task-file>..." >&2
}

timestamp() {
  date -u '+%Y-%m-%dT%H:%M:%S.%3NZ'
}

status_line() {
  printf '%s\n' "$1" | tee -a "$status"
}

control_line() {
  printf '%s\n' "$1" | tee -a "$log" "$status"
}

audit_lane() {
  [ -n "${label:-}" ] || return 0
  [ -n "${log:-}" ] && [ -f "$log" ] || return 0

  local audit_dir="$AUDIT_ROOT/$label"
  mkdir -p "$audit_dir" 2>/dev/null || return 0
  cp "$log" "$audit_dir/$label.log" 2>/dev/null || true
  if [ -f "$status" ]; then
    cp "$status" "$audit_dir/$label.status" 2>/dev/null || true
  fi
  if [ -f "$pid_file" ]; then
    cp "$pid_file" "$audit_dir/$label.pid" 2>/dev/null || true
  fi
  if [ -f "$status" ]; then
    awk '/ after / { verdict = $0 } END { if (verdict != "") print verdict }' \
      "$status" > "$audit_dir/$label.verdict"
  fi
}

finish_lane() {
  local rc="$1" started="$2" secs="$3" elapsed head dirty verdict
  elapsed=$(( $(date -u +%s) - started ))
  head=$(git -C "$wt" log --oneline -1 2>/dev/null || true)
  dirty=$(git -C "$wt" status --porcelain 2>/dev/null | wc -l)

  control_line "[$(timestamp)] $label exit $rc"
  if [ "$rc" -eq 124 ] || [ "$rc" -eq 137 ]; then
    verdict="TIMED OUT at ${secs}s -- killed, work left uncommitted in the worktree"
  elif [ "$rc" -ne 0 ]; then
    verdict="FAILED rc=$rc"
  else
    verdict="exited clean"
  fi
  control_line "[$(timestamp)] $label $verdict after ${elapsed}s | head: $head | dirty: $dirty"
  audit_lane
  return "$rc"
}

resolve_git_dir() {
  local worktree="$1" dot_git="$1/.git" git_dir
  if [ -f "$dot_git" ]; then
    git_dir=$(sed -n 's/^gitdir:[[:space:]]*//p' "$dot_git")
    [ -n "$git_dir" ] || return 2
    case "$git_dir" in
      /*) ;;
      *) git_dir=$(builtin cd -- "$(dirname "$dot_git")" && builtin cd -- "$git_dir" && pwd -P) || return 2 ;;
    esac
  elif [ -d "$dot_git" ]; then
    git_dir=$(builtin cd -- "$dot_git" && pwd -P) || return 2
  else
    return 2
  fi
  printf '%s\n' "$git_dir"
}

prepare_add_dirs() {
  local worktree="$1" npm_cache git_dir
  git_dir=$(resolve_git_dir "$worktree") || {
    control_line "[$(timestamp)] $label: cannot resolve $worktree/.git"
    return 2
  }

  add_dirs=(
    "$HOME_DIR/.cache/go-build"
    "$HOME_DIR/go/pkg/mod"
  )
  npm_cache=""
  if command -v npm >/dev/null 2>&1; then
    npm_cache=$(npm config get cache 2>/dev/null || true)
    case "$npm_cache" in
      ""|*"$'\n'"*) npm_cache="" ;;
    esac
  fi
  [ -n "$npm_cache" ] && add_dirs+=("$npm_cache")
  add_dirs+=("$git_dir")

  add_dir_args=()
  local add_dir
  for add_dir in "${add_dirs[@]}"; do
    add_dir_args+=(--add-dir "$add_dir")
  done
}

run_lane() {
  label="$1"
  wt="$2"
  local secs="$3" task="$4" task_dir
  local started agent last_size last_change reported size now quiet rc

  log="$LOG_ROOT/$label.log"
  status="$LOG_ROOT/$label.status"
  pid_file="$LOG_ROOT/$label.pid"
  trap audit_lane EXIT INT TERM

  started=$(date -u +%s)
  control_line "[$(timestamp)] $label started, bound ${secs}s, worktree $wt"
  if [ ! -r "$task" ]; then
    control_line "[$(timestamp)] $label task file unavailable: $task"
    finish_lane 2 "$started" "$secs"
    return 2
  fi
  task_dir=$(builtin cd -- "$(dirname "$task")" && pwd -P) || {
    control_line "[$(timestamp)] $label task directory unavailable: $task"
    finish_lane 2 "$started" "$secs"
    return 2
  }
  task="$task_dir/$(basename "$task")"
  cd "$wt" || {
    control_line "[$(timestamp)] $label: no such worktree"
    finish_lane 2 "$started" "$secs"
    return 2
  }
  prepare_add_dirs "$wt" || {
    finish_lane 2 "$started" "$secs"
    return 2
  }

  KILL_AFTER_S=${KILL_AFTER_S:-30}
  timeout --kill-after="${KILL_AFTER_S}s" "$secs" \
    codex exec -s danger-full-access -C "$wt" --color never \
      "${add_dir_args[@]}" - >>"$log" 2>&1 <"$task" &
  agent=$!
  printf '%s\n' "$agent" > "$pid_file"

  STALL_S=${STALL_S:-600}
  STALL_KILL_S=${STALL_KILL_S:-1800}
  POLL_S=${POLL_S:-30}
  last_size=0
  last_change=$(date -u +%s)
  reported=0
  while kill -0 "$agent" 2>/dev/null; do
    sleep "$POLL_S"
    size=$(wc -c < "$log" 2>/dev/null || echo 0)
    now=$(date -u +%s)
    if [ "$size" -ne "$last_size" ]; then
      last_size="$size"
      last_change="$now"
      reported=0
    elif [ $(( now - last_change )) -ge "$STALL_S" ]; then
      quiet=$(( now - last_change ))
      if [ "$quiet" -ge "$STALL_KILL_S" ]; then
        control_line "[$(timestamp)] $label STALL-KILL: silent ${quiet}s >= ${STALL_KILL_S}s, killing"
        kill -TERM "$agent" 2>/dev/null || true
        sleep "${STALL_TERM_WAIT_S:-10}"
        kill -KILL "$agent" 2>/dev/null || true
        break
      elif [ "$reported" -eq 0 ]; then
        control_line "[$(timestamp)] $label STALLED (stalled): log unchanged for ${quiet}s (alive, $size bytes)"
        last_size=$(wc -c < "$log" 2>/dev/null || echo 0)
        reported=1
      fi
    fi
  done
  wait "$agent"
  rc=$?
  finish_lane "$rc" "$started" "$secs"
  return "$rc"
}

prepare_lane() {
  label="$1"
  log="$LOG_ROOT/$label.log"
  status="$LOG_ROOT/$label.status"
  pid_file="$LOG_ROOT/$label.pid"
  : > "$log"
  : > "$status"
  rm -f "$pid_file"
  status_line "[$(timestamp)] $label queued"
}

status_has_exit() {
  awk '$0 ~ / exit [0-9]+$/ { found = 1 } END { exit !found }' "$1"
}

run_batch() {
  local secs="$1" worktree="$2"
  shift 2
  local total="$#" slots nproc_count next active i lane_rc highest_rc
  local candidate task_base base_label duplicate j
  local -a task_files labels controllers

  task_files=("$@")
  nproc_count=$(nproc 2>/dev/null || echo 2)
  slots=$(( nproc_count - 1 ))
  [ "$slots" -lt 1 ] && slots=1
  [ "$slots" -gt "$total" ] && slots="$total"

  base_label="$requested_label"
  for i in "${!task_files[@]}"; do
    if [ -n "$base_label" ]; then
      if [ "$total" -eq 1 ]; then
        candidate="$base_label"
      else
        candidate="${base_label}-$((i + 1))"
      fi
    else
      task_base=${task_files[$i]##*/}
      candidate=${task_base%.*}
      [ -n "$candidate" ] || candidate="lane-$((i + 1))"
    fi
    while :; do
      duplicate=0
      for j in "${!labels[@]}"; do
        [ "${labels[$j]}" = "$candidate" ] && duplicate=1
      done
      [ "$duplicate" -eq 0 ] && break
      candidate="$candidate-$((i + 1))"
    done
    case "$candidate" in
      ""|*[!A-Za-z0-9._-]*) usage; return 2 ;;
    esac
    labels[i]="$candidate"
    label="$candidate"
    prepare_lane "$label"
  done

  next=0
  active=0
  highest_rc=0
  while [ "$next" -lt "$total" ] && [ "$active" -lt "$slots" ]; do
    label="${labels[$next]}"
    log="$LOG_ROOT/$label.log"
    status="$LOG_ROOT/$label.status"
    pid_file="$LOG_ROOT/$label.pid"
    run_lane "$label" "$worktree" "$secs" "${task_files[$next]}" &
    controllers[next]=$!
    next=$((next + 1))
    active=$((active + 1))
  done

  while [ "$active" -gt 0 ]; do
    i=0
    while [ "$i" -lt "$total" ]; do
      if [ -n "${controllers[$i]-}" ]; then
        label="${labels[$i]}"
        status="$LOG_ROOT/$label.status"
        if status_has_exit "$status"; then
          wait "${controllers[$i]}"
          lane_rc=$?
          [ "$lane_rc" -gt "$highest_rc" ] && highest_rc="$lane_rc"
          unset 'controllers[i]'
          active=$((active - 1))
          if [ "$next" -lt "$total" ]; then
            label="${labels[$next]}"
            log="$LOG_ROOT/$label.log"
            status="$LOG_ROOT/$label.status"
            pid_file="$LOG_ROOT/$label.pid"
            run_lane "$label" "$worktree" "$secs" "${task_files[$next]}" &
            controllers[next]=$!
            next=$((next + 1))
            active=$((active + 1))
          fi
        fi
      fi
      i=$((i + 1))
    done
    [ "$active" -gt 0 ] && sleep "${QUEUE_POLL_S:-0.1}"
  done
  return "$highest_rc"
}

while [ "$#" -gt 0 ]; do
  case "$1" in
    --deadline)
      [ "$#" -ge 2 ] || { usage; exit 2; }
      deadline="$2"
      shift 2
      ;;
    --label)
      [ "$#" -ge 2 ] || { usage; exit 2; }
      requested_label="$2"
      shift 2
      ;;
    --*)
      usage
      exit 2
      ;;
    *)
      break
      ;;
  esac
done

if [ -z "$deadline" ] || ! [[ "$deadline" =~ ^[0-9]+([.][0-9]+)?$ ]] || [ "$#" -lt 2 ]; then
  usage
  exit 2
fi

worktree=$(builtin cd -- "$1" && pwd -P) || {
  usage
  exit 2
}
shift

mkdir -p "$LOG_ROOT" || exit 2
run_batch "$deadline" "$worktree" "$@"
exit $?
