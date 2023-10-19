package jira

import (
	"fmt"
	"os"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

type Renderer struct {
	listItemPrefix string
}

func NewRenderer() renderer.NodeRenderer {
	return &Renderer{}
}

func (r *Renderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	// blocks
	reg.Register(ast.KindDocument, r.renderDocument)
	reg.Register(ast.KindHeading, r.renderHeading)
	reg.Register(ast.KindBlockquote, r.renderBlockquote)
	reg.Register(ast.KindCodeBlock, r.renderCodeBlock)
	reg.Register(ast.KindFencedCodeBlock, r.renderFencedCodeBlock)
	reg.Register(ast.KindHTMLBlock, r.renderDefault)
	reg.Register(ast.KindList, r.renderList)
	reg.Register(ast.KindListItem, r.renderListItem)
	reg.Register(ast.KindParagraph, r.renderParagraph)
	reg.Register(ast.KindTextBlock, r.renderTextBlock)
	reg.Register(ast.KindThematicBreak, r.renderDefault)

	// inlines

	reg.Register(ast.KindAutoLink, r.renderDefault)
	reg.Register(ast.KindCodeSpan, r.renderCodeSpan)
	reg.Register(ast.KindEmphasis, r.renderEmphasis)
	reg.Register(ast.KindImage, r.renderDefault)
	reg.Register(ast.KindLink, r.renderLink)
	reg.Register(ast.KindRawHTML, r.renderDefault)
	reg.Register(ast.KindText, r.renderText)
	reg.Register(ast.KindString, r.renderDefault)
}

func (r *Renderer) renderDefault(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	nodeKind := node.Kind()
	err := fmt.Errorf("no render function for node kind '%s(%d)'", nodeKind.String(), nodeKind)

	return ast.WalkStop, err
}

func (r *Renderer) renderDocument(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (r *Renderer) renderHeading(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Heading)
	if entering {
		_ = w.WriteByte('h')
		_ = w.WriteByte("0123456"[n.Level])
		if n.Attributes() != nil {
			fmt.Fprintf(os.Stderr, "attributes are not yet supported")
		}
		_, _ = w.WriteString(". ")
	} else {
		_ = w.WriteByte('\n')
		if node.NextSibling() != nil {
			_ = w.WriteByte('\n')
		}
	}

	return ast.WalkContinue, nil
}

func (r *Renderer) renderBlockquote(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("{quote}\n")
	} else {
		_, _ = w.WriteString("{quote}\n")
		if node.NextSibling() != nil {
			_ = w.WriteByte('\n')
		}
	}

	return ast.WalkContinue, nil
}

func (r *Renderer) renderCodeBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.CodeBlock)
	attr := n.Attributes()
	if attr != nil {
		fmt.Fprintf(os.Stderr, "attributes=%#v\n", attr)
	}
	if entering {
		_, _ = w.WriteString("{code}\n")
		r.writeLines(w, source, n)
	} else {
		_, _ = w.WriteString("{code}\n")
		if node.NextSibling() != nil {
			_ = w.WriteByte('\n')
		}
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderFencedCodeBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.FencedCodeBlock)
	attr := n.Attributes()
	if attr != nil {
		fmt.Fprintf(os.Stderr, "attributes=%#v\n", attr)
	}
	if entering {
		_, _ = w.WriteString("{code")
		// TODO: probably need to map language to a supported Jira language
		// https://jira.atlassian.com/secure/WikiRendererHelpAction.jspa?section=advanced
		language := n.Language(source)
		if language != nil {
			_ = w.WriteByte(':')
			_, _ = w.Write(language)
		}
		_, _ = w.WriteString("}\n")
		r.writeLines(w, source, n)
	} else {
		_, _ = w.WriteString("{code}\n")
		if node.NextSibling() != nil {
			_ = w.WriteByte('\n')
		}
	}

	return ast.WalkContinue, nil
}

func (r *Renderer) renderList(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.List)
	if !entering {
		if node.NextSibling() != nil && node.FirstChild() != nil {
			_, _ = w.WriteString("\n\n")
		}
		r.listItemPrefix = r.listItemPrefix[:len(r.listItemPrefix)-1]
		return ast.WalkContinue, nil
	}

	if n.IsOrdered() {
		r.listItemPrefix += "#"
	} else {
		r.listItemPrefix += "*"
	}

	return ast.WalkContinue, nil
}

func (r *Renderer) renderListItem(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	// FIXME: list items don't know if they are ordered or unordered
	if entering {
		_, _ = w.WriteString(r.listItemPrefix + " ")
	} else {
		if node.NextSibling() != nil && node.FirstChild() != nil {
			_ = w.WriteByte('\n')
		}
	}

	return ast.WalkContinue, nil
}

func (r *Renderer) renderParagraph(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		if node.NextSibling() != nil {
			_, _ = w.WriteString("\n\n")
		} else {
			_ = w.WriteByte('\n')
		}
	}

	return ast.WalkContinue, nil
}

func (r *Renderer) renderTextBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		if node.NextSibling() != nil && node.FirstChild() != nil {
			_ = w.WriteByte('\n')
		}
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderText(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	n := node.(*ast.Text)
	segment := n.Segment
	_, _ = w.Write(segment.Value(source))
	if n.SoftLineBreak() {
		_ = w.WriteByte(' ')
	}

	return ast.WalkContinue, nil
}

func (r *Renderer) renderCodeSpan(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("{{")
	} else {
		_, _ = w.WriteString("}}")
	}

	return ast.WalkContinue, nil
}

func (r *Renderer) renderEmphasis(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {

	n := node.(*ast.Emphasis)
	tag := byte('_')
	if n.Level == 2 {
		tag = '*'
	}
	_ = w.WriteByte(tag)

	return ast.WalkContinue, nil
}

func (r *Renderer) renderLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {

	n := node.(*ast.Link)
	if entering {
		_ = w.WriteByte('[')
	} else {
		_ = w.WriteByte('|')
		if !html.IsDangerousURL(n.Destination) {
			_, _ = w.Write(util.EscapeHTML(util.URLEscape(n.Destination, true)))
		}
		_ = w.WriteByte(']')
	}

	return ast.WalkContinue, nil
}

func (r *Renderer) writeLines(w util.BufWriter, source []byte, n ast.Node) {
	l := n.Lines().Len()
	for i := 0; i < l; i++ {
		line := n.Lines().At(i)
		_, _ = w.Write(line.Value(source))
	}
}
