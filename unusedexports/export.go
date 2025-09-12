package unusedexports

import "fmt"

type Export struct {
	ExportName string
	FileName   string
	LineNumber int
}

// provide a string interface for Export
func (e Export) String() string {
	return fmt.Sprintf("%s:%d\t%s", e.FileName, e.LineNumber, e.ExportName)
}
