#!/usr/bin/env bats

load helpers

@test "stop command after timeout" {
  YARGS=(-stop-after 2s)
  ytester -loop -serve "stop$YTAG"

  for i in 1 2 3; {
    knock
    running ytester

    ylog | grep 'starting: ' | lines $i

    # Use it.
    for j in 1 2 3; {
      ysend hello | grep -qFx "stop$YTAG"
      sleep 1
    }

    # Let it expire.
    sleep 3

    ! running ytester

    ylog | grep 'stopping ' | lines $i
  }
}
