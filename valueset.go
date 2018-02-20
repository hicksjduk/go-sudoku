package sudoku

import (
	"fmt"
	"sync"
)

type locker func(func())

type ValueSet struct {
	lock   locker
	min    int
	max    int
	values []uint64
}

type ValueOutOfRangeError struct {
	message string
}

func (vs *ValueSet) error(value int) *ValueOutOfRangeError {
	msg := fmt.Sprintf("Value %d is out of range (supported bounds are %d and %d)", value, vs.min, vs.max)
	return &ValueOutOfRangeError{msg}
}

func (e *ValueOutOfRangeError) Error() string {
	return e.message
}

func newValueSet(min int, max int) ValueSet {
	size := max - min + 1
	return ValueSet{
		lock:   newLocker(),
		min:    min,
		max:    max,
		values: make([]uint64, size/64+1),
	}
}

func newLocker() locker {
	lock := sync.Mutex{}
	return func(action func()) {
		lock.Lock()
		defer lock.Unlock()
		action()
	}
}

func (vs *ValueSet) address(value int) (int, uint) {
	vIndex := value - vs.min
	return vIndex / 64, uint(vIndex % 64)
}

func (vs *ValueSet) checkInBounds(values ...int) (err *ValueOutOfRangeError) {
	for _, v := range values {
		if outOfBounds := v < vs.min || v > vs.max; outOfBounds {
			err = vs.error(v)
			return
		}
	}
	return
}

func (vs *ValueSet) Add(values ...int) (added []int, err error) {
	if err := vs.checkInBounds(values...); err != nil {
		return
	}
	addFunc := func(slot uint64, bit uint) uint64 {
		return slot | (1 << bit)
	}
	added = vs.apply(addFunc, values)
	return
}

func (vs *ValueSet) Remove(values ...int) (added []int, err error) {
	if err := vs.checkInBounds(values...); err != nil {
		return
	}
	addFunc := func(slot uint64, bit uint) uint64 {
		return slot &^ (1 << bit)
	}
	added = vs.apply(addFunc, values)
	return
}

func (vs *ValueSet) apply(f func(uint64, uint) uint64, values []int) (changed []int) {
	changed = make([]int, 0, len(values))
	vs.lock(func() {
		vals := vs.values
		for _, v := range values {
			slot, bit := vs.address(v)
			if sval := f(vals[slot], bit); vals[slot] != sval {
				changed = append(changed, v)
				vals[slot] = sval
			}
		}
		vs.values = vals
	})
	return
}

func (vs *ValueSet) Values() (answer []int) {
	answer = make([]int, vs.max-vs.min+1)
	var vals []uint64
	vs.lock(func() {
		vals = vs.values
	})
	num := vs.min
	for _, slot := range vals {
		mask := uint64(1)
		for b := 0; num <= vs.max && b < 63; b++ {
			mask <<= 1
			if bitIsSet := (slot & mask) != 0; bitIsSet {
				answer = append(answer, num)
			}
			num++
		}
	}
	return
}

func (vs *ValueSet) Contains(value int) bool {
	if err := vs.checkInBounds(value); err != nil {
		return false
	}
	var vals []uint64
	vs.lock(func() {
		vals = vs.values
	})
	slot, bit := vs.address(value)
	return (vals[slot] & (1 << bit)) != 0
}
