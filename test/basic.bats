#!/usr/bin/env bats

load helpers

@test "command starts when port is used" {
  ynetbash 'while true; do serve "wave$YTAG"; done'
  running ynetd
  ! running ytester
  ysend hello | grep -qFx "wave$YTAG"
  running ytester
  close
  ! running ynetd
  ! running ytester
}
