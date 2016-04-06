#!/bin/bash

OSLIST=(linux_x86_64 darwin_x86_64 windows_x86_64.exe)
SUMS=("md5" "shasum -a 1" "shasum -a 256")

for OS in "${OSLIST[@]}"
do
  for SUM in "${SUMS[@]}"
  do
    $SUM "cig_$OS" > "cig_$OS.
  done
done

