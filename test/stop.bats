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
      is "`ysend hello`" = "stop$YTAG"
      sleep 1
    }

    # Let it expire.
    sleep 3

    ! running ytester

    ylog | grep 'stopping ' | lines $i
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

    ylog | grep 'stopping ' | lines $i
    last="$current"
  }
}
