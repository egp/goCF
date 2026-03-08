#!/usr/bin/env zsh
# tools/test_coverage.sh

go test -coverprofile=./tmp/cover.out ./cf
go tool cover -func=./tmp/cover.out | sort -k3 > ./tmp/cover.func.txt
grep -v "100.0%" ./tmp/cover.func.txt | head -30 > ./tmp/cover.top30.txt

#eof test_coverage.sh