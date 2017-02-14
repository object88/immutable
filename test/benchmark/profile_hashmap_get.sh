#!/bin/bash
# Performs the go benchmark tests for hashmap get, with cpu analysis

go test -bench=Benchmark_Hashmap_Get_NotPackedBlock/n=5000000 -run=NONE -cpuprofile=./cpu.prof

go tool pprof -pdf ./benchmark.test ./cpu.prof > ./cpu.pdf
