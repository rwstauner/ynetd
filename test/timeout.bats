#!/usr/bin/env bats

load helpers

@test "works before timeout" {
  YARGS=(-t 6s)
  ytester -before 4s -serve "timely$YTAG"

  start=`date +%s`
  # This will wait.
  is "`ysend "hello"`" = "timely$YTAG"
  end=`date +%s`

  is $end -ge $((start + 4))
}

@test "times out" {
  YARGS=(-t 10ms)
  ytester -before 10s
  knock
  knock
  knock

  sleep 1
  running ytester

  # One for each attempt.

  ylog -y | grep 'timed out after 10ms' | lines 3
}

@test "times out, works later" {
  YARGS=(-t 3s)
  ytester -before 5s -serve "timely$YTAG"
  knock

  # Command is running.
  running ytester

  # Wait for first sleep.
  sleep 4

  ylog -y | grep 'timed out after 3s'

  # Wait for listen to start.
  sleep 2

  is "`ysend hello`" = "timely$YTAG"
}
