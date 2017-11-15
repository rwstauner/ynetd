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
  ${YTESTER:-build/ytester} -knock -port "$LISTEN_PORT"
}

no_zombies () {
  # Ignore "bash" as some bats helpers can be temporarily zombied.
  ! (ps -o state,args | grep -v bash | grep -E '^Z|defunct')
}

running () {
  # Use subshell to help command terminate.
  (ps -o args | grep -E "^$1$YTAG")
}

signal () {
  kill `ps -o pid,args | awk -v CMD="$1$YTAG" '$2 ~ CMD { print $1 }'`
}

ylog () {
  cat $YLOG
}

ynetd () {
  YLOG="$logdir/test$YTAG.log"
  # Use exec to separate from bats and set $0.
  (YTAG="$YTAG" exec -a "ynetd$YTAG" "${YNETD:-ynetd}" "$@" &> "$YLOG") &
  YPID=$!
}

ytester () {
  ynetd -listen "localhost:$LISTEN_PORT" -proxy "localhost:$PROXY_PORT" "${YARGS[@]}" \
    bash -c 'exec -a ytester$YTAG "$@"' -- \
      "${YTESTER:-build/ytester}" -port "$PROXY_PORT" "$@"
}

ysend () {
  ${YTESTER:-build/ytester} -send "$*" -port "$LISTEN_PORT"
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
