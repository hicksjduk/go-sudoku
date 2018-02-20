package sudoku

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValueSet_Add_None(t *testing.T) {
	vs := NewValueSet(1, 9)
	ch, err := vs.Add()
	assert.Nil(t, err)
	assertEqual(t, []int{}, ch)
	assertEqual(t, []int{}, vs.Values())
}

func TestValueSet_Add_SomeToEmpty(t *testing.T) {
	vs := NewValueSet(1, 9)
	ch, err := vs.Add(1, 3)
	if assert.Nil(t, err) {
		assertEqual(t, []int{1, 3}, ch)
		assertEqual(t, []int{1, 3}, vs.Values())
	}
}

func TestValueSet_Add_SomeToNonEmptyNoChange(t *testing.T) {
	vs := NewValueSet(1, 9)
	vs.Add(1, 3, 5, 4)
	ch, err := vs.Add(1, 3)
	if assert.Nil(t, err) {
		assertEqual(t, []int{}, ch)
		assertEqual(t, []int{1, 3, 4, 5}, vs.Values())
	}
}

func TestValueSet_Add_SomeToNonEmptyWithChange(t *testing.T) {
	vs := NewValueSet(1, 9)
	vs.Add(1, 3, 5, 4)
	ch, err := vs.Add(1, 6, 3)
	if assert.Nil(t, err) {
		assertEqual(t, []int{6}, ch)
		assertEqual(t, []int{1, 3, 4, 5, 6}, vs.Values())
	}
}

func TestValueSet_Add_JustTooSmall(t *testing.T) {
	vs := NewValueSet(1, 9)
	_, err := vs.Add(1, 6, 3, 0)
	if assert.NotNil(t, err) {
		assertEqual(t, []int{}, vs.Values())
	}
}

func TestValueSet_Add_JustTooBig(t *testing.T) {
	vs := NewValueSet(1, 9)
	_, err := vs.Add(10)
	if assert.NotNil(t, err) {
		assertEqual(t, []int{}, vs.Values())
	}
}

func TestValueSet_Remove_FromEmpty(t *testing.T) {
	vs := NewValueSet(1, 9)
	ch, err := vs.Remove(5)
	if assert.Nil(t, err) {
		assertEqual(t, []int{}, ch)
		assertEqual(t, []int{}, vs.Values())
	}
}

func TestValueSet_Remove_FromNonEmptyNoChange(t *testing.T) {
	vs := NewValueSet(1, 9)
	vs.Add(1, 3, 5, 7)
	ch, err := vs.Remove(2, 4, 6, 8)
	if assert.Nil(t, err) {
		assertEqual(t, []int{}, ch)
		assertEqual(t, []int{1, 3, 5, 7}, vs.Values())
	}
}

func TestValueSet_Remove_FromNonEmptyWithChange(t *testing.T) {
	vs := NewValueSet(1, 9)
	vs.Add(1, 3, 5, 7)
	ch, err := vs.Remove(2, 3, 4)
	if assert.Nil(t, err) {
		assertEqual(t, []int{3}, ch)
		assertEqual(t, []int{1, 5, 7}, vs.Values())
	}
}

func TestValueSet_Remove_JustTooSmall(t *testing.T) {
	vs := NewValueSet(1, 9)
	_, err := vs.Remove(0)
	if assert.NotNil(t, err) {
		assertEqual(t, []int{}, vs.Values())
	}
}

func TestValueSet_Remove_JustTooBig(t *testing.T) {
	vs := NewValueSet(1, 9)
	_, err := vs.Remove(10)
	if assert.NotNil(t, err) {
		assertEqual(t, []int{}, vs.Values())
	}
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

func assertEqual(t *testing.T, expected, actual []int) {
	if assert.Equal(t, len(expected), len(actual)) {
		for i, v := range expected {
			assert.Equal(t, v, actual[i])
		}
	}
}
