#!/usr/bin/env bats

load helpers

@test "command starts when port is used" {
  ynetbash 'while echo "wave$YTAG" | listen; do sleep 1; done'
  running ynetd
  ! running ynetbash
  ysend hello | grep -qFx "wave$YTAG"
  running ynetbash
  close
  ! running ynetd
  ! running ynetbash
}
