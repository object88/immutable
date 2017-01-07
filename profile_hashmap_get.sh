#!/bin/bash
# Performs the go benchmark tests for hashmap get, with cpu analysis

go test -bench=Benchmark_Hashmap_Get_LargeBlock -run=NONE -cpuprofile=./cpu.prof

go tool pprof -pdf ./immutable.test ./cpu.prof > ./cpu.pdf
