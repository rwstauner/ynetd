#!/usr/bin/env bats

yversion () {
  $YNETD --version
}

@test "displays version" {
  yversion | grep -qEx 'ynetd [0-9]+\.[0-9]+\.[0-9]+'
  [ 1 -eq `yversion | wc -l` ]
}
