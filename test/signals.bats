#!/usr/bin/env bats

load helpers

ystate () {
  ps -p "$YPID" -o state= | cut -c 1
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
