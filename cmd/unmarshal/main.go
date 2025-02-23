// Copyright (c) VirtualTam
// SPDX-License-Identifier: MIT

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/virtualtam/netscape-go/v2"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("missing input filename")
	}

	filePath := os.Args[1]

	document, err := netscape.UnmarshalFile(filePath)
	if err != nil {
		fmt.Println("failed to unmarshal file:", err)
		os.Exit(1)
	}

	jsonData, err := json.MarshalIndent(document, "", "  ")
	if err != nil {
		fmt.Println("failed to marshal data as JSON:", err)
		os.Exit(1)
	}

	fmt.Println(string(jsonData))
}
