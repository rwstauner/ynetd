#!/usr/bin/env bats

load helpers

@test "command is not required" {
  "${YTESTER:-build/ytester}" -port "$PROXY_PORT" -serve "nocmd$YTAG" &
  tester="$!"

  # Nothing on listen port.
  ysend -timeout 1s hello | lines 0

  # Tester on proxy port.
  ysend -port "$PROXY_PORT" foo | grep -q "nocmd$YTAG"

  ynetd -listen ":$LISTEN_PORT" -proxy "localhost:$PROXY_PORT"
  running ynetd

  ysend -timeout 1s hello | grep -q "nocmd$YTAG"
  kill "$tester" || :

  ylog -y | grep -qF 'cmd: <nil>'
  ylog -y | grep starting: | lines 0
}
