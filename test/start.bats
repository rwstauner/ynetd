#!/usr/bin/env bats

load helpers

@test "auto start" {
  YARGS=()
  ytester -loop -serve "start$YTAG"

  sleep 0.5
  ! running ytester

  # Nothing on listen port.
  ysend -timeout 1s -port "$PROXY_PORT" hello | lines 0

  close

  YARGS=(-auto-start)
  ytester -loop -serve "start$YTAG"

  sleep 0.5
  running ytester

  # Circumvent ynetd to prove the command is already running.
  is "`ysend -timeout 1s -port "$PROXY_PORT" hello`" = "start$YTAG"
}

@test "auto start from yaml" {
  tmp=`mktemp -t ynetd.XXXXXX`
  cat <<YAML > "$tmp"
---
services:
  - proxy: {":$LISTEN_PORT": "localhost:$PROXY_PORT"}
    command: ["$YAS", "ytester$YTAG", "$YTESTER", "-port", "$PROXY_PORT", "-loop", "-serve", "start$YTAG"]
    auto_start: true
YAML

  ynetd -config "$tmp"

  sleep 0.5
  running ytester

  # Circumvent ynetd to prove the command is already running.
  is "`ysend -timeout 1s -port "$PROXY_PORT" hello`" = "start$YTAG"

  rm "$tmp"
}
