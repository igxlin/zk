package main

import (
	"bytes"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

var md = goldmark.New(
	goldmark.WithExtensions(extension.TaskList),
)

type Doc struct {
	path      string
	content   []byte
	links     Set[string]
	backlinks Set[string]
}

func (doc *Doc) UpdateLinks(source, target string) {
	dirname := filepath.Dir(doc.path)
	relpath, err := filepath.Rel(dirname, target)
	if err != nil {
		return
	}

	for link := range doc.links {
		if filepath.Join(dirname, link) != source {
			continue
		}

		doc.content = bytes.Replace(doc.content, []byte(link), []byte(relpath), -1)
		// TODO: Update doc.links
	}

	doc.Overwrite()
}

func (doc *Doc) Overwrite() {
	os.WriteFile(doc.path, doc.content, 0644)
}

func NewDoc(filename string) *Doc {
	// NOTE: so far, only markdown is supported
	if !strings.HasSuffix(filename, ".md") {
		return nil
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		return nil
	}

	root := md.Parser().Parse(
		text.NewReader(content),
		parser.WithContext(parser.NewContext()),
	)

	links := Set[string]{}
	err = ast.Walk(root, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			switch link := n.(type) {
			case *ast.Link:
				href, err := url.PathUnescape(string(link.Destination))
				if err != nil {
					return ast.WalkStop, err
				}

				links.Insert(href)
			}
		}
		return ast.WalkContinue, nil
	})

	if err != nil {
		return nil
	}

	return &Doc{
		path:      filename,
		content:   content,
		links:     links,
		backlinks: Set[string]{},
	}
}

type Docs Map[string, *Doc]

func (docs Docs) GetOrNew(filename string) *Doc {
	if doc, ok := docs[filename]; ok {
		return doc
	}

	doc := NewDoc(filename)
	docs[filename] = doc
	return doc
}

func NewDocs(dirname string) Docs {
	docs := Docs{}

	dirname, err := filepath.Abs(dirname)
	if err != nil {
		return nil
	}

	if err := filepath.Walk(dirname, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// TODO: Add skip dir config
		if info.IsDir() && (info.Name() == ".git" || info.Name() == "vendor") {
			return filepath.SkipDir
		}

		if doc := NewDoc(path); doc != nil {
			docs[path] = doc
		}

		return nil
	}); err != nil {
		return nil
	}

	for filename, doc := range docs {
		for link := range doc.links {
			abslink := filepath.Join(filepath.Dir(filename), link)
			if _, ok := docs[abslink]; !ok {
				continue
			}

			docs[abslink].backlinks.Insert(filename)
		}
	}

	return docs
}

func (docs Docs) Contain(filename string) bool {
	_, ok := docs[filename]
	return ok
}

var DocCollection Docs
