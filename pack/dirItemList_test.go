package lss

import (
	"testing"
)

func TestDirItemList_NewFromSlice(t *testing.T) {
	il := NewDirItemListFromSlice([]string{
		"foo.0001.mb",
		"foo.0002.mb",
	})
	if len(il) != 2 {
		t.Error("Wrong number of items constructed:", len(il), "Should Be:2. Items:", il)
	}
}

func testEq(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func TestDirItemList_SortDirItemList(t *testing.T) {

	tm := map[string][]string{
		"t1": []string{
			"foo.2.mb",
			"foo.4.mb",
			"foo.10.mb",
			"foo.0001.mb",
			"foo.0002.mb",
		},

		"t1r": []string{
			"foo.0001.mb",
			"foo.0002.mb",
			"foo.2.mb",
			"foo.4.mb",
			"foo.10.mb",
		},
		"t2": []string{
			"foo.2.mb",
			"foo.4.mb",
			"foo.10.mb",
			"foo.9.mb",
			"foo.0100.mb",
			"foo.0101.mb",
			"foo.0102.mb",
		},
		"t2r": []string{
			"foo.0100.mb",
			"foo.0101.mb",
			"foo.0102.mb",
			"foo.2.mb",
			"foo.4.mb",
			"foo.9.mb",
			"foo.10.mb",
		},
	}

	test := func(contents []string, results []string) {
		sz := len(contents)
		Stringlist(contents).NaturalSort()

		il := NewDirItemListFromSlice(contents)
		if len(il) != sz {
			t.Error("Wrong number of items constructed:", len(il),
				"Should Be:", sz, ".Number of Items:", il)
		}

		padded, unpadded := SortDirItemList(il)
		pdil := NewSliceFromDirItemList(padded)
		updil := NewSliceFromDirItemList(unpadded)
		dil := append(pdil, updil...)

		if !testEq(dil, results) {
			t.Error("results:", dil, "Does not equal expected results:", results)
			print(padded.String())
			println("--------")
			println(unpadded.String())
		}
	}

	test(tm["t1"], tm["t1r"])
	test(tm["t2"], tm["t2r"])

}

func TestDirItemList_BuildRangeString(t *testing.T) {
	test := []string{
		"foo.0002.mb",
		"foo.0004.mb",
		"foo.0003.mb",
		"foo.0010.mb",
		"foo.0009.mb",
		"foo.0100.mb",
		"foo.0101.mb",
		"foo.0105.mb",
		"foo.0102.mb",
	}
	sz := len(test)
	Stringlist(test).NaturalSort()
	il := NewDirItemListFromSlice(test)
	if len(il) != sz {
		t.Error("Wrong number of items constructed:", len(il),
			"Should Be:", sz, ".Number of Items:", il)
	}

	rs := BuildRangeString(il)
	println(rs)
}

func TestDirItemList_DivideByType(t *testing.T) {
	test := []string{
		"foo.0002.mb",
		"foo.0004.mb",
		"foo.0003.mb",
		"foobar.10.mb",
		"foobar.9.mb",
		"foobar.0100.mb",
		"foobar.0101.mb",
		"foobar.0105.mb",
		"foobarba.0103.mb",
		"foobarba.0104.mb",
		"foobarba.0102.mb",
	}
	sz := len(test)
	Stringlist(test).NaturalSort()
	il := NewDirItemListFromSlice(test)
	if len(il) != sz {
		t.Error("Wrong number of items constructed:", len(il),
			"Should Be:", sz, ".Number of Items:", il)
	}

	for x := range DivideByType(il) {
		println(BuildRangeString(x))
		/*for _, y := range x {
			println(y.String())
		}*/
		println("")
	}
}
