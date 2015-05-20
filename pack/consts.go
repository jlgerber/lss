package lss

var AppHelpTemplate string = `NAME:
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

var Usage = `Print a directory listing, with file ranges presented in terse, 
	standard VFX form: 

	<prefix>.%[0#]d[.ext] <range>.
	
	For example, given a directory with the following contents:
	
	foo.01.exr foo.02.exr foo.03.exr
	
	print:

	foo.%02d.exr 1-3

	The user may pass an explicit directory to the command. If no directory is provided, lss uses
	the current working directory.
	`
