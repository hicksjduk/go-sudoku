package valueset

import (
	"sync"
)

type ValueSet struct {
	lock   func(func())
	min    int
	max    int
	values []uint64
}

func NewValueSet(min int, max int) ValueSet {
	size := max - min + 1
	return ValueSet{
		lock:   newLocker(),
		min:    min,
		max:    max,
		values: make([]uint64, size/64+1),
	}
}

func newLocker() func(func()) {
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

func (vs *ValueSet) filterOutOfBounds(values ...int) []int {
	answer := make([]int, 0, len(values))
	for _, v := range values {
		if vs.isInBounds(v) {
			answer = append(answer, v)
		}
	}
	return answer
}

func (vs *ValueSet) isInBounds(v int) bool {
	return v >= vs.min && v <= vs.max
}

func (vs *ValueSet) Add(values ...int) []int {
	addFunc := func(slot uint64, bit uint) uint64 {
		return slot | (1 << bit)
	}
	return vs.apply(addFunc, vs.filterOutOfBounds(values...))
}

func (vs *ValueSet) Remove(values ...int) []int {
	addFunc := func(slot uint64, bit uint) uint64 {
		return slot &^ (1 << bit)
	}
	return vs.apply(addFunc, vs.filterOutOfBounds(values...))
}

func (vs *ValueSet) apply(f func(uint64, uint) uint64, values []int) []int {
	changed := make([]int, 0, len(values))
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
	return changed
}

func (vs *ValueSet) Values() (answer []int) {
	answer = make([]int, 0, vs.max-vs.min+1)
	var vals []uint64
	vs.lock(func() {
		vals = vs.values
	})
	num := vs.min
	for _, slot := range vals {
		for b, mask := 0, uint64(1); num <= vs.max && b < 64; b, num, mask = b+1, num+1, mask<<1 {
			if bitIsSet := (slot & mask) != 0; bitIsSet {
				answer = append(answer, num)
			}
		}
	}
	return
}

func (vs *ValueSet) Contains(value int) bool {
	if !vs.isInBounds(value) {
		return false
	}
	var vals []uint64
	vs.lock(func() {
		vals = vs.values
	})
	slot, bit := vs.address(value)
	return (vals[slot] & (1 << bit)) != 0
}
