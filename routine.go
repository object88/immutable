package immutable

type routine struct {
	options *RoutineOptions
	base    Base
	funcs  []
	err *error
}

type RoutineOptions struct {
	Width int
}

type HashRoutine struct {
	routine
}

type ListRoutine struct {
	routine
}

func (*HashRoutine) GoFilter(f FilterPredicate) *HashRoutine {
	return nil
}

func (hr *HashRoutine) GoMap(predicate MapPredicate) *HashRoutine {
	// mutated := b.instantiate(b.Size(), nil)
	// abort := make(chan struct{})
	// ch := hr.base.iterate(abort)
	// for kvp := range ch {
	// 	newV, err := go predicate(kvp.key, kvp.value)
	// 	if err != nil {
	// 		hr.err = &err
	// 		close(abort)
	//
	// 		return hr
	// 	}
	// 	mutated.internalSet(kvp.key, newV)
	// }
	// return hr
}

func (*HashRoutine) GoReduce(acc Value, f ReducePredicate) (Value, error) {
	return nil, nil
}

func (*HashRoutine) Result() (*HashMap, error) {
	return nil, nil
}

func (*HashRoutine) ToList() (*ListRoutine, error) {
	return nil, nil
}

func (*ListRoutine) GoFilter(f FilterPredicate) *ListRoutine {
	return nil
}

func (*ListRoutine) GoMap(f MapPredicate) *ListRoutine {
	return nil
}

func (*ListRoutine) GoReduce(acc Value, f ReducePredicate) (Value, error) {
	return nil, nil
}

func (*ListRoutine) Result() (interface{}, error) {
	return nil, nil
}

func (*ListRoutine) ToHashmap() (*HashMap, error) {
	return nil, nil
}

func (*routine) Error() error {
	return nil
}
