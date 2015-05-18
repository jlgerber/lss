package main

import (
	//"fmt"
	"github.com/jlgerber/lss/pack"
	"os"
)

func getPath() string {
	var path string
	if len(os.Args) > 1 {
		path = os.Args[1]
	} else {
		path, _ = os.Getwd()
	}
	return path
}

func main() {
	path := getPath()
	println("Path:", path)

	// unsorted path contents
	err, contents := lss.FilteredListingFromPath(path, nil)
	if err != nil {
		println("shit")
	}

	// cast to a Stringlist and call NaturalSort()
	lss.Stringlist(contents).NaturalSort()
	//sz := len(contents)
	// convert to DirItems
	println("CONTENTS")
	println(lss.Stringlist(contents).String())
	println("DONE")
	dil := lss.NewDirItemListFromSlice(contents)
	/*if len(dil) != sz {
		t.Error("Wrong number of items constructed:", len(dil),
			"Should Be:", sz, ".Number of Items:", dil)
	}*/

	// sort DirItems into padded and nonPadded
	padded, unpadded := lss.SortDirItemList(dil)

	for x := range lss.DivideByType(padded) {
		println(lss.BuildRangeString(x))

	}
	for x := range lss.DivideByType(unpadded) {
		println(lss.BuildRangeString(x))

	}
}
