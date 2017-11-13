#!/usr/bin/env bats

yversion () {
  $YNETD --version
}

@test "displays version" {
  yversion >&2
  yversion | grep -qEx 'ynetd v?[0-9]+(\.[0-9]+)+(-g[0-9a-f]+)?'
  [ 1 -eq `yversion | wc -l` ]
}
