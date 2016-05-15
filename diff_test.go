package jsondiff

import (
	"reflect"
	"testing"
)

// func TestA(t *testing.T) {
// 	a := map[string]interface{}{
// 		"test1": 12,
// 		"test3": 23,
// 	}
// 	b := map[string]interface{}{
// 		"test1": 12,
// 		"test2": 12,
// 	}

// 	diffA := NewDiffMap(a)
// 	diffB := NewDiffMap(b)

// for _, k := range diffA.Keys() {
// 	vA, _ := diffA.Get(k)
// 	vA.Mark(k)

// 	vB, ok := diffB.Get(k)
// 	if !ok {
// 		diffA.Mark(k, GREEN)
// 		continue
// 	}

// 	if !compare(vA, vB) {
// 		diffA.Mark(k, YELLOW)
// 		diffB.Mark(k, YELLOW)
// 	}
// }

// for _, k := range diffB.Unmarked() {
// 	diffB.Mark(k, GREEN)
// }

// }

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

	expectedResult := []DiffItem{
		{"fieldA", 12, TypeEquals, nil},
		{"fieldB", "foo", TypeNotEquals, "bar"},
		{"fieldC", 34.55, TypeRemoved, nil},
		{"fieldD", 11.22, TypeNotEquals, 11.33},
		{"fieldE", "baz", TypeNotEquals, 100500},
		{"fieldX", nil, TypeAdded, "bazfoo"},
		{"fieldZ", "foobar", TypeRemoved, nil},

		{"fieldArray1", []int{1, 2, 3, 4, 5}, TypeEquals, nil},
		{"fieldArray2", []int{1, 2, 3, 4, 5}, TypeNotEquals, []int{1, 2, 3, 4, 5}},
		{"fieldArray3", []int{1, 2, 3}, TypeNotEquals, []int{1, 4, 5}},
		{"fieldArray4", []string{"foo", "bar", "baz"}, TypeEquals, nil},
		{"fieldArray5", []string{"foo", "bar", "baz"}, TypeNotEquals, []string{"baz", "bar", "foo"}},
		{"fieldArray6", []string{"bar", "baz"}, TypeNotEquals, []int{1, 2}},

		{"fieldMap1", map[string]int{"mapA": 10, "mapB": 100, "mapC": 1001}, TypeEquals, nil},
		{"fieldMap2", nil, TypeDiff, []DiffItem{
			{"mapA", 10, TypeRemoved, nil},
			{"mapB", 100, TypeEquals, nil},
			{"mapC", nil, TypeAdded, 100},
		}},

		{"fieldNested1", nil, TypeDiff, []DiffItem{
			{"fieldA", 12, TypeEquals, nil},
			{"fieldB", "foo", TypeNotEquals, "bar"},
			{"fieldC", 34.55, TypeRemoved, nil},
			{"fieldD", 11.22, TypeNotEquals, 11.33},
			{"fieldE", "baz", TypeNotEquals, 100500},
			{"fieldX", nil, TypeAdded, "bazfoo"},
			{"fieldZ", "foobar", TypeRemoved, nil},

			{"fieldArray1", []int{1, 2, 3, 4, 5}, TypeEquals, nil},
			{"fieldArray2", []int{1, 2, 3, 4, 5}, TypeNotEquals, []int{1, 2, 3, 4, 5}},
			{"fieldArray3", []int{1, 2, 3}, TypeNotEquals, []int{1, 4, 5}},
			{"fieldArray4", []string{"foo", "bar", "baz"}, TypeEquals, nil},
			{"fieldArray5", []string{"foo", "bar", "baz"}, TypeNotEquals, []string{"baz", "bar", "foo"}},
			{"fieldArray6", []string{"bar", "baz"}, TypeNotEquals, []int{1, 2}},

			{"fieldMap1", map[string]int{"mapA": 10, "mapB": 100, "mapC": 1001}, TypeEquals, nil},
			{"fieldMap2", nil, TypeDiff, []DiffItem{
				{"mapA", 10, TypeRemoved, nil},
				{"mapB", 100, TypeEquals, nil},
				{"mapC", nil, TypeAdded, 100},
			}},
		}},
	}

	actual := Diff(mapA, mapB)

	if !reflect.DeepEqual(expectedResult, actual) {
		t.Fatalf("expected result not equal to actual")
	}
}
