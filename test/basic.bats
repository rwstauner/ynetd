#!/usr/bin/env bats

load helpers

@test "command starts when port is used" {
  ynetbash 'while true; do serve "wave$YTAG"; done'
  running ynetd
  ! running ynetbash
  ysend hello | grep -qFx "wave$YTAG"
  running ynetbash
  close
  ! running ynetd
  ! running ynetbash
}
