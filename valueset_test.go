package sudoku

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"sort"
)

func TestValueSet_Add_None(t *testing.T) {
	vs := NewValueSet(1, 9)
	ch := vs.Add()
	assertEqual(t, []int{}, ch)
	assertEqual(t, []int{}, vs.Values())
}

func TestValueSet_Add_SomeToEmpty(t *testing.T) {
	vs := NewValueSet(1, 9)
	ch := vs.Add(1, 3)
	assertEqual(t, []int{1, 3}, ch)
	assertEqual(t, []int{1, 3}, vs.Values())
}

func TestValueSet_Add_SomeToNonEmptyNoChange(t *testing.T) {
	vs := NewValueSet(1, 9)
	vs.Add(1, 3, 5, 4)
	ch := vs.Add(1, 3)
	assertEqual(t, []int{}, ch)
	assertEqual(t, []int{1, 3, 4, 5}, vs.Values())
}

func TestValueSet_Add_SomeToNonEmptyWithChange(t *testing.T) {
	vs := NewValueSet(1, 9)
	vs.Add(1, 3, 5, 4)
	ch := vs.Add(1, 6, 3)
	assertEqual(t, []int{6}, ch)
	assertEqual(t, []int{1, 3, 4, 5, 6}, vs.Values())
}

func TestValueSet_Add_WithInvalid(t *testing.T) {
	vs := NewValueSet(1, 9)
	ch := vs.Add(1, 6, 10, -4, 1131, 3, 0)
	assertEqual(t, []int{1, 3, 6}, ch)
	assertEqual(t, []int{1, 3, 6}, vs.Values())
}

func TestValueSet_Remove_FromEmpty(t *testing.T) {
	vs := NewValueSet(1, 9)
	ch := vs.Remove(5)
	assertEqual(t, []int{}, ch)
	assertEqual(t, []int{}, vs.Values())
}

func TestValueSet_Remove_FromNonEmptyNoChange(t *testing.T) {
	vs := NewValueSet(1, 9)
	vs.Add(1, 3, 5, 7)
	ch := vs.Remove(2, 4, 6, 8)
	assertEqual(t, []int{}, ch)
	assertEqual(t, []int{1, 3, 5, 7}, vs.Values())
}

func TestValueSet_Remove_FromNonEmptyWithChange(t *testing.T) {
	vs := NewValueSet(1, 9)
	vs.Add(1, 3, 5, 7)
	ch := vs.Remove(2, 3, 4)
	assertEqual(t, []int{3}, ch)
	assertEqual(t, []int{1, 5, 7}, vs.Values())
}

func TestValueSet_Remove_WithInvalid(t *testing.T) {
	vs := NewValueSet(1, 9)
	vs.Add(1, 3, 2, 5)
	ch := vs.Remove(2, 0, -225, 147192547, 10, 6, 5)
	assertEqual(t, []int{2, 5}, ch)
	assertEqual(t, []int{1, 3}, vs.Values())
}

func TestValueSet_Contains_NotFound(t *testing.T) {
	vs := NewValueSet(1, 9)
	vs.Add(1, 3)
	found := vs.Contains(4)
	assert.False(t, found)
}

func TestValueSet_Contains_Found(t *testing.T) {
	vs := NewValueSet(1, 9)
	vs.Add(1, 3)
	found := vs.Contains(3)
	assert.True(t, found)
}

func TestValueSet_Contains_JustTooSmall(t *testing.T) {
	vs := NewValueSet(1, 9)
	vs.Add(1, 3)
	found := vs.Contains(0)
	assert.False(t, found)
}

func TestValueSet_Contains_JustTooBig(t *testing.T) {
	vs := NewValueSet(1, 9)
	vs.Add(1, 3)
	found := vs.Contains(10)
	assert.False(t, found)
}

func TestAddWithLargeRange(t *testing.T) {
	vs := NewValueSet(-1415, 124714)
	ch := vs.Add(vs.min, vs.max, 44)
	assertEqual(t, []int{vs.min, 44, vs.max}, ch)
	assertEqual(t, []int{vs.min, 44, vs.max}, vs.Values())
}

func TestContainsWithLargeRange(t *testing.T) {
	vs := NewValueSet(-1415, 124714)
	vs.Add(vs.min, vs.max, 44)
	assert.True(t, vs.Contains(vs.min))
	assert.False(t, vs.Contains(vs.max - 1))
}

func TestRemoveWithLargeRange(t *testing.T) {
	vs := NewValueSet(-1415, 124714)
	vs.Add(vs.min, vs.max, 44)
	ch := vs.Remove(44, vs.min)
	assertEqual(t, []int{vs.min, 44}, ch)
	assertEqual(t, []int{vs.max}, vs.Values())
}

func assertEqual(t *testing.T, expected, actual []int) {
	sort.Ints(expected)
	sort.Ints(actual)
	if assert.Equal(t, len(expected), len(actual)) {
		for i, v := range expected {
			assert.Equal(t, v, actual[i])
		}
	}
}
