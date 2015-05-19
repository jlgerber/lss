package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/jlgerber/lss/pack"
	"os"
)

func getCwdPath() string {
	path, _ := os.Getwd()
	return path
}

func PrintContents(contents []string, debug bool) {
	rangePadding := 25 // this needs to be replaced

	// cast to a Stringlist and call NaturalSort()
	lss.Stringlist(contents).NaturalSort()
	//sz := len(contents)
	// convert to DirItems
	if debug {
		println("Natural Sort Contents:")
		println(lss.Stringlist(contents).String())
		println("-------------------------")
	}
	dil := lss.NewDirItemListFromSlice(contents)

	// sort DirItems into padded and nonPadded
	padded, unpadded := lss.SortDirItemList(dil)

	for x := range lss.DivideByType(unpadded) {
		println(lss.BuildRangeString(x, rangePadding))

	}
	for x := range lss.DivideByType(padded) {
		println(lss.BuildRangeString(x, rangePadding))

	}
}

func main() {

	cli.AppHelpTemplate = `NAME:
    {{.Name}} - {{.Usage}}

	USAGE:
	   {{.Name}} {{if .Flags}}[global options]{{end}} [Path]

	VERSION:
	   {{.Version}}{{if len .Authors}}

	AUTHOR(S): 
	   {{range .Authors}}{{ . }}{{end}}{{end}}

	COMMANDS:
	   {{range .Commands}}{{join .Names ", "}}{{ "\t" }}{{.Usage}}
	   {{end}}{{if .Flags}}
	GLOBAL OPTIONS:
	   {{range .Flags}}{{.}}
	   {{end}}{{end}}
`
	app := cli.NewApp()
	app.Name = "lss"
	app.Usage = `Print a directory listing, with file ranges presented in terse, 
	standard VFX form: 

	<prefix>.%[0#]d[.ext] <range>.
	
	For example, given a directory with the following contents:
	
	foo.01.exr foo.02.exr foo.03.exr
	
	print:

	foo.%02d.exr 1-3

	The user may pass an explicit directory to the command. If no directory is provided, lss uses
	the current working directory.
	`

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
		path := getCwdPath()
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

			PrintContents(contents, debug)
		}
	}

	app.Run(os.Args)
}
