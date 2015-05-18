package lss

import (
	"os"
)

// NewDirItemListFromPath returns an error object and a DirItemList
func FilteredListingFromPath(path string, filter func(string) bool) (error, []string) {
	vsf := make([]string, 0)
	dir, err := os.Open(path)
	if err == nil {
		filenames, err := dir.Readdirnames(0)
		if err == nil {

			if filter == nil {
				filter = func(nm string) bool { return true }
			}
			for _, v := range filenames {
				if filter(v) {
					vsf = append(vsf, v)
				}
			}
		}

	}
	return err, vsf

}
