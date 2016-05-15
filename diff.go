package jsondiff

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

func Diff(a, b map[string]interface{}) []DiffItem {
	return compareStringMaps(a, b)
}

func compare(A, B interface{}) (ResolutionType, []DiffItem) {
	return TypeNotEquals, nil
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

	for kB, vB := range B {
		if _, ok := B[kB]; !ok {
			result = append(result, DiffItem{kB, nil, TypeAdded, vB})
		}
	}

	return result
}
