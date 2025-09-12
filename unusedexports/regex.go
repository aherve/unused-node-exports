package unusedexports

const hasExportRegexPattern = `export (?:async )?(?:function|const)`
const exportNameRegexPattern = `export (?:async )?(?:function|const) (\w+)`

const removeExportRegexPattern = `export (.*)`

const importRegexPattern = `import (?:type )?\{\s*([\s\S]*?)\s*\}`
const removeAliasRegexPattern = `(\w+)\s+as\s+(\w+)`
