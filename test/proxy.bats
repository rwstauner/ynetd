#!/usr/bin/env bats

load helpers

@test "arg: -proxy" {
  ynetd -proxy ":5001 localhost:5002 :5003 localhost:5004"
  close
  ylog | grep -qF -- 'proxy :5001 -> localhost:5002 cmd'
  ylog | grep -qF -- 'proxy :5003 -> localhost:5004 cmd'
}

@test "arg: -proxy-sep" {
  ynetd -proxy-sep , -proxy :4000,localhost:5000,:6000,localhost:7000
  close
  ylog | grep -qF -- 'proxy :4000 -> localhost:5000 cmd'
  ylog | grep -qF -- 'proxy :6000 -> localhost:7000 cmd'
}
