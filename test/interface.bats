#!/usr/bin/env bats

load helpers

@test "proxy with interface" {
  if ! which ifconfig; then
    skip "ifconfig not present"
  fi
  iface=`ifconfig | grep -Eo '^(lo0?)\b'`
  if [[ -z "$iface" ]]; then
    skip "cannot find loopback"
  fi

  ynetd -proxy "interface:$iface:$LISTEN_PORT localhost:5002"
  close
  ylog | grep -qF -- "proxy 127.0.0.1:$LISTEN_PORT -> localhost:5002 cmd"
}

@test "proxy with bad interface" {
  iface=narf
  ynetd -proxy "interface:$iface:$LISTEN_PORT localhost:5002"
  ylog | grep -qF -- "no such network interface"
  ! running ynetd
}

