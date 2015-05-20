package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/jlgerber/lss/pack"
	"os"
)

func PrintContents(contents []string, debug bool) {

	// cast to a Stringlist and call NaturalSort()
	lss.Stringlist(contents).NaturalSort()

	if debug {
		println("Natural Sort Contents:")
		println(lss.Stringlist(contents).String())
		println("-------------------------")
	}

	// build a dirItemList from the contents slice
	dil := lss.NewDirItemListFromSlice(contents)

	for item := range lss.RangeStringsFromSortedItemList(dil) {
		fmt.Println(item)
	}

}

func main() {

	cli.AppHelpTemplate = lss.AppHelpTemplate
	app := cli.NewApp()
	app.Name = "lss"
	app.Usage = lss.Usage

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name: "debug, d",
			//Value: false,
			Usage: "turn debugging on.",
		},
		cli.BoolFlag{
			Name:  "all, a",
			Usage: "show hidden files.",
		},
	}

	app.Action = func(c *cli.Context) {
		args := c.Args()
		path := lss.GetCwdPath()
		if len(args) > 0 {
			path = args[0]
		}

		debug := c.Bool("debug")

		if debug {
			println("Path:", path)
			println("-----------")
		}

		showHidden := func(nm string) bool {
			if c.Bool("all") == true {
				return true
			}
			if string(nm[0]) == "." {
				return false
			}
			return true
		}
		// unsorted path contents
		err, contents := lss.FilteredListingFromPath(path, showHidden)
		if err != nil {
			fmt.Println(err)
		} else {
			//			PrintContents(contents, debug)

			for value := range lss.RangesChanFromStringSlice(contents) {
				fmt.Println(value)
			}
		}
	}

	app.Run(os.Args)
}
