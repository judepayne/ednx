# ednx

A Go library for converting between EDN (Extensible Data Notation) and JSON formats, with an optional CLI tool.

## Library Usage

```bash
go get github.com/judepayne/ednx
```

```go
import "github.com/judepayne/ednx/ednx"

// JSON to EDN
opts := &ednx.EdnConvertOptions{
    KeywordizeKeys: true,
    PrettyPrint:    true,
    WidthLimit:     80,
}
ednData, err := ednx.JsonToEdn(jsonBytes, opts)

// EDN to JSON
jsonOpts := &ednx.JsonConvertOptions{PrettyPrint: true}
jsonData, err := ednx.EdnToJson(ednBytes, jsonOpts)
```

### Options

**EdnConvertOptions:**
- `KeywordizeKeys` - Convert string keys to EDN keywords
- `PrettyPrint` - Format output with proper indentation
- `WidthLimit` - Character width limit for pretty printing (default: 80)

**JsonConvertOptions:**
- `PrettyPrint` - Format JSON with indentation

## CLI Tool

```bash
go install github.com/judepayne/ednx/cmd@latest
```

```bash
# Convert JSON to EDN
ednx -e < input.json

# Convert EDN to JSON  
ednx -j < input.edn

# Pretty print with keywordized keys
ednx -e -p -k < input.json
```

## Features

- Bidirectional EDN ↔ JSON conversion
- Width-aware pretty printing algorithm
- Keywordization of map keys for EDN
- Handles nested data structures
- Clean, library-first API