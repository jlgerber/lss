package lss

import (
	"strings"
	"testing"
)

func equalStringSlices(s1 []string, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

func TestSortItems_PaddedOrder(t *testing.T) {
	strs := []string{
		"this_is_v01.1000.exr",
		"this_is_v01.1001.exr",
		"this_is_v01.0008.exr",
		"this_is_v01.0009.exr",
	}

	correct := []string{
		"this_is_v01.0008.exr",
		"this_is_v01.0009.exr",
		"this_is_v01.1000.exr",
		"this_is_v01.1001.exr",
	}

	if !equalStringSlices(NaturalSort(strs), correct) {
		t.Fatal("failed to sort", strs)
	}
}

func TestSortItems_UnpaddedOrder(t *testing.T) {
	strs := []string{
		"this_is_v01.2.exr",
		"this_is_v01.1.exr",
		"this_is_v01.11.exr",
		"this_is_v01.10.exr",
	}

	correct := []string{
		"this_is_v01.1.exr",
		"this_is_v01.2.exr",
		"this_is_v01.10.exr",
		"this_is_v01.11.exr",
	}

	if !equalStringSlices(NaturalSort(strs), correct) {
		t.Fatal("failed to sort", strs)
	}
}

func TestSortItems_CaseOrder(t *testing.T) {
	strs := []string{
		"SHOUT",
		"shout",
		"SHout",
	}

	correct := []string{
		"SHOUT",
		"SHout",
		"shout",
	}

	if !equalStringSlices(NaturalSort(strs), correct) {
		t.Fatal("failed to sort", strs)
	}
}

func TestSortItems_PaddedVsUnpadded(t *testing.T) {
	strs := []string{
		"0001",
		"1",
		"2",
	}

	correct := []string{
		"0001",
		"1",
		"2",
	}

	if !equalStringSlices(NaturalSort(strs), correct) {
		t.Fatal("failed to sort", strs)
	}
}

func TestSortItems_chunking(t *testing.T) {
	strs := Stringlist{
		"this_is_v01.2.exr",
		"this_is_v01.1.exr",
		"this_is_v01.11.exr",
		"this_is_v01.10.exr",
	}
	chunked := make([]string, 0)
	for tst := range ChunkStringsToChan([]string(strs)) {
		dirname := strings.Join(tst, "")
		chunked = append(chunked, dirname)
	}

	for i, chunkval := range chunked {
		if chunkval != strs[i] {
			t.Error("improper chunking", i, chunkval)
		}
	}

}

/*

for x in range sortDirItems(genDirItems(dirItemStrings) {
	fmt.Println(x)
}

*/
