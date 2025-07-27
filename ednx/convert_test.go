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
