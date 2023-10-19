package main

import (
	"fmt"
	"os"

	"github.com/bradfordboyle/md2jira/jira"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

func DefaultRenderer() renderer.Renderer {
	return renderer.NewRenderer(renderer.WithNodeRenderers(util.Prioritized(jira.NewRenderer(), 1000)))
}

func main() {
	if len(os.Args) < 2 {
		fail("missing required argument")
	}
	md := goldmark.New(goldmark.WithRenderer(DefaultRenderer()))
	document, err := os.ReadFile(os.Args[1])
	failIf(err, "cannot read file %s", os.Args[1])
	err = md.Convert(document, os.Stdout)
	failIf(err, "cannot convert document")
}

func fail(format string, a ...any) {
	fmt.Fprintf(os.Stderr, "ERROR - "+format+"\n", a...)
	os.Exit(1)
}

func failIf(err error, format string, a ...any) {
	if err != nil {
		msg := fmt.Sprintf(format, a...)
		fmt.Fprintf(os.Stderr, "ERROR - %s : %s\n", msg, err.Error())
		os.Exit(1)
	}
}
