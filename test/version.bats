#!/usr/bin/env bats

load helpers

yversion () {
  "$YNETD" --version
}

@test "displays version" {
  yversion >&2
  yversion | grep -qEx 'ynetd v?[0-9]+(\.[0-9]+)+(-g[0-9a-f]+)?'
  yversion | lines 1
}
