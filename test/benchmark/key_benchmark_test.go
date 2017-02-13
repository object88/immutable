package benchmark

import (
	"encoding/binary"
	"hash/adler32"
	"hash/fnv"
	"testing"
)

var hash32 uint32
var hash64 uint64

func Benchmark_Hashing_FNV32(b *testing.B) {
	var result uint32
	for i := 0; i < b.N; i++ {
		hasher := fnv.New32()

		binary.Write(hasher, binary.LittleEndian, uint32(i))
		result = hasher.Sum32()
	}
	hash32 = result
}

func Benchmark_Hashing_FNV32a(b *testing.B) {
	var result uint32
	for i := 0; i < b.N; i++ {
		hasher := fnv.New32a()

		binary.Write(hasher, binary.LittleEndian, uint32(i))
		result = hasher.Sum32()
	}
	hash32 = result
}

func Benchmark_Hashing_FNV64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hasher := fnv.New64()

		binary.Write(hasher, binary.LittleEndian, uint64(i))
		hash64 = hasher.Sum64()
	}
}

func Benchmark_Hashing_FNV64a(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hasher := fnv.New64a()

		binary.Write(hasher, binary.LittleEndian, uint64(i))
		hash64 = hasher.Sum64()
	}
}

func Benchmark_Hashing_Adler32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hasher := adler32.New()

		binary.Write(hasher, binary.LittleEndian, uint32(i))
		hash32 = hasher.Sum32()
	}
}
