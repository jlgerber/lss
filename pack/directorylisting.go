package lss

import (
	"errors"
	"os"
)

// NewDirItemListFromPath returns an error object and a DirItemList
func FilteredListingFromPath(path string, filter func(string) bool) (error, []string) {
	vsf := make([]string, 0)
	dir, err := os.Open(path)
	if err == nil {

		defer func() {
			if err := dir.Close(); err != nil {
				panic(err)
			}
		}()

		fileInfo, err := dir.Stat()
		if err != nil {
			return err, vsf
		}

		if !fileInfo.IsDir() {
			return errors.New("Supplied path:'" + path + "' is not a directory"), vsf
		}

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

func GetCwdPath() string {
	path, _ := os.Getwd()
	return path
}
