package immutable

// BoolValue is a boolean.
type BoolValue bool

// StringValue is a string.
type StringValue string

// func DehydrateBool(v Value) (result uint64, err error) {
// 	if bool(v.(BoolValue)) {
// 		return 1, nil
// 	}
// 	return 0, nil
// }
//
// func HydrateBool(value uint64) (result Value, err error) {
// 	return BoolValue(value == 1), nil
// }
//
func (b BoolValue) String() string {
	if bool(b) {
		return "true"
	}
	return "false"
}

//
// func DehydrateString(value Value) (result unsafe.Pointer, err error) {
// 	return unsafe.Pointer(&value), nil
// }
//
// func HydrateString(value unsafe.Pointer) (result Value, err error) {
// 	return nil, nil
// }

func (s StringValue) String() string {
	return string(s)
}
