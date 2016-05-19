package jsondiff

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"sort"

	"github.com/mgutz/ansi"
)

type ResolutionType int

const (
	TypeEquals ResolutionType = iota
	TypeNotEquals
	TypeAdded
	TypeRemoved
	TypeDiff
)

var (
	colorStartYellow = ansi.ColorCode("yellow")
	colorStartRed    = ansi.ColorCode("red")
	colorStartGreen  = ansi.ColorCode("green")
	colorReset       = ansi.ColorCode("reset")
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

func Format(items []DiffItem) []byte {
	buf := bytes.Buffer{}

	writeItems(&buf, "", items)

	return buf.Bytes()
}

func writeItems(writer io.Writer, prefix string, items []DiffItem) {
	writer.Write([]byte{'{'})
	last := len(items) - 1

	prefixNotEqualsA := prefix + "<> "
	prefixNotEqualsB := prefix + "** "
	prefixAdded := prefix + "<< "
	prefixRemoved := prefix + ">> "

	for i, item := range items {
		writer.Write([]byte{'\n'})

		switch item.Resolution {
		case TypeEquals:
			writeItem(writer, prefix, item.Key, item.ValueA, i < last)
		case TypeNotEquals:
			writer.Write([]byte(colorStartYellow))

			writeItem(writer, prefixNotEqualsA, item.Key, item.ValueA, i < last)
			writer.Write([]byte{'\n'})
			writeItem(writer, prefixNotEqualsB, item.Key, item.ValueB, i < last)

			writer.Write([]byte(colorReset))
		case TypeAdded:
			writer.Write([]byte(colorStartGreen))
			writeItem(writer, prefixAdded, item.Key, item.ValueB, i < last)
			writer.Write([]byte(colorReset))
		case TypeRemoved:
			writer.Write([]byte(colorStartRed))
			writeItem(writer, prefixRemoved, item.Key, item.ValueA, i < last)
			writer.Write([]byte(colorReset))
		case TypeDiff:
			subdiff := item.ValueB.([]DiffItem)
			fmt.Fprintf(writer, "%s\"%s\": ", prefix, item.Key)
			writeItems(writer, prefix+"    ", subdiff)
			if i < last {
				writer.Write([]byte{','})
			}
		}

	}

	fmt.Fprintf(writer, "\n%s}", prefix)
}

func writeItem(writer io.Writer, prefix, key string, value interface{}, isNotLast bool) {
	fmt.Fprintf(writer, "%s\"%s\": ", prefix, key)
	serialized, _ := json.Marshal(value)
	writer.Write(serialized)
	if isNotLast {
		writer.Write([]byte{','})
	}
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
