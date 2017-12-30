#!/usr/bin/env bats

load helpers

@test "command starts when port is used" {
  ytester -loop -serve "wave$YTAG"
  running ynetd
  ! running ytester
  is "`ysend hello`" = "wave$YTAG"
  running ytester
  close
  ! running ynetd
  ! running ytester
}
