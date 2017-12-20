#!/usr/bin/env bats

load helpers

@test "no services" {
  "$YNETD" | grep -q 'no services configured'
}
