package main

import (
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
)

type Export struct {
	FuncName string
	FileName string
}

func main() {

	repo, err := git.PlainOpen("~/Bobsled/bobsled")
	if err != nil {
		log.Fatal(err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("building exports map")
	exports := buildExports(worktree)

	log.Println("building imports map")
	imports := buildImports(worktree)

	for _, export := range exports {
		if _, found := imports[export.FuncName]; !found {
			log.Printf("exported function %s in file %s is not imported anywhere\n", export.FuncName, export.FileName)
		}
	}

}

func buildImports(workTree *git.Worktree) map[string]struct{} {

	allFiles := []string{}
	allImportFiles, err := workTree.Grep(&git.GrepOptions{
		Patterns: []*regexp.Regexp{
			regexp.MustCompile("import {"),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, result := range allImportFiles {
		if strings.HasSuffix(result.FileName, ".ts") {
			allFiles = append(allFiles, "/Users/aherve/Bobsled/bobsled/"+result.FileName)
		}
	}

	importsMap := make(map[string]struct{})
	for _, file := range allFiles {
		imports := findImportsInFile(file)
		for _, imp := range imports {
			importsMap[imp] = struct{}{}
		}
	}

	log.Printf("found %d unique imports across %d files\n", len(importsMap), len(allFiles))

	return importsMap
}

func findImportsInFile(filePath string) []string {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	// Read the file content
	content, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("failed to read file: %s", err)
	}

	// find imports
	re := regexp.MustCompile(`import (?:type )?\{\s*([\s\S]*?)\s*\}`)

	// Find all matches
	matches := re.FindAllStringSubmatch(string(content), -1)

	res := []string{}

	for _, match := range matches {
		names := strings.SplitSeq(match[1], ",")
		for name := range names {
			trimmed := strings.TrimSpace(name)
			if trimmed != "" {
				res = append(res, trimmed)
			}
		}
	}

	return res
}

func buildExports(workTree *git.Worktree) []Export {

	results, err := workTree.Grep(&git.GrepOptions{
		Patterns: []*regexp.Regexp{
			regexp.MustCompile("export (async )?function"),
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	res := []Export{}

	for _, result := range results {
		funcName := regexp.MustCompile(`export (async )?function (\w+)`).FindStringSubmatch(result.Content)
		if len(funcName) > 1 {
			res = append(res, Export{FuncName: funcName[2], FileName: result.FileName})
		}
	}
	log.Printf("found %d exported functions\n", len(res))
	return res
}
