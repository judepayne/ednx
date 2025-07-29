package ednx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"olympos.io/encoding/edn"
)

// removeEdnEscaping removes HTML escaping from EDN/JSON output
func removeEdnEscaping(data []byte) []byte {
	data = bytes.ReplaceAll(data, []byte("\\u003e"), []byte(">"))
	data = bytes.ReplaceAll(data, []byte("\\u003c"), []byte("<"))
	data = bytes.ReplaceAll(data, []byte("\\u0026"), []byte("&"))
	return data
}

// JsonToEdn converts JSON data to EDN format with configurable options.
// If the WidthLimit is omitted, it defaults to 80 characters.
// If opts is nil, default options are used.
func JsonToEdn(data []byte, opts *EdnConvertOptions) ([]byte, error) {
	// Handle nil options by creating default options
	if opts == nil {
		opts = &EdnConvertOptions{}
	}
	
	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, ErrJsonUnmarshal
	}
	ednData := convertToEdnValue(v, opts)
	var widthLimit int
	if opts.WidthLimit == 0 {
		widthLimit = 80
	} else {
		widthLimit = opts.WidthLimit
	}
	if opts.PrettyPrint {
		prettyEdn, err := prettifyEdn(ednData, 0, widthLimit)
		if err != nil {
			return nil, ErrEdnPrettify
		}
		return removeEdnEscaping(prettyEdn), nil

	} else {
		normalEdn, err := edn.Marshal(ednData)
		if err != nil {
			return nil, ErrEdnMarshal
		}
		return removeEdnEscaping(normalEdn), nil
	}
}

// convertToEdnValue recursively converts JSON values to EDN values
func convertToEdnValue(v any, opts *EdnConvertOptions) any {
	// Handle nil options by creating default options
	if opts == nil {
		opts = &EdnConvertOptions{}
	}
	
	switch val := v.(type) {
	case map[string]any:
		ednMap := make(map[any]any)
		for k, v := range val {
			var key any
			if opts.KeywordizeKeys {
				key = edn.Keyword(k)
			} else {
				key = k
			}
			ednMap[key] = convertToEdnValue(v, opts)
		}
		return ednMap
	case []any:
		for i, elem := range val {
			val[i] = convertToEdnValue(elem, opts)
		}
		return val
	default:
		return val
	}
}

// prettifyEdn is a width-aware EDN pretty printer that formats output to fit within character limits.
// Algorithm:
// 1. For primitives: marshal directly
// 2. For collections (arrays/maps):
//   - First attempt single-line formatting
//   - Calculate total length including current indentation
//   - If exceeds width limit, switch to multi-line with proper indentation
//   - Recursively format nested elements with updated indent levels
//
// 3. Uses "compact-first" strategy similar to Clojure's pprint
func prettifyEdn(value any, currentIndent int, widthLimit int) ([]byte, error) {
	switch val := value.(type) {
	case string, int, int64, float64, bool, nil:
		// Primitives - marshal directly
		return edn.Marshal(val)

	case []any:
		// Try single line first
		elements := make([]string, len(val))
		totalLen := 2 // for brackets []

		for i, elem := range val {
			elemBytes, err := prettifyEdn(elem, currentIndent+1, widthLimit)
			if err != nil {
				return nil, err
			}
			elements[i] = string(elemBytes)
			totalLen += len(elements[i])
			if i > 0 {
				totalLen += 1 // space separator
			}
		}

		// Check if single line fits
		if totalLen+currentIndent <= widthLimit {
			singleLine := "[" + strings.Join(elements, " ") + "]"
			return []byte(singleLine), nil
		}

		// Multi-line format
		var buf bytes.Buffer
		buf.WriteString("[")
		for i, elem := range val {
			if i == 0 {
				elemBytes, err := prettifyEdn(elem, currentIndent+1, widthLimit)
				if err != nil {
					return nil, err
				}
				buf.Write(elemBytes)
			} else {
				buf.WriteString("\n")
				buf.WriteString(strings.Repeat(" ", currentIndent+1))
				elemBytes, err := prettifyEdn(elem, currentIndent+1, widthLimit)
				if err != nil {
					return nil, err
				}
				buf.Write(elemBytes)
			}
		}
		buf.WriteString("]")
		return buf.Bytes(), nil

	case map[any]any:
		// Collect all key-value pairs with consistent formatting
		pairs := make([]string, 0, len(val))
		totalLen := 2 // for braces {}

		for k, v := range val {
			keyBytes, err := edn.Marshal(k)
			if err != nil {
				return nil, err
			}
			valueBytes, err := prettifyEdn(v, currentIndent+1, widthLimit)
			if err != nil {
				return nil, err
			}
			pair := string(keyBytes) + " " + string(valueBytes)
			pairs = append(pairs, pair)
			totalLen += len(pair)
			if len(pairs) > 1 {
				totalLen += 2 // ", " separator
			}
		}

		// Check if single line fits
		if totalLen+currentIndent <= widthLimit {
			singleLine := "{" + strings.Join(pairs, ", ") + "}"
			return []byte(singleLine), nil
		}

		// Multi-line format - use same consistent formatting as single-line
		var buf bytes.Buffer
		buf.WriteString("{")
		for i, pair := range pairs {
			if i == 0 {
				buf.WriteString(pair)
			} else {
				buf.WriteString(",\n" + strings.Repeat(" ", currentIndent+1))
				buf.WriteString(pair)
			}
		}
		buf.WriteString("}")
		return buf.Bytes(), nil

	default:
		// Fallback to regular marshal
		return edn.Marshal(value)
	}
}

// EdnToJson converts EDN data to JSON format with configurable options.
// If opts is nil, default options are used.
func EdnToJson(data []byte, opts *JsonConvertOptions) ([]byte, error) {
	// Handle nil options by creating default options
	if opts == nil {
		opts = &JsonConvertOptions{}
	}
	
	var v any
	if err := edn.Unmarshal(data, &v); err != nil {
		return nil, ErrEdnUnmarshal
	}
	jsonData := convertToJsonValue(v)
	if opts.PrettyPrint {
		prettyJson, err := json.MarshalIndent(jsonData, "", " ")
		if err != nil {
			return nil, ErrJsonPrettify
		} else {
			return removeEdnEscaping(prettyJson), err
		}
	} else {
		normalJson, err := json.Marshal(jsonData)
		if err != nil {
			return nil, ErrJsonMarshal
		} else {
			return removeEdnEscaping(normalJson), err
		}
	}
}

// convertToJsonValue recursively converts EDN values to JSON values
func convertToJsonValue(v any) any {
	switch val := v.(type) {
	case map[any]any:
		jsonMap := make(map[string]any)
		for k, v := range val {
			var keyStr string
			if keyword, ok := k.(edn.Keyword); ok {
				keyStr = string(keyword)
			} else if str, ok := k.(string); ok {
				keyStr = str
			} else {
				keyStr = fmt.Sprintf("%v", k) // fallback
			}
			jsonMap[keyStr] = convertToJsonValue(v)
		}
		return jsonMap
	case []any:
		{
			for i, elem := range val {
				val[i] = convertToJsonValue(elem)
			}
		}
		return val
	default:
		return val

	}
}