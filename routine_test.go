package immutable

import "testing"

func Test_Routine(t *testing.T) {
	contents := map[Key]Value{
		IntKey(1): "a",
		IntKey(2): "b",
	}
	original := NewHashMap(contents, nil)
	original.GoMap(func(k Key, v Value) (Value, error) {
		return v, nil
	}, nil).GoFilter(func(k Key, v Value) (bool, error) {
		return true, nil
	})
}
