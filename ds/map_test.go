package ds

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MapFilterFunc_1(t *testing.T) {
	result := MapFilterFunc(map[int]bool{1: true, 2: false, 3: true}, func(v bool, _ int) bool {
		return v == true
	})
	assert.Equal(t, map[int]bool{1: true, 3: true}, result)
}

func Test_MapFilterFunc_2(t *testing.T) {
	result := MapFilterFunc(map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}, func(v int, _ string) bool {
		return v <= 3
	})
	assert.Equal(t, map[string]int{"a": 1, "b": 2, "c": 3}, result)
}

func Test_MapKeys_1(t *testing.T) {
	result := MapKeys(map[string]string{"a": "1", "b": "2"})
	assert.ElementsMatch(t, []string{"a", "b"}, result)
}

func Test_MapKeys_2(t *testing.T) {
	result := MapKeys(map[int]bool{1: true, 2: false})
	assert.ElementsMatch(t, []int{1, 2}, result)
}

func Test_MapValues_1(t *testing.T) {
	result := MapValues(map[string]string{"a": "1", "b": "2"})
	assert.ElementsMatch(t, []string{"1", "2"}, result)
}

func Test_MapValues_2(t *testing.T) {
	result := MapValues(map[int]bool{1: true, 2: false})
	assert.ElementsMatch(t, []bool{true, false}, result)
}
