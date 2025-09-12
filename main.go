package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aherve/unused-node-exports/v2/unusedexports"
	"github.com/urfave/cli/v3"
)

const version = "v3.0.0"

func showVersion() {
	fmt.Println(version)
}

func main() {

	pathArg := "."

	cmd := &cli.Command{
		Name:  "unused-node-exports",
		Usage: "find unused exports in a nodejs/typescript project",
		Commands: []*cli.Command{
			{
				Name:  "version",
				Usage: "show version",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					showVersion()
					return nil
				},
			},
			{
				Name:  "scan",
				Usage: "scan git directory find unused exports",
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name:        "path",
						Value:       ".",
						Destination: &pathArg,
					},
				},
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:    "file-extensions",
						Aliases: []string{"e"},
						Usage:   "List of file extensions to consider. If provided, only files with these extensions will be scanned.",
						Value:   []string{".ts", ".tsx", ".js", ".jsx", ".mjs", ".cjs", ".mts", ".cts"},
					},
					&cli.StringFlag{
						Name:    "output",
						Aliases: []string{"o"},
						Usage:   "If provided, the results will be written to this file in CSV format",
						Value:   "",
					},
					&cli.StringFlag{
						Name:    "export-prefix",
						Aliases: []string{"p", "prefix"},
						Usage:   "If provided, only exports starting with this prefix will be considered. This is useful to find unused exports in a specific namespace.",
						Value:   "",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					log.Printf("scanning path %s with extensions %+v", pathArg, cmd.StringSlice("file-extensions"))
					res, err := unusedexports.FindUnusedExports(pathArg, cmd.StringSlice("file-extensions"), cmd.String("export-prefix"))
					if err != nil {
						return err
					}

					if outFile := cmd.String("output"); outFile != "" {
						return CSVExport(res.UnusedExports, outFile)
					}

					for _, exp := range res.UnusedExports {
						log.Printf("%s\t%s", exp.FileName, exp.ExportName)
					}

					log.Printf("found %d unused exports amongst %d imports and %d exports", len(res.UnusedExports), res.NumberOfImports, res.NumberOfExports)
					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}

}

func CSVExport(exports []unusedexports.Export, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	file.WriteString("file,exportName\n")
	for _, entry := range exports {
		if _, err := file.WriteString(entry.FileName + "," + entry.ExportName + "\n"); err != nil {
			return err
		}
	}
	return nil
}
