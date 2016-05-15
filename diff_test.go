package jsondiff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiff(t *testing.T) {
	mapA := map[string]interface{}{
		"fieldA": 12,
		"fieldB": "foo",
		"fieldC": 34.55,
		"fieldD": 11.22,
		"fieldE": "baz",
		"fieldZ": "foobar",

		"fieldArray1": []int{1, 2, 3, 4, 5},
		"fieldArray2": []int{1, 4, 2, 3, 5},
		"fieldArray3": []int{1, 2, 3},
		"fieldArray4": []string{"foo", "bar", "baz"},
		"fieldArray5": []string{"foo", "bar", "baz"},
		"fieldArray6": []string{"bar", "baz"},

		"fieldMap1": map[string]int{
			"mapA": 10,
			"mapC": 1001,
			"mapB": 100,
		},
		"fieldMap2": map[string]int{
			"mapA": 10,
			"mapB": 100,
		},

		"fieldNested1": map[string]interface{}{
			"fieldNested1A": 12,
			"fieldNested1B": "foo",
			"fieldNested1C": 34.55,
			"fieldNested1D": 11.22,
			"fieldNested1E": "baz",
			"fieldNested1Z": "foobar",

			"fieldNested1Array1": []int{1, 2, 3, 4, 5},
			"fieldNested1Array2": []int{1, 4, 2, 3, 5},
			"fieldNested1Array3": []int{1, 2, 3},
			"fieldNested1Array4": []string{"foo", "bar", "baz"},
			"fieldNested1Array5": []string{"foo", "bar", "baz"},
			"fieldNested1Array6": []string{"bar", "baz"},

			"fieldNested1Map1": map[string]int{
				"mapA": 10,
				"mapC": 1001,
				"mapB": 100,
			},
			"fieldNested1Map2": map[string]int{
				"mapA": 10,
				"mapB": 100,
			},
		},
	}

	mapB := map[string]interface{}{
		"fieldA": 12,
		"fieldB": "bar",
		"fieldD": 11.33,
		"fieldE": 100500,
		"fieldX": "bazfoo",

		"fieldArray1": []int{1, 2, 3, 4, 5},
		"fieldArray2": []int{1, 2, 3, 4, 5},
		"fieldArray3": []int{1, 4, 5},
		"fieldArray4": []string{"foo", "bar", "baz"},
		"fieldArray5": []string{"baz", "bar", "foo"},
		"fieldArray6": []int{1, 2},

		"fieldMap1": map[string]int{
			"mapB": 100,
			"mapA": 10,
			"mapC": 1001,
		},
		"fieldMap2": map[string]int{
			"mapB": 100,
			"mapC": 100,
		},

		"fieldNested1": map[string]interface{}{
			"fieldNested1A": 12,
			"fieldNested1B": "bar",
			"fieldNested1D": 11.33,
			"fieldNested1E": 100500,
			"fieldNested1X": "bazfoo",

			"fieldNested1Array1": []int{1, 2, 3, 4, 5},
			"fieldNested1Array2": []int{1, 2, 3, 4, 5},
			"fieldNested1Array3": []int{1, 4, 5},
			"fieldNested1Array4": []string{"foo", "bar", "baz"},
			"fieldNested1Array5": []string{"baz", "bar", "foo"},
			"fieldNested1Array6": []string{"foo", "baz"},

			"fieldNested1Map1": map[string]int{
				"mapB": 100,
				"mapA": 10,
				"mapC": 1001,
			},
			"fieldNested1Map2": map[string]int{
				"mapB": 100,
				"mapC": 100,
			},
		},
	}

	// Limitations and gotchas:
	// 1. all numeric types converted to interface{}
	// 2. all slice types converted to []interface{}
	// 3. all map types converted to map[string]interface{}
	// 4. at the end keys are sorted alphabetically
	expected := []DiffItem{
		{"fieldA", 12., TypeEquals, nil},

		{"fieldArray1", []interface{}{1., 2., 3., 4., 5.}, TypeEquals, nil},
		{"fieldArray2", []interface{}{1., 4., 2., 3., 5.}, TypeNotEquals, []interface{}{1., 2., 3., 4., 5.}},
		{"fieldArray3", []interface{}{1., 2., 3.}, TypeNotEquals, []interface{}{1., 4., 5.}},
		{"fieldArray4", []interface{}{"foo", "bar", "baz"}, TypeEquals, nil},
		{"fieldArray5", []interface{}{"foo", "bar", "baz"}, TypeNotEquals, []interface{}{"baz", "bar", "foo"}},
		{"fieldArray6", []interface{}{"bar", "baz"}, TypeNotEquals, []interface{}{1., 2.}},

		{"fieldB", "foo", TypeNotEquals, "bar"},
		{"fieldC", 34.55, TypeRemoved, nil},
		{"fieldD", 11.22, TypeNotEquals, 11.33},
		{"fieldE", "baz", TypeNotEquals, 100500.},

		{"fieldMap1", map[string]interface{}{"mapA": 10., "mapB": 100., "mapC": 1001.}, TypeEquals, nil},
		{"fieldMap2", nil, TypeDiff, []DiffItem{
			{"mapA", 10., TypeRemoved, nil},
			{"mapB", 100., TypeEquals, nil},
			{"mapC", nil, TypeAdded, 100.},
		}},

		{"fieldNested1", nil, TypeDiff, []DiffItem{
			{"fieldNested1A", 12., TypeEquals, nil},

			{"fieldNested1Array1", []interface{}{1., 2., 3., 4., 5.}, TypeEquals, nil},
			{"fieldNested1Array2", []interface{}{1., 4., 2., 3., 5.}, TypeNotEquals, []interface{}{1., 2., 3., 4., 5.}},
			{"fieldNested1Array3", []interface{}{1., 2., 3.}, TypeNotEquals, []interface{}{1., 4., 5.}},
			{"fieldNested1Array4", []interface{}{"foo", "bar", "baz"}, TypeEquals, nil},
			{"fieldNested1Array5", []interface{}{"foo", "bar", "baz"}, TypeNotEquals, []interface{}{"baz", "bar", "foo"}},
			{"fieldNested1Array6", []interface{}{"bar", "baz"}, TypeNotEquals, []interface{}{"foo", "baz"}},

			{"fieldNested1B", "foo", TypeNotEquals, "bar"},
			{"fieldNested1C", 34.55, TypeRemoved, nil},
			{"fieldNested1D", 11.22, TypeNotEquals, 11.33},
			{"fieldNested1E", "baz", TypeNotEquals, 100500.},

			{"fieldNested1Map1", map[string]interface{}{"mapA": 10., "mapB": 100., "mapC": 1001.}, TypeEquals, nil},
			{"fieldNested1Map2", nil, TypeDiff, []DiffItem{
				{"mapA", 10., TypeRemoved, nil},
				{"mapB", 100., TypeEquals, nil},
				{"mapC", nil, TypeAdded, 100.},
			}},

			{"fieldNested1X", nil, TypeAdded, "bazfoo"},
			{"fieldNested1Z", "foobar", TypeRemoved, nil},
		}},

		{"fieldX", nil, TypeAdded, "bazfoo"},
		{"fieldZ", "foobar", TypeRemoved, nil},
	}

	actual := Diff(mapA, mapB)

	assert.Equal(t, expected, actual)
}
