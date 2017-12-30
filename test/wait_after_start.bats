#!/usr/bin/env bats

load helpers

# Verify that ytester waits properly.
@test "no waiting" {
  YARGS=()
  ytester -serve "wait$YTAG" -serve-after 2s

  ! running ytester

  start=`date +%s`
  is "`ysend ready`" = "not yet"
  end=`date +%s`

  is $((end - start)) -lt 2
}

@test "wait after start" {
  YARGS=(-wait-after-start 2500ms)
  ytester -loop -serve "wait$YTAG" -serve-after 2s

  ! running ytester

  # First connection waits.
  start=`date +%s`
  is "`ysend ready`" = "wait$YTAG"
  end=`date +%s`

  is $((end - start)) -ge 2

  # Second does not.
  start=`date +%s`
  is "`ysend ready`" = "wait$YTAG"
  end=`date +%s`

  is $((end - start)) -lt 2

  kill -s INT `ypidof ytester`
  ! running ytester

  # Each new start should wait.
  start=`date +%s`
  is "`ysend ready`" = "wait$YTAG"
  end=`date +%s`

  is $((end - start)) -ge 2
}
