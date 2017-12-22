#!/usr/bin/env bats

load helpers

@test "arg: -proxy" {
  listen2="$PROXY_PORT" # ssh...
  ynetd -proxy ":$LISTEN_PORT localhost:5002 :$listen2 localhost:5004"
  close
  ylog | grep -qF -- "proxy :$LISTEN_PORT -> localhost:5002 cmd"
  ylog | grep -qF -- "proxy :$listen2 -> localhost:5004 cmd"
}

@test "arg: -proxy-sep" {
  listen2="$PROXY_PORT" # ssh...
  ynetd -proxy-sep , -proxy ":$LISTEN_PORT,localhost:5000,:$listen2,localhost:7000"
  close
  ylog | grep -qF -- "proxy :$LISTEN_PORT -> localhost:5000 cmd"
  ylog | grep -qF -- "proxy :$listen2 -> localhost:7000 cmd"
}
