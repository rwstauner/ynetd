LISTEN_PORT=$((63000 + (RANDOM % 1000)))
PROXY_PORT=$((64000 + (RANDOM % 1000)))

YTAG=
YPID=

logdir=tmp
mkdir -p "$logdir"

debug () {
  echo " # $*" >&3
}

knock () {
  nc -z localhost "$LISTEN_PORT"
}

running () {
  # Use subshell to help command terminate.
  (ps -o args | grep -E "^$1$YTAG")
}

ylog () {
  cat $YLOG
}

ynetd () {
  YTAG=":$((RANDOM))"
  YLOG="$logdir/test$YTAG.log"
  # Use exec to separate from bats and set $0.
  (YTAG="$YTAG" exec -a "ynetd$YTAG" "${YNETD:-ynetd}" "$@" &> "$YLOG") &
  YPID=$!
  sleep 1
}

ynetbash () {
  # Last arg is script.
  args=()
  while [[ $# -gt 1 ]]; do
    args+=("$1")
    shift
  done

  ynetd -listen "localhost:$LISTEN_PORT" -proxy "localhost:$PROXY_PORT" "${args[@]}" \
    bash -c 'exec -a ynetbash$YTAG bash -c "$@"' -- \
      "listen () { nc -l -p $PROXY_PORT localhost; }; $1"
}

ysend () {
  echo "$*" | nc localhost $LISTEN_PORT
}

close () {
  if [[ -n "$YPID" ]]; then
    # Don't count these exit statuses as errors.
    kill $YPID || :
    wait $YPID || :
  fi
  YPID=
}

teardown () {
  close
  if [[ -n "$YLOG" ]]; then
    rm -f "$YLOG"
  fi
}
