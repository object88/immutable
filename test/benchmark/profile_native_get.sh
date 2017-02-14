#!/bin/bash
# Performs the go benchmark tests for hashmap get, with cpu analysis

go test -bench=Benchmark_Hashmap_Get_Native/n=5000000 -run=NONE -benchtime=10s -cpuprofile=./cpu_native_get.prof

go tool pprof -pdf ./benchmark.test ./cpu_native_get.prof > ./cpu_native_get.pdf
