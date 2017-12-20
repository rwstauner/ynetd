#!/usr/bin/env bats

load helpers

@test "config file" {
  tmp=`mktemp -t ynetd.XXXXXX`
  listen2=$((LISTEN_PORT+1))
  proxy2=$((PROXY_PORT+1))
  cat <<JSON > "$tmp"
{
  "Services": [
    {
      "Proxy": {
        ":$LISTEN_PORT": "localhost:$PROXY_PORT"
      },
      "Command": ["$YAS", "ytester1$YTAG", "$YTESTER", "-port", "$PROXY_PORT", "-loop", "-serve", "json1$YTAG"],
      "Timeout": "3s"
    },
    {
      "Proxy": {
        ":$listen2": "localhost:$proxy2"
      },
      "Command": ["$YAS", "ytester2$YTAG", "$YTESTER", "-port", "$proxy2", "-loop", "-serve", "json2$YTAG"],
      "Timeout": "4s"
    }
  ]
}
JSON
  cat "$tmp" >&2

  ynetd -config "$tmp"

  running ynetd
  ! running ytester1
  ! running ytester2

  # Wait for service to be done (better than a sleep).
  is "`ysend hello`" = "json1$YTAG"

  running ytester1
  ! running ytester2

  is "`ysend -port "$listen2" hello`" = "json2$YTAG"

  running ytester1
  running ytester2

  close
  ylog -y | grep starting: | lines 2

  rm "$tmp"
}
