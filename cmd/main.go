package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/judepayne/ednx/ednx"
)

func main() {
	// flags:
	// -e to edn
	// -j to json
	// -p prettify
	// -k keywordize keys. must be used with -e
	// -width widthLimit. must be used with -e and -p

	ednPtr := flag.Bool("e", false, "Convert to Edn")
	jsonPtr := flag.Bool("j", false, "Convert to Json")
	prettyPtr := flag.Bool("p", false, "Prettify Edn/Json")
	keywordPtr := flag.Bool("k", false, "Keywordize Edn keys")
	widthPtr := flag.Int("width", 80, "Character width limit for prettifying Edn")

	flag.Parse()

	var data []byte
	var err error

	if flag.NArg() == 0 {
		// Read from stdin
		data, err = io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Println("Error reading stdin:", err)
			os.Exit(1)
		}
	} else {
		// Read from file
		filename := flag.Arg(0)
		f, err := os.Open(filename)
		if err != nil {
			fmt.Println("Could not open file", filename)
			os.Exit(1)
		}
		defer f.Close()
		data, err = io.ReadAll(f)
		if err != nil {
			fmt.Println("Error reading file:", err)
			os.Exit(1)
		}
	}

	if *ednPtr && *jsonPtr {
		fmt.Println("Error: cannot specify both -e and -j")
		os.Exit(1)
	} else if !*ednPtr && !*jsonPtr {
		fmt.Println("Error: must specify either -e or -j")
		os.Exit(1)
	} else if *jsonPtr && *keywordPtr {
		fmt.Println("Error: The keyword flag can only be used when converting to edn")
		os.Exit(1)
	}

	var ednOpts = ednx.EdnConvertOptions{
		KeywordizeKeys: *keywordPtr,
		PrettyPrint:    *prettyPtr,
		WidthLimit:     *widthPtr,
	}

	var jsonOpts = ednx.JsonConvertOptions{
		PrettyPrint: *prettyPtr,
	}

	if *ednPtr {
		ednData, err := ednx.JsonToEdn(data, &ednOpts)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(string(ednData))
	} else {
		jsonData, err := ednx.EdnToJson(data, &jsonOpts)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(string(jsonData))
	}
}
