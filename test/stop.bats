#!/usr/bin/env bats

load helpers

wait-each () {
  # Wait individually to ensure we get a non-zero status if any child returned one.
  while [[ $# -gt 0 ]]; do
    wait "$1"
    shift
  done
}

@test "stop command after timeout" {
  YARGS=(-stop-after 2s)
  ytester -loop -serve "stop$YTAG"

  for i in 1 2 3; {
    knock
    running ytester

    ylog | grep 'starting: ' | lines $i

    # Use it.
    for j in 1 2 3; {
      is "`ysend hello`" = "stop$YTAG"
      sleep 1
    }

    # Let it expire.
    sleep 3

    ! running ytester

    ylog | grep 'stopping: ' | lines $i
  }
}

@test "stop process group" {
  ynetd -proxy "localhost:$LISTEN_PORT localhost:$PROXY_PORT" -stop-after 2s \
    "$YAS" "ytester$YTAG" \
      bash -c "echo something; $YTESTER -port $PROXY_PORT -loop -serve \$(date +%s); exit \$?"

  last='none'
  for i in 1 2 3; {
    knock
    running ytester

    ylog | grep 'starting: ' | lines $i

    current=`ysend hello`
    # Different than last time.
    is "$current" != "$last"

    for j in 1 2; {
      # Still the same for this invocation.
      is "`ysend x`" = "$current"
      sleep 0.25
    }

    # Let it expire.
    sleep 3

    ! running ytester

    ylog | grep 'stopping: ' | lines $i
    last="$current"
  }
}

@test "stop after client disconnect" {
  YARGS=(-stop-after 2s)
  ytester -loop -serve "stop$YTAG"

  knock
  running ytester

  sleep 3

  ! running ytester

  is "`ysend -delay 3s slow`" = "stop$YTAG"

  running ytester

  # Let it expire.
  sleep 3

  ! running ytester

  ylog | grep 'stopping: ' | lines 2
}

@test "stop after multiple clients disconnect" {
  YARGS=(-stop-after 2s)
  ytester -loop -serve "multi$YTAG"

  slowly () {
    is "`ysend -delay ${1}s slow$1`" = "multi$YTAG"
  }

  pids=()
  for i in 1 2 1 2; {
    slowly 4 &
    pids+=($!)
    sleep $i
  }

  running ytester

  wait-each "${pids[@]}"

  # Let it expire.
  sleep 3

  ! running ytester

  ylog | grep 'stopping: ' | lines 1
}

@test "stop after longest client disconnects" {
  YARGS=(-stop-after 2s)
  ytester -loop -serve "long$YTAG"

  slowly () {
    is "`ysend -delay ${1}s slow$1`" = "long$YTAG"
  }

  slowly 9 &
  pids=($!)

  slowly 3 &
  pids+=($!)

  sleep 1

  running ytester

  sleep 3

  running ytester

  wait-each "${pids[@]}"

  running ytester

  sleep 3

  ! running ytester

  ylog | grep 'stopping: ' | lines 1
}

@test "stop and restart correctly" {
  YARGS=(-stop-after 4s)
  ytester -loop -serve "restart$YTAG"

  slowly () {
    is "`ysend -delay ${1}s slow$1`" = "restart$YTAG"
  }

  # Start the timer.
  slowly 1

  running ytester

  # Stop (and restart) the timer.
  pids=()
  for i in 1 2 3; {
    slowly 2 &
    pids+=($!)
    sleep 1
  }

  running ytester

  wait-each "${pids[@]}"

  sleep 4

  ! running ytester

  knock

  running ytester

  sleep 5

  ! running ytester

  ylog | grep 'stopping: ' | lines 2
}
