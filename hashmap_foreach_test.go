package immutable

import "testing"

func Test_Hashmap_ForEach(t *testing.T) {
	data := map[Key]Value{
		IntKey(0): false,
		IntKey(1): false,
		IntKey(2): false,
		IntKey(3): false,
	}
	original := NewHashMap(data)
	if original == nil {
		t.Fatal("Failed to create hashmap")
	}
	original.ForEach(func(k Key, v Value) {
		if v.(bool) {
			t.Fatalf("At %s, already visited\n", k)
		}
		data[k] = true
	})

	for k, v := range data {
		if !v.(bool) {
			t.Fatalf("At %s, not visited\n", k)
		}
	}
}
