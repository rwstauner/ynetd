#!/usr/bin/env bats

load helpers

ystate () {
  ps -p "$YPID" -o state= | cut -c 1
}

@test "exits without ever having done anything" {
  ytester -loop -serve "wave$YTAG"
  running ynetd
  ! running ytester
  kill -s INT $YPID
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
