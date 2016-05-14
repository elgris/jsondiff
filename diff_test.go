package diff

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
			"fieldNested1Array4": []int{"foo", "bar", "baz"},
			"fieldNested1Array5": []int{"foo", "bar", "baz"},
			"fieldNested1Array6": []int{"bar", "baz"},

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
			"fieldNested1Array4": []int{"foo", "bar", "baz"},
			"fieldNested1Array5": []int{"baz", "bar", "foo"},
			"fieldNested1Array6": []int{"foo", "baz"},

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
		{"fieldA", 12, Equals},
		{"fieldB", "foo", NotEquals, "bar"},
		{"fieldC", 34.55, Removed},
		{"fieldD", 11.22, NotEquals, 11.33},
		{"fieldE", "baz", NotEquals, 100500},
		{"fieldX", nil, Added, "bazfoo"},
		{"fieldZ", "foobar", Removed},

		{"fieldArray1", []int{1, 2, 3, 4, 5}, Equals},
		{"fieldArray2", []int{1, 2, 3, 4, 5}, NotEquals, []int{1, 2, 3, 4, 5}},
		{"fieldArray3", []int{1, 2, 3}, NotEquals, []int{1, 4, 5}},
		{"fieldArray4": []string{"foo", "bar", "baz"}, Equals},
		{"fieldArray5": []string{"foo", "bar", "baz"}, NotEquals, []string{"baz", "bar", "foo"}},
		{"fieldArray6": []string{"bar", "baz"}, NotEquals, []int{1, 2}},

		{"fieldMap1", map[string]int{"mapA": 10, "mapB": 100, "mapC": 1001}, Equals},
		{"fieldMap2", nil, Diff, []DiffItem{
			{"mapA", 10, Removed},
			{"mapB", 100, Equals},
			{"mapC", nil, Added, 100},
		}},

		{"fieldNested1", nil, Diff, []DiffItem{
			{"fieldA", 12, Equals},
			{"fieldB", "foo", NotEquals, "bar"},
			{"fieldC", 34.55, Removed},
			{"fieldD", 11.22, NotEquals, 11.33},
			{"fieldE", "baz", NotEquals, 100500},
			{"fieldX", nil, Added, "bazfoo"},
			{"fieldZ", "foobar", Removed},

			{"fieldArray1", []int{1, 2, 3, 4, 5}, Equals},
			{"fieldArray2", []int{1, 2, 3, 4, 5}, NotEquals, []int{1, 2, 3, 4, 5}},
			{"fieldArray3", []int{1, 2, 3}, NotEquals, []int{1, 4, 5}},
			{"fieldArray4": []string{"foo", "bar", "baz"}, Equals},
			{"fieldArray5": []string{"foo", "bar", "baz"}, NotEquals, []string{"baz", "bar", "foo"}},
			{"fieldArray6": []string{"bar", "baz"}, NotEquals, []int{1, 2}},

			{"fieldMap1", map[string]int{"mapA": 10, "mapB": 100, "mapC": 1001}, Equals},
			{"fieldMap2", nil, Diff, []DiffItem{
				{"mapA", 10, Removed},
				{"mapB", 100, Equals},
				{"mapC", nil, Added, 100},
			}},
		}},
	}

	actual := Diff(mapA, mapB)

	if !reflect.DeepEqual(m1, m2) {
		t.Fatalf("expected result not equal to actual")
	}
}
