# unused-node-exports

A fast CLI tool to find unused exports in Node.js and TypeScript projects. It scans your git repository for exported functions, that are not imported anywhere else in your codebase.

## Installation

### Using go:

```
go install github.com/aherve/unused-node-exports/v2@latest
```

### Downloading the binary:

Precompiled binaries are available on the [releases page](https://github.com/aherve/unused-node-exports/releases)

## Usage

Scan your current git repository for unused exports:

```
unused-node-exports scan
```

Options:

- `--help, -h`: Show help message.
- `--path, -p`: Path to the git repository to scan (default: current directory).
- `--file-extensions, -e`: Comma-separated list of file extensions to scan (default: .ts, .tsx, .js, .jsx, .mjs, .cjs, .mts, .cts).
- `--output, -o`: Output results to a CSV file.

Example:

```
unused-node-exports scan -p ./my-project -e .ts,.tsx -o unused.csv
```

## Performance and Limitations

This tool prioritizes speed over perfection. It uses regular expressions to find exports and imports in all git tracked files, but does not use a full parser or AST analysis. As a result, it may not catch all edge cases, such as dynamic imports or exports, or certain complex import/export patterns.
