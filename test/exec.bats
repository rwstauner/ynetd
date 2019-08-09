#!/usr/bin/env bats

load helpers

@test "get dest addr from exec" {
  cmd=$PWD/test/dest-addr
  export PROXY_PORT
  ynetd -proxy ":$LISTEN_PORT exec:$cmd"
  "$YTESTER" -port "$PROXY_PORT" -serve "exec$YTAG" &
  tester="$!"
  is "`ysend -timeout 1s hello`" = "exec$YTAG"
  kill "$tester" || :
  ylog | grep -qF -- "proxy :$LISTEN_PORT -> exec:$cmd cmd"
}

@test "show error from exec" {
  cmd=$PWD/test/dest-addr-fail
  ynetd -proxy ":$LISTEN_PORT exec:$cmd"
  knock
  ylog | grep -qE -- "failed to get address from $cmd \(.+?\): stderr stuff"
}
