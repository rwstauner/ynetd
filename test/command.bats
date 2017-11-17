#!/usr/bin/env bats

load helpers

@test "only runs once" {
  ytester -loop -serve "wave$YTAG"
  running ynetd
  ! running ytester

  # Make several simultaneous requests.
  for i in 1 2 3; { knock & }
  # Wait for service to be done (better than a sleep).
  ysend hello

  running ytester

  close
  ylog -y | grep Starting: | lines 1
}
