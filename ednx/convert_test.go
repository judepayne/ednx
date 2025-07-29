package ednx

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"
)

func readFile(filename string) []byte {
	// Read from file
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("Could not open file", filename)
		os.Exit(1)
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}
	return data
}

func testJsonToEdnRoundTrip(filename string) (bool, error) {
	original := readFile(filename)

	// JSON → EDN → JSON
	edn, err := JsonToEdn(original, &EdnConvertOptions{})
	if err != nil {
		return false, err
	}

	backToJson, err := EdnToJson(edn, &JsonConvertOptions{})
	if err != nil {
		return false, err
	}

	// Compare JSON structures (not strings)
	var orig, result map[string]interface{}
	json.Unmarshal(original, &orig)
	json.Unmarshal(backToJson, &result)

	return reflect.DeepEqual(orig, result), nil
}

func TestJsonToEdnRoundTrip(t *testing.T) {
	successBoard, err := testJsonToEdnRoundTrip("../examples/board.json")
	if err != nil {
		t.Fatal(err)
	}
	if !successBoard {
		t.Errorf("Round-trip failed for board.json")
	}

	successContainers, err := testJsonToEdnRoundTrip("../examples/containers.json")
	if err != nil {
		t.Fatal(err)
	}
	if !successContainers {
		t.Errorf("Round-trip failed for containers.json")
	}
}

func TestNilOptionsHandling(t *testing.T) {
	// Test JsonToEdn with nil options
	jsonData := []byte(`{"key": "value", "number": 42}`)
	
	ednResult, err := JsonToEdn(jsonData, nil)
	if err != nil {
		t.Fatalf("JsonToEdn with nil options failed: %v", err)
	}
	if len(ednResult) == 0 {
		t.Error("JsonToEdn with nil options returned empty result")
	}
	
	// Test EdnToJson with nil options
	jsonResult, err := EdnToJson(ednResult, nil)
	if err != nil {
		t.Fatalf("EdnToJson with nil options failed: %v", err)
	}
	if len(jsonResult) == 0 {
		t.Error("EdnToJson with nil options returned empty result")
	}
	
	// Verify round-trip works with nil options
	var original, result map[string]interface{}
	json.Unmarshal(jsonData, &original)
	json.Unmarshal(jsonResult, &result)
	
	if !reflect.DeepEqual(original, result) {
		t.Error("Round-trip with nil options failed to preserve data")
	}
}
