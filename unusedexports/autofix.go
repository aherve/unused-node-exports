package unusedexports

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func Autofix(exports []Export) error {

	//1. group results by filename
	exportsByFile := make(map[string][]Export)
	for _, exp := range exports {
		exportsByFile[exp.FileName] = append(exportsByFile[exp.FileName], exp)
	}

	//2. for each file, sort exports by line number ascending
	for _, exps := range exportsByFile {
		sort.Slice(exps, func(i, j int) bool {
			return exps[i].LineNumber < exps[j].LineNumber
		})
	}

	// Process each file
	for filePath, exps := range exportsByFile {
		if err := fixFileExports(filePath, exps); err != nil {
			return fmt.Errorf("could not fix file %s: %w", filePath, err)
		}
	}
}

func fixFileExports(filePath string, exports []Export) error {
	// open the file for reading
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("could not open file %s: %w", filePath, err)
	}
	defer file.Close()

	newFileLines := []string{}
	scanner := bufio.NewScanner(file)
	lineNumber := 1
	nextChangeIndex := 0

	for scanner.Scan() {
		line := scanner.Text()
	}

	return nil
}
