package immutable_test

import (
	"testing"

	"github.com/object88/immutable"
)

func Test_Hashmap_ForEach(t *testing.T) {
	contents := map[int]string{
		0: "false",
		1: "false",
		2: "false",
		3: "false",
	}
	original := immutable.NewIntToStringHashmap(contents)
	if original == nil {
		t.Fatal("Failed to create hashmap")
	}
	original.ForEach(func(k int, v string) {
		if contents[k] == "true" {
			t.Fatalf("At %s, already visited\n", k)
		}
		contents[k] = "true"
	})

	for k, v := range contents {
		if v != "true" {
			t.Fatalf("At %s, not visited\n", k)
		}
	}
}
