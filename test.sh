#!/bin/bash

org=${PWD%/*}
org=${org##*/}
repository=${PWD##*/}
echo "** org:$org"
echo "** repository:$repository"
echo

# 重新造一遍 go mod
sh ./shell/gen-proto.sh
sh ./shell/configure.sh

# test
cd ./src
go test -v -bench=".*" ./database/database_test.go ./database/database.go