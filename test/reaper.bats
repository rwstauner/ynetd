#!/usr/bin/env bats

load helpers

@test "command restarts after exiting" {
  ynetbash 'serve foo; sleep 2'
  running ynetd
  ! running ynetbash
  knock
  running ynetbash
  sleep 2
  ! running ynetbash
  no_zombies
  knock
  running ynetbash
}

@test "command restarts when killed" {
  ynetbash 'while true; do serve "reap$YTAG"; done'
  running ynetd
  ! running ynetbash

  ysend hello | grep -qFx "reap$YTAG"
  running ynetbash
  signal ynetbash

  # killed.
  ! running ynetbash
  no_zombies

  ysend hello | grep -qFx "reap$YTAG"
  running ynetbash
}
