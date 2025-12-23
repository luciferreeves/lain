package tags

import (
	"github.com/flosch/pongo2/v6"
)

func init() {
	pongo2.RegisterTag("static", static)
}

func static(doc *pongo2.Parser, start *pongo2.Token, arguments *pongo2.Parser) (pongo2.INodeTag, *pongo2.Error) {
	pathToken := arguments.MatchType(pongo2.TokenString)
	if pathToken == nil {
		return nil, arguments.Error("Expected a string path", nil)
	}

	return &tagStaticNode{
		path: pathToken.Val,
	}, nil
}

type tagStaticNode struct {
	path string
}

func (node *tagStaticNode) Execute(ctx *pongo2.ExecutionContext, writer pongo2.TemplateWriter) *pongo2.Error {
	writer.WriteString("/static/" + node.path)
	return nil
}
