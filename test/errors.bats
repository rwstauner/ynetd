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
