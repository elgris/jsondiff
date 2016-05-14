package diff

import "sort"

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

func NewDiffMap(m map[string]interface{}) *diffMap {
	mapLen := len(m)

	keys := make([]string, 0, mapLen)
	marked := make([]string, 0, mapLen)
	unmarked := make(map[string]interface{}, 0, mapLen)

	for k, _ := range m {
		keys := append(keys, k)
		unmarked[k] = nil
	}
	sort.Strings(keys)

	return &diffMap{
		data:     m,
		keys:     keys,
		marked:   marked,
		unmarked: unmarked,
	}
}

func Diff(a, b map[string]interface{}) []DiffItem {
	return compareStringMaps(a, b)
}

type diffMap struct {
	data     map[string]interface{}
	keys     []string
	marked   []string
	unmarked map[string]interface{}
}

func (m *diffMap) Mark(key string) {
	if _, ok := m.data[key]; ok {
		m.marked[key] = nil
		delete(m.unmarked, key)
	}
}

func (m diffMap) Keys() []string { return m.keys }

func (m diffMap) Marked() []string {
	marked := make([]string, 0, len(m.marked))
	for _, k := range m.keys {
		if _, ok := m.marked[k]; ok {
			marked = append(marked, k)
		}
	}

	return marked
}

func (m diffMap) Unmarked() []string {
	unmarked := make([]string, 0, len(m.unmarked))
	for _, k := range m.keys {
		if _, ok := m.unmarked[k]; ok {
			unmarked = append(unmarked, k)
		}
	}

	return unmarked
}

func compare(A, B interface{}) ResolutionType {
	/*
				Cases:
		        - A,B - simple types (int, double, string)
		        - A,B - array types (array, slice)
		        - A,B - map[string]* types
		        - The rest of cases

	*/

}

func compareStringMaps(A, B map[string]interface{}) []DiffItem {
	result := []DiffItem{}
	for kA, vA := range A {

		vB, ok := B[kA]
		if !ok {
			result = append(result, DiffItem{kA, vA, Removed})
			continue
		}

		resolutionType, diffs := compare(vA, vB)

		switch resolutionType {
		case Equals:
			result = append(result, DiffItem{kA, vA, Equals})
		case NotEquals:
			result = append(result, DiffItem{kA, vA, NotEquals, vB})
		case Diff:
			result = append(result, DiffItem{kA, nil, Diff, diffs})
		}
	}

	for kB, vB := range B {
		if _, ok := B[kB]; !ok {
			result = append(result, DiffItem{kB, nil, Added, vB})
		}
	}
}
