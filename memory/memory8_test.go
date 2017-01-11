package memory

import (
	"math"
	"math/rand"
	"testing"
	"time"
)

func Test_Small_Random(t *testing.T) {
	src := rand.NewSource(time.Now().UnixNano())
	random := rand.New(src)

	width := uint32(64 - 11)
	max := int64(math.Pow(2.0, float64(width)))

	count := 8
	contents := make([]uint64, count)
	for i := 0; i < count; i++ {
		contents[i] = uint64(random.Int63n(max))
	}

	m := AllocateMemories(SmallBlock, width, uint32(count))
	for k, v := range contents {
		m.Assign(uint64(k), v)
	}

	for k, v := range contents {
		result := m.Read(uint64(k))
		if result != v {
			t.Fatalf("At %d\nexpected %064b\nreceived %064b\n", k, v, result)
		}
	}
}

// func Test_Read(t *testing.T) {
// 	m := AllocateMemories(4, 8)
// 	m.m[0] = 0x0f
// 	result := m.Read(0)
// 	if result != 0x0f {
// 		t.Fatalf("Incorrect result from read; expected %x, got %x\n", 0x0f, result)
// 	}
// }
//
// func Test_Read_AtSecondUint32(t *testing.T) {
// 	m := AllocateMemories(4, 16)
// 	m.m[4] = 0x0f
// 	result := m.Read(8)
// 	if result != 0x0f {
// 		t.Fatalf("Incorrect result from read; expected %x, got %x\n", 0x0f, result)
// 	}
// }
//
// func Test_Read_Offset(t *testing.T) {
// 	m := AllocateMemories(4, 8)
// 	m.m[0] = 0xf0
// 	result := m.Read(1)
// 	if result != 0x0f {
// 		t.Fatalf("Incorrect result from read; expected %x, got %x\n", 0x0f, result)
// 	}
// }
//
// func Test_Read_AtSecondByte(t *testing.T) {
// 	m := AllocateMemories(4, 16)
// 	m.m[4] = 240
// 	result := m.Read(9)
// 	if result != 15 {
// 		t.Fatalf("Incorrect result from read; expected %x, got %x\n", 15, result)
// 	}
// }
//
// func Test_Read_Over2Bytes(t *testing.T) {
// 	// Allocate 12 bits of memory across 2 bytes
// 	m := AllocateMemories(6, 2)
// 	m.m[0] = 0xc0
// 	m.m[1] = 0x0f
// 	result := m.Read(1)
// 	if result != 0x3f {
// 		t.Fatalf("Incorect result from read; expected %x, got %x\n", 0x3f, result)
// 	}
// }
//
// func Test_Read_Over3Bytes(t *testing.T) {
// 	// Allocate 22 bits of memory across 3 byte
// 	m := AllocateMemories(11, 2)
// 	m.m[0] = 0x00
// 	m.m[1] = 0xf8
// 	m.m[2] = 0x3f
// 	result := m.Read(1)
// 	if result != 0x7ff {
// 		t.Fatalf("Incorect result from read; expected %x, got %x\n", 0x7ff, result)
// 	}
// }
//
// func Test_Assign(t *testing.T) {
// 	m := AllocateMemories(4, 1)
// 	m.Assign(0, 0xf)
// 	result := m.m[0]
// 	if result != 0xf {
// 		t.Fatalf("Incorrect result from write; expected %x, got %x\n", 0xf, result)
// 	}
// }
//
// func Test_Assign_Offset(t *testing.T) {
// 	m := AllocateMemories(4, 2)
// 	m.Assign(1, 0xf)
// 	result := m.m[0]
// 	if result != 0xf0 {
// 		t.Fatalf("Incorrect result from write; expected %x, got %x\n", 0xf0, result)
// 	}
// }
//
// func Test_Assign_AtSecondByte(t *testing.T) {
// 	m := AllocateMemories(4, 4)
// 	m.Assign(2, 0xf)
// 	result := m.m[1]
// 	if result != 0xf {
// 		t.Fatalf("Incorrect result from write; expected %x, got %x\n", 0xf, result)
// 	}
// }
//
// func Test_Assign_Over2Bytes(t *testing.T) {
// 	m := AllocateMemories(6, 2)
// 	m.Assign(1, 0x3f)
// 	result1 := m.m[0]
// 	if result1 != 0xc0 {
// 		t.Fatalf("Incorrect result from write; expected %x, got %x\n", 0xc0, result1)
// 	}
// 	result2 := m.m[1]
// 	if result2 != 0x0f {
// 		t.Fatalf("Incorrect result from write; expected %x, got %x\n", 0x0f, result2)
// 	}
// }
//
// func Test_Assign_Over3Bytes(t *testing.T) {
// 	m := AllocateMemories(11, 2)
// 	m.Assign(1, 0x7ff)
// 	result1 := m.m[0]
// 	if result1 != 0x00 {
// 		t.Fatalf("Incorrect result from write; expected %x, got %x\n", 0x00, result1)
// 	}
// 	result2 := m.m[1]
// 	if result2 != 0xf8 {
// 		t.Fatalf("Incorrect result from write; expected %x, got %x\n", 0xf8, result2)
// 	}
// 	result3 := m.m[2]
// 	if result2 != 0x3f {
// 		t.Fatalf("Incorrect result from write; expected %x, got %x\n", 0x3f, result3)
// 	}
// }
//
// func Test_ReadBits(t *testing.T) {
// 	mask := generateMask(2, 2)
//
// 	evaluateReadBits(t, 12, 255, mask)
// 	evaluateReadBits(t, 0, 0, mask)
// 	evaluateReadBits(t, 4, 4, mask)
// 	evaluateReadBits(t, 8, 8, mask)
// }
//
// func evaluateReadBits(t *testing.T, expected, initial, mask byte) {
// 	result := readBits(initial, mask)
// 	if expected != result {
// 		t.Fatalf("Expected %d from readBits(%d, %d); got %d", expected, initial, mask, result)
// 	}
// }
//
// func Test_WriteBits(t *testing.T) {
// 	evaluateWriteBits(t, 1, 0, 1, 1, 0)
// 	evaluateWriteBits(t, 2, 0, 1, 1, 1)
// 	evaluateWriteBits(t, 6, 0, 3, 2, 1)
// 	evaluateWriteBits(t, 253, 255, 0, 1, 1)
// }
//
// func evaluateWriteBits(t *testing.T, expected, original, newValue, width, shift byte) {
// 	result := writeBits(original, newValue, width, shift)
// 	if expected != result {
// 		t.Fatalf("Expected %d from writeBits(%d, %d, %d, %d); got %d", expected, original, newValue, width, shift, result)
// 	}
// }
//
// func Test_GenerateMask(t *testing.T) {
// 	evaluateGenerateMask(t, 1, 1, 0)  // 00001
// 	evaluateGenerateMask(t, 3, 2, 0)  // 00011
// 	evaluateGenerateMask(t, 7, 3, 0)  // 00111
// 	evaluateGenerateMask(t, 2, 1, 1)  // 00010
// 	evaluateGenerateMask(t, 6, 2, 1)  // 00110
// 	evaluateGenerateMask(t, 14, 3, 1) // 01110
// 	evaluateGenerateMask(t, 4, 1, 2)  // 00100
// 	evaluateGenerateMask(t, 12, 2, 2) // 01100
// 	evaluateGenerateMask(t, 28, 3, 2) // 11100
// }
//
// func evaluateGenerateMask(t *testing.T, expected, width, shift byte) {
// 	result := generateMask(width, shift)
// 	if expected != result {
// 		t.Fatalf("Expected %d from generateMask(%d, %d); got %d", expected, width, shift, result)
// 	}
// }
