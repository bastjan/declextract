package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/bastjan/declextract/extract"
)

func main() {
	flag.Parse()

	path := flag.Arg(0)
	if path == "" {
		fmt.Fprintln(os.Stderr, "usage: declextract <path> <declaration>")
		os.Exit(5)
	}
	name := flag.Arg(1)
	if name == "" {
		fmt.Fprintln(os.Stderr, "usage: declextract <path> <declaration>")
		os.Exit(5)
	}

	v, err := extract.ExtractDeclarationFromFile(path, name)
	if err != nil {
		if errors.Is(err, &extract.NotFoundError{}) {
			fmt.Fprintf(os.Stderr, "declaration %q not found in %q\n", name, path)
			os.Exit(3)
		}
		fmt.Fprintln(os.Stderr, "Failed to extract declaration", err)
		os.Exit(1)
	}

	fmt.Print(v)
}
