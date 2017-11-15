#!/usr/bin/env bats

load helpers

@test "command restarts after exiting" {
  ytester -before 2s
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
  ytester -loop -serve "reap$YTAG"
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
