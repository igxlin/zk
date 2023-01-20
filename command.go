package main

import (
	"os"
	"path/filepath"
)

func Rename(source, target string) {
	source, err := filepath.Abs(source)
	if err != nil {
		panic("Unable to get abspath of " + source)
	}

	target, err = filepath.Abs(target)
	if err != nil {
		panic("Unable to get abspath of " + target)
	}

	if !DocCollection.Contain(source) {
		panic("Database doesn't contain " + source)
	}

	if err := os.Rename(source, target); err != nil {
		panic(err)
	}

	doc := DocCollection[source]
	for backlink := range doc.backlinks {
		DocCollection[backlink].UpdateLinks(source, target)
	}

}
