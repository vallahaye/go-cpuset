#!/usr/bin/env bats

bats_require_minimum_version 1.5.0

function cpuset() {
  go -C "$BATS_TEST_DIRNAME" run . "$@"
}

@test "print help and exit" {
  run -0 cpuset -help
  [[ "$lines[0]" = 'Usage: cpuset '* ]]
}

@test "print version and exit" {
  run -0 cpuset -version
  [[ "$output" = 'cpuset '* ]]
}

@test "flag provided but not defined: -undefined" {
  run -1 cpuset -undefined difference 0-32 8-16
  [[ "${lines[0]}" = 'flag provided but not defined: -undefined' ]]
  [[ "${lines[-1]}" = 'exit status 2' ]]
}

@test "flag provided but invalid: -format" {
  run -1 cpuset -format invalid difference 0-32 8-16
  [[ "${lines[0]}" = 'flag provided but invalid: -format' ]]
  [[ "${lines[-1]}" = 'exit status 2' ]]
}

@test "invalid number of arguments" {
  run -1 cpuset
  [[ "${lines[0]}" = 'invalid number of arguments' ]]
  [[ "${lines[-1]}" = 'exit status 2' ]]
}

@test "command provided but not defined: undefined" {
  run -1 cpuset undefined 0-32 8-16
  [[ "${lines[0]}" = 'command provided but not defined: undefined' ]]
  [[ "${lines[-1]}" = 'exit status 2' ]]
}

@test "s1 provided but invalid" {
  run -1 cpuset difference s1 8-16
  [[ "${lines[0]}" =~ 's1 provided but invalid:'* ]]
  [[ "${lines[-1]}" = 'exit status 2' ]]
}

@test "s2 provided but invalid" {
  run -1 cpuset difference 0-32 s2
  [[ "${lines[0]}" =~ 's2 provided but invalid:'* ]]
  [[ "${lines[-1]}" = 'exit status 2' ]]
}

@test "compute the difference of the two cpusets (default format)" {
  run -0 cpuset difference 0-32 8-16
  [[ "$output" = '0-7,17-32' ]]
}

@test "compute the difference of the two cpusets (-format list)" {
  run -0 cpuset -format list difference 0-32 8-16
  [[ "$output" = '0-7,17-32' ]]
}

@test "compute the difference of the two cpusets (-format mask)" {
  run -0 cpuset -format mask difference 00000001,ffffffff 0000ff00
  [[ "$output" = '00000001,ffff00ff' ]]
}

@test "compute the intersection of the two cpusets (default format)" {
  run -0 cpuset intersection 0-32 8-16
  [[ "$output" = '8-16' ]]
}

@test "compute the intersection of the two cpusets (-format list)" {
  run -0 cpuset -format list intersection 0-32 8-16
  [[ "$output" = '8-16' ]]
}

@test "compute the intersection of the two cpusets (-format mask)" {
  run -0 cpuset -format mask intersection 00000001,ffffffff 0000ff00
  [[ "$output" = '0000ff00' ]]
}

@test "compute the union of the two cpusets (default format)" {
  run -0 cpuset union 0-32 8-16
  [[ "$output" = '0-32' ]]
}

@test "compute the union of the two cpusets (-format list)" {
  run -0 cpuset -format list union 0-32 8-16
  [[ "$output" = '0-32' ]]
}

@test "compute the union of the two cpusets (-format mask)" {
  run -0 cpuset -format mask union 00000001,ffffffff 0000ff00
  [[ "$output" = '00000001,ffffffff' ]]
}
