package main

import (
	"context"
	"log"
	"os"

	"github.com/aherve/unused-node-exports/v2/unusedexports"
	"github.com/urfave/cli/v3"
)

func main() {

	cmd := &cli.Command{
		Name:  "unused-node-exports",
		Usage: "find unused exports in a nodejs/typescript project",
		Commands: []*cli.Command{
			{
				Name:  "scan",
				Usage: "scan git directory find unused exports",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "path",
						Aliases: []string{"p"},
						Usage:   "Path to the git repository to scan. Defaults to the current directory.",
						Value:   ".",
					},
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
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					log.Printf("scanning path %s with extensions %+v", cmd.String("path"), cmd.StringSlice("file-extensions"))
					res, err := unusedexports.FindUnusedExports(cmd.String("path"), cmd.StringSlice("file-extensions"))
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
