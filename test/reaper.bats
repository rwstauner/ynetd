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

  # kill and restart multiple times.
  for i in 1 2 3; {
    is "`ysend hello`" = "reap$YTAG"
    running ytester
    kill `ypidof ytester`

    # killed.
    ! running ytester
    no_zombies
  }
}

@test "ynetd remains responsive after SIGCHLD" {
  if ! kill -l | grep -q CHLD; then
    skip "SIGCHLD not supported on this platform"
  fi

  ytester -loop -serve "reap$YTAG"
  running ynetd
  ! running ytester

  is "`ysend hello`" = "reap$YTAG"
  running ytester
  kill -s CHLD $YPID

  # The signal was fake, the process should still be running.
  running ytester
  # kill it, try again

  is "`ysend hello`" = "reap$YTAG"
}
