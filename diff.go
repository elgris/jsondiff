package jsondiff

import (
	"encoding/json"
	"reflect"
	"sort"
)

/*
   - sort keys
   - for each key from A check if the key is existing in B:
      - if existing, then compare A[k] and B[k]:
          - if they are basic type (comparable), then just compare
          - if they are different types, then mark them as "non-equal" (yellow) immediately
          - if they are arrays - compare them according arrays comparison (position-wise)
          - if they are maps - run this algorithm recursively
      - if B[k] is not existing - mark it as "green" in A
      - mark k as "checked"
   - for each key k2 from B:
       - if k2 is "checked" - skip
       - otherwise - mark it as "green" in B
*/

type ResolutionType int

const (
	TypeEquals ResolutionType = iota
	TypeNotEquals
	TypeAdded
	TypeRemoved
	TypeDiff
)

type DiffItem struct {
	Key        string
	ValueA     interface{}
	Resolution ResolutionType
	ValueB     interface{}
}

type byKey []DiffItem

func (m byKey) Len() int           { return len(m) }
func (m byKey) Less(i, j int) bool { return m[i].Key < m[j].Key }
func (m byKey) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }

func Diff(a, b interface{}) []DiffItem {
	mapA := map[string]interface{}{}
	mapB := map[string]interface{}{}

	jsonA, _ := json.Marshal(a)
	jsonB, _ := json.Marshal(b)

	json.Unmarshal(jsonA, &mapA)
	json.Unmarshal(jsonB, &mapB)

	return compareStringMaps(mapA, mapB)
}

func compare(A, B interface{}) (ResolutionType, []DiffItem) {
	equals := reflect.DeepEqual(A, B)
	if equals {
		return TypeEquals, nil
	}

	mapA, okA := A.(map[string]interface{})
	mapB, okB := B.(map[string]interface{})

	if okA && okB {
		diff := compareStringMaps(mapA, mapB)
		return TypeDiff, diff
	}

	return TypeNotEquals, nil
}

func compareStringMaps(A, B map[string]interface{}) []DiffItem {
	keysA := sortedKeys(A)
	keysB := sortedKeys(B)

	result := []DiffItem{}

	for _, kA := range keysA {
		vA := A[kA]

		vB, ok := B[kA]
		if !ok {
			result = append(result, DiffItem{kA, vA, TypeRemoved, nil})
			continue
		}

		resolutionType, diffs := compare(vA, vB)

		switch resolutionType {
		case TypeEquals:
			result = append(result, DiffItem{kA, vA, TypeEquals, nil})
		case TypeNotEquals:
			result = append(result, DiffItem{kA, vA, TypeNotEquals, vB})
		case TypeDiff:
			result = append(result, DiffItem{kA, nil, TypeDiff, diffs})
		}
	}

	for _, kB := range keysB {
		if _, ok := A[kB]; !ok {
			result = append(result, DiffItem{kB, nil, TypeAdded, B[kB]})
		}
	}

	sort.Sort(byKey(result))

	return result
}

func sortedKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k, _ := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}
