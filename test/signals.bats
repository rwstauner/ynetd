#!/usr/bin/env bats

load helpers

ystate () {
  ps -p "$YPID" -o state= | cut -c 1
}

@test "exits without ever having done anything" {
  ytester -loop -serve "wave$YTAG"
  running ynetd
  ! running ytester
  kill $YPID
  sleep 1 # With -race this takes slightly longer.
  ! running ynetd
}

@test "can be suspended and resumed" {
  ytester -loop -serve "wave$YTAG"
  ysend hello
  is `ystate` = S

  kill -s STOP $YPID
  is `ystate` = T

  kill -s CONT $YPID
  is `ystate` = S
}

@test "ignores INT when already ignored" {
  "$YAS" "ynetd$YTAG" "$YNETD" -proxy "localhost:$LISTEN_PORT localhost:$PROXY_PORT" &
  YPID=$!
  running ynetd

  kill -s INT $YPID
  sleep 1
  running ynetd

  kill -s TERM $YPID
  sleep 1
  ! running ynetd
}

@test "respects INT in the foreground" {
  # For the default bash (3.2) on MacOS (High Sierra)
  # this is true when running more than one bats file (bats-exec-suite).
  if "$YTESTER" -int-ignored 2>&1 | grep -q "int ignored: true"; then
    skip "SIGINT is already ignored"
  fi

  (pid=; while [[ -z "$pid" ]]; do sleep 1; pid=`ypidof ynetd`; done; kill -s INT "$pid") & # bg
  "$YAS" "ynetd$YTAG" "$YNETD" -proxy "localhost:$LISTEN_PORT localhost:$PROXY_PORT" # fg
  # First (bg) proc should kill the second (fg) one.
  ! running ynetd
}

@test "signal process group" {
  ynetd -proxy "localhost:$LISTEN_PORT localhost:$PROXY_PORT" -stop-after 2s \
    "$YAS" "ytester$YTAG" \
      /bin/bash -c "($YAS ychild$YTAG $YTESTER -port $PROXY_PORT -loop -serve main)"

  knock
  running ytester
  running ychild
  close
  ! running tester
  ! running ychild
}
