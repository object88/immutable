package immutable

import "fmt"

type RoutineOptions struct {
	Width int
}

func (s *MappedStep) GoMap(f MapPredicate) *MappedStep {
	hr.funcs = append(hr.funcs, &mapStep{predicate})
	return hr
}

func (s *MappedStep) GoReduce(f ReducePredicate) *ReducedStep {
	return nil
}

func (s *MappedStep) ToHashmap() *HashMap {
	return nil
}

// func (s *FilteredStep) ToList() interface{} {}
// Func (*FilteredStep) ToSet() *Set {}

type MappedStep struct {
}

type ReducedStep struct {
}

type routine struct {
	options *RoutineOptions
	base    Base
	funcs   []step
	err     *error
}

type step interface {
	execute(kvp <-chan keyValuePair)
}

type mapStep struct {
	f MapPredicate
}

func (ms *mapStep) execute(kvp <-chan keyValuePair) {
	x := <-kvp
	key, value := x.key, x.value
	result, err := ms.f(key, value)
}

func (hr *HashRoutine) Result() (*HashMap, error) {
	routineCount := 10
	if routineCount > hr.base.Size() {
		routineCount = hr.base.Size()
	}

	var err error
	chans := make([]chan keyValuePair, routineCount)
	abort := make(chan error)
	ch := make(chan *HashMap)
	go func() {
		defer close(ch)

		for _, v := range chans {
		}

		for i := 0; i < hr.base.Size(); i++ {
			select {
			case res := <-resc:
				fmt.Println(res)
			case err = <-abort:
				return
			}
		}

	}()

	result := <-ch

	return result, nil
}

func (*HashRoutine) ToList() (*ListRoutine, error) {
	return nil, nil
}

func (*routine) Error() error {
	return nil
}
