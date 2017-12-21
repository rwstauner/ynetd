#!/usr/bin/env bats

load helpers

@test "deprecated: -listen" {
  ynetd -listen ":5000" -proxy "localhost:5001"
  close
  ylog | grep -qE -- '^-listen is deprecated'
  ylog | grep -qF -- 'proxy :5000 -> localhost:5001 cmd'
}
