#!/usr/bin/env bats

load helpers

@test "no services" {
  "$YNETD" | grep -q 'no services configured'
}

@test "no address" {
  ynetd -proxy ":$LISTEN_PORT "
  knock
  # Should be instant, rather than timing out.
  ylog | grep -qF -- "error starting listener: destination address is required"
}

@test "no port" {
  ynetd -proxy ":$LISTEN_PORT foobar"
  knock
  # Should be instant, rather than timing out.
  ylog | grep -qF -- "dial tcp: address foobar: missing port in address"
}

@test "bad host" {
  ynetd -timeout 0s -proxy ":$LISTEN_PORT .:1"
  knock
  sleep 1
  ylog | grep -qF -- "dial tcp: lookup .: no such host"
}

@test "command not found" {
  cmd="ynetd-test-cmd-that-should-not-exist-$YTAG"
  ynetd -proxy ":$LISTEN_PORT localhost:$PROXY_PORT" -timeout 5s -wait-after-start 5s -auto-start $cmd arg

  logmsg="error starting [$cmd arg]: exec: \"$cmd\": executable file not found in \$PATH"
  # auto-start should have already errored
  ylog | grep -qF -- "$logmsg"

  start=`date +%s`
  is "`ysend ready`" = "" # connection closed
  end=`date +%s`
  is $((end - start)) -lt 2 # faster than timeout

  # 1 for auto-start and 1 for ysend.
  ylog | grep -F -- "$logmsg" | lines 2
}
