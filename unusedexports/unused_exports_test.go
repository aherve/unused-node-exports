package unusedexports

import (
	"regexp"
	"testing"
)

var tests = []struct {
	Name       string
	Input      string
	ExportName string
}{
	{
		Name:       "named export function",
		Input:      "export function myFunction() {}",
		ExportName: "myFunction",
	},
	{
		Name:       "async function",
		Input:      "export async function myFunction() {}",
		ExportName: "myFunction",
	},
	{
		Name:       "exported const arrow function",
		Input:      "export const myFunction = () => {};",
		ExportName: "myFunction",
	},
	{
		Name:       "exported const",
		Input:      "export const myConst = 42;",
		ExportName: "myConst",
	},
	{
		Name:       "not an export (function)",
		Input:      "function myFunction() {}",
		ExportName: "",
	},
	{
		Name:       "not an export (const)",
		Input:      "const myConst = () => {}",
		ExportName: "",
	},
}

func TestExportRegexes(t *testing.T) {

	hasExport, err := regexp.Compile(hasExportRegexPattern)
	if err != nil {
		t.Errorf("could not compile hasExport regex: %v", err)
	}

	exportName, err := regexp.Compile(exportNameRegexPattern)
	if err != nil {
		t.Errorf("could not compile exportName regex: %v", err)
	}

	for _, test := range tests {

		hasExportRes := hasExport.MatchString(test.Input)
		shouldHaveExport := test.ExportName != ""
		if hasExportRes != shouldHaveExport {
			t.Errorf("test %s: expected hasExport to be %v, got %v", test.Name, shouldHaveExport, hasExportRes)
		}

		if !shouldHaveExport {
			continue
		}

		exportNameRes := exportName.FindStringSubmatch(test.Input)
		if len(exportNameRes) < 2 {
			t.Errorf("test %s: expected exportName to find a match, got none", test.Name)
			continue
		}
		if exportNameRes[1] != test.ExportName {
			t.Errorf("test %s: expected exportName to be %s, got %s", test.Name, test.ExportName, exportNameRes[2])
		}

	}
}

func TestFindImportsInContent(t *testing.T) {
	input := `
	import { A, B,
C } from 'module1';
	import { D } from 'module2';
	import type { E
	,
	F } from 'module3';
	import { G as aliasG,    H  } from 'module4';
	
	import type { I as aliasI, J } from 'module5';
	`

	found, err := findImportsInContent([]byte(input))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	expected := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}
	if len(found) != len(expected) {
		t.Errorf("expected %d imports, got %d", len(expected), len(found))
	}
	for i, exp := range expected {
		if found[i] != exp {
			t.Errorf("expected import %s, got %s", exp, found[i])
		}
	}
}
