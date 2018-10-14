LISTEN_PORT=$((63000 + (RANDOM % 1000)))
PROXY_PORT=$((64000 + (RANDOM % 1000)))

YTAG=
YPID=

YNETD="${YNETD:-build/ynetd}"
YTESTER="${YTESTER:-build/ytester}"
YAS="${BASH_SOURCE[0]%/*}/as"

logdir=tmp
mkdir -p "$logdir"

debug () {
  echo " # $*" >&3
}

dgrep () {
  lines=`cat`
  echo "$lines" | grep "$@" || { echo "$lines" >&2; false; }
}

is () {
  echo "[ $* ]" >&2 # for debugging
  [ "$@" ]
}

knock () {
  "$YTESTER" -knock -port "$LISTEN_PORT"
}

lines () {
  is "$1" -eq "$(cat | wc -l)"
}

no_zombies () {
  # Ignore "bash" as some bats helpers can be temporarily zombied.
  ! (ps -o state,args | grep -vE 'bash|grep' | grep -E '^Z|defunct')
}

running () {
  # Use subshell to help command terminate.
  (ps -o args | dgrep -E "^$1$YTAG")
}

ypidof () {
  ps -o pid,args | awk -v CMD="$1$YTAG" '$2 ~ CMD { print $1 }'
}

ylog () {
  cmd=(cat)
  if [[ $1 == "-y" ]]; then
    cmd=(grep -E ^ynetd)
  fi
  "${cmd[@]}" $YLOG
}

ynetd () {
  YLOG="$logdir/test$YTAG.log"
  # Use exec to separate from bats and set $0.
  (YTAG="$YTAG" exec -a "ynetd$YTAG" "$YNETD" "$@" &> "$YLOG") &
  YPID=$!
  # Wait for it to start.
  for i in 1 2 3 4; {
    if [[ -s "$YLOG" ]]; then break; fi
    sleep 0.25 || :
  }
}

ytester () {
  ynetd -proxy "localhost:$LISTEN_PORT localhost:$PROXY_PORT" "${YARGS[@]}" \
    "$YAS" "ytester$YTAG" \
      "$YTESTER" -port "$PROXY_PORT" "$@"
}

ysend () {
  args=()
  while [[ $# -gt 1 ]]; do
    args+=("$1")
    shift
  done

  "$YTESTER" -port "$LISTEN_PORT" "${args[@]}" -send "$1"
}

close () {
  if [[ -n "$YPID" ]]; then
    # Don't count these exit statuses as errors.
    kill $YPID || :
    wait $YPID || :
  fi
  YPID=
}

setup () {
  YTAG=":$((RANDOM))"
  YARGS=()
}

teardown () {
  close
  if [[ -n "$YLOG" ]]; then
    # Dump to STDERR so that if the test fails we see the output.
    ylog >&2
    rm -f "$YLOG"
  fi
}
