LISTEN_PORT=$((63000 + (RANDOM % 1000)))
PROXY_PORT=$((64000 + (RANDOM % 1000)))

YTAG=
YPID=

debug () {
  echo " # $*" >&3
}

running () {
  # Use subshell to help command terminate.
  (ps -o args | grep -E "^$1$YTAG")
}


ynetd () {
  YTAG=":$((RANDOM))"
  # Use exec to separate from bats and set $0.
  (YTAG="$YTAG" exec -a "ynetd$YTAG" "${YNETD:-ynetd}" "$@") &
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
}
