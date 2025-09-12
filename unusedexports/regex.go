package unusedexports

const hasExportRegexPattern = `export (?:async )?(?:function|const)`
const exportNameRegexPattern = `export (?:async )?(?:function|const) (\w+)`

const removeExportRegexPattern = `export (.*)`
