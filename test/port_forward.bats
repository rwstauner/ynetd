#!/usr/bin/env bats

load helpers

@test "command is not required" {
  "$YTESTER" -port "$PROXY_PORT" -loop -serve "nocmd$YTAG" &
  tester="$!"

  # Nothing on listen port.
  ysend -timeout 1s hello | lines 0

  # Tester on proxy port.
  is "`ysend -port "$PROXY_PORT" foo`" = "nocmd$YTAG"

  ynetd -proxy ":$LISTEN_PORT localhost:$PROXY_PORT"
  running ynetd

  is "`ysend -timeout 1s hello`" = "nocmd$YTAG"
  kill "$tester" || :

  ylog -y | grep -qF 'cmd: <nil>'
  ylog -y | grep starting: | lines 0
}
