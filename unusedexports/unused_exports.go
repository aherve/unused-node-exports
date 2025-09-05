package unusedexports

import (
	"fmt"
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

type UnusedExportResult struct {
	UnusedExports   []Export
	NumberOfImports int
	NumberOfExports int
}

func FindUnusedExports(worktreePath string, fileSuffixFilter []string) (*UnusedExportResult, error) {

	repo, err := git.PlainOpen(worktreePath)
	if err != nil {
		return nil, fmt.Errorf("could not open git repo at %s: %w", worktreePath, err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("could not open git repo at %s: %w", worktreePath, err)
	}

	log.Println("building exports map")
	exports, err := buildExports(worktree, fileSuffixFilter)
	if err != nil {
		return nil, fmt.Errorf("could not build exports map: %w", err)
	}

	log.Println("building imports map")
	imports, err := buildImports(worktree, fileSuffixFilter)
	if err != nil {
		return nil, fmt.Errorf("could not build imports map: %w", err)
	}

	unusedExports := []Export{}
	for _, export := range exports {
		if _, found := imports[export.FuncName]; !found {
			unusedExports = append(unusedExports, export)
		}
	}

	return &UnusedExportResult{
		UnusedExports:   unusedExports,
		NumberOfExports: len(exports),
		NumberOfImports: len(imports),
	}, nil

}

func buildImports(workTree *git.Worktree, fileSuffixFilter []string) (map[string]struct{}, error) {
	res := make(map[string]struct{})

	// git grep all files that contain "import" key
	allImportFiles, err := workTree.Grep(&git.GrepOptions{
		Patterns: []*regexp.Regexp{
			regexp.MustCompile("import {"),
		},
	})
	if err != nil {
		return res, err
	}

	// Filter files by provided suffix
	filteredFiles := []string{}
	for _, result := range allImportFiles {
		for _, suffix := range fileSuffixFilter {
			if !strings.HasSuffix(result.FileName, suffix) {
				continue
			}
			filteredFiles = append(filteredFiles, workTree.Filesystem.Root()+"/"+result.FileName)
			break
		}
	}

	// Now parse each file for imports
	for _, file := range filteredFiles {
		imports, err := findImportsInFile(file)
		if err != nil {
			return res, fmt.Errorf("could not find imports in file %s: %w", file, err)
		}
		for _, imp := range imports {
			res[imp] = struct{}{}
		}
	}

	return res, nil
}

func findImportsInFile(filePath string) ([]string, error) {
	res := []string{}
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return res, fmt.Errorf("failed to open file: %s", err)
	}
	defer file.Close()

	// Read the file content
	content, err := io.ReadAll(file)
	if err != nil {
		return res, fmt.Errorf("failed to read file: %s", err)
	}

	// find imports
	re := regexp.MustCompile(`import (?:type )?\{\s*([\s\S]*?)\s*\}`)

	// Find all matches
	matches := re.FindAllStringSubmatch(string(content), -1)

	for _, match := range matches {
		names := strings.SplitSeq(match[1], ",")
		for name := range names {
			trimmed := strings.TrimSpace(name)
			if trimmed != "" {
				res = append(res, trimmed)
			}
		}
	}

	return res, nil
}

func buildExports(workTree *git.Worktree, fileSuffixFilter []string) ([]Export, error) {
	res := []Export{}

	results, err := workTree.Grep(&git.GrepOptions{
		Patterns: []*regexp.Regexp{
			regexp.MustCompile("export (async )?function"),
		},
	})

	if err != nil {
		return res, fmt.Errorf("could not grep for exports: %w", err)
	}

	for _, result := range results {
		for _, suffix := range fileSuffixFilter {
			if !strings.HasSuffix(result.FileName, suffix) {
				continue
			}
			funcName := regexp.MustCompile(`export (async )?function (\w+)`).FindStringSubmatch(result.Content)
			if len(funcName) > 1 {
				res = append(res, Export{FuncName: funcName[2], FileName: result.FileName})
			}
			break
		}
	}
	return res, nil
}
