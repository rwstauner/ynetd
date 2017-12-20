#!/usr/bin/env bats

load helpers

@test "config file" {
  tmp=`mktemp -t ynetd.XXXXXX`
  cat <<JSON > "$tmp"
{
  "Services": [
    {
      "Proxy": {
        ":$LISTEN_PORT": "localhost:$PROXY_PORT"
      },
      "Command": ["$YAS", "ytester1$YTAG", "$YTESTER", "-port", "$PROXY_PORT", "-loop", "-serve", "json1$YTAG"],
      "Timeout": "3s"
    }
  ]
}
JSON
  cat "$tmp" >&2

  ynetd -config "$tmp"

  running ynetd
  ! running ytester1

  # Wait for service to be done (better than a sleep).
  is "`ysend hello`" = "json1$YTAG"

  running ytester1

  close
  ylog -y | grep starting: | lines 1

  rm "$tmp"
}
