package main

import "os"

func main() {
	// filename := "/home/gx/notes/wiki/my-book-list.md"
	// getLinksFromFile(filename)
	DocCollection = NewDocs(".")
	if DocCollection == nil {
		panic("Unable to initialize collection")
	}

	if len(os.Args) > 3 && os.Args[1] == "mv" {
		Rename(os.Args[2], os.Args[3])
	}

}
