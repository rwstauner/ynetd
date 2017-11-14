#!/usr/bin/env bats

load helpers

@test "command restarts after exiting" {
  ynetbash 'serve foo; sleep 2'
  running ynetd
  ! running ytester
  knock
  running ytester
  sleep 2
  ! running ytester
  no_zombies
  knock
  running ytester
}

@test "command restarts when killed" {
  ynetbash 'while true; do serve "reap$YTAG"; done'
  running ynetd
  ! running ytester

  ysend hello | grep -qFx "reap$YTAG"
  running ytester
  signal ytester

  # killed.
  ! running ytester
  no_zombies

  ysend hello | grep -qFx "reap$YTAG"
  running ytester
}
