package unusedexports

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
)

func Autofix(rootPath string, exports []Export) error {

	//1. group results by filename
	exportsByFile := make(map[string][]Export)
	for _, exp := range exports {
		exportsByFile[exp.FileName] = append(exportsByFile[exp.FileName], exp)
	}

	// 2. for each file, sort exports by line number, ascending
	var wg sync.WaitGroup
	for _, exps := range exportsByFile {
		wg.Add(1)
		go func(exps []Export) {
			defer wg.Done()
			sort.Slice(exps, func(i, j int) bool {
				return exps[i].LineNumber < exps[j].LineNumber
			})
		}(exps)
	}
	wg.Wait()

	// Process each file
	for filePath, exps := range exportsByFile {
		if err := fixFileExports(rootPath+"/"+filePath, exps); err != nil {
			return fmt.Errorf("could not fix file %s: %w", filePath, err)
		}
	}
	log.Printf("autofixed %d exports amongst %d files", len(exports), len(exportsByFile))
	return nil
}

func fixFileExports(filePath string, exports []Export) error {
	changeRegex := regexp.MustCompile(removeExportRegexPattern)

	log.Println("Fixing file", filePath)

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("could not open file %s: %w", filePath, err)
	}

	newFileLines := []string{}
	scanner := bufio.NewScanner(file)
	lineNumber := 0

	nextChangeIndex := 0
	var nextChange *Export = &exports[nextChangeIndex]

	for scanner.Scan() {
		line := scanner.Text()
		lineNumber++

		if nextChange == nil {
			newFileLines = append(newFileLines, line)
			continue
		}

		if lineNumber == nextChange.LineNumber {

			changed := changeRegex.ReplaceAllString(line, "$1")
			newFileLines = append(newFileLines, changed)

			nextChangeIndex++
			if nextChangeIndex < len(exports) {
				nextChange = &exports[nextChangeIndex]
			} else {
				nextChange = nil
			}
		} else {
			newFileLines = append(newFileLines, line)
		}
	}
	file.Close()

	// Now replace the file content
	newContent := strings.Join(newFileLines, "\n") + "\n"
	return os.WriteFile(filePath, []byte(newContent), 0644)
}
