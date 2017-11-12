#!/usr/bin/env bats

load helpers

@test "works before timeout" {
  ynetbash -t 6s 'sleep 4; echo "timely$YTAG" | listen'

  start=`date +%s`
  # This will wait.
  ysend | grep -qFx "timely$YTAG"
  end=`date +%s`

  [ $end -ge $((start + 4)) ]
}

@test "times out" {
  ynetbash -t 10ms 'sleep 10'
  knock
  knock
  knock

  sleep 1
  running ynetbash

  # One for each attempt.
  [ 3 -eq `ylog | grep 'ynetd: timed out after 10ms' | wc -l` ]
}

@test "times out, works later" {
  ynetbash -t 3s 'sleep 5; echo "timely$YTAG" | listen'
  knock

  # Command is running.
  running ynetbash

  # Wait for first sleep.
  sleep 4

  ylog | grep 'ynetd: timed out after 3s'

  # Wait for listen to start.
  sleep 2

  ysend hello | grep -qFx "timely$YTAG"
}
