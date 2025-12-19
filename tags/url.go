package tags

import (
	"fmt"
	"lain/utils/urls"
	"strings"

	"github.com/flosch/pongo2/v6"
)

type urlNode struct {
	routeName string
	params    map[string]pongo2.IEvaluator
}

func url(doc *pongo2.Parser, start *pongo2.Token, arguments *pongo2.Parser) (pongo2.INodeTag, *pongo2.Error) {
	routeNameToken := arguments.MatchType(pongo2.TokenString)
	if routeNameToken == nil {
		return nil, arguments.Error("expected route name string", nil)
	}
	routeName := routeNameToken.Val

	params := make(map[string]pongo2.IEvaluator)

	for arguments.Remaining() > 0 {
		keyToken := arguments.MatchType(pongo2.TokenIdentifier)
		if keyToken == nil {
			return nil, arguments.Error("expected param key identifier", nil)
		}

		if arguments.Match(pongo2.TokenSymbol, "=") == nil {
			return nil, arguments.Error("expected '=' after param key", nil)
		}

		valueExpr, err := arguments.ParseExpression()
		if err != nil {
			return nil, err
		}

		params[keyToken.Val] = valueExpr
	}

	return &urlNode{
		routeName: routeName,
		params:    params,
	}, nil
}

func (n *urlNode) Execute(ctx *pongo2.ExecutionContext, writer pongo2.TemplateWriter) *pongo2.Error {
	path, ok := urls.GetFullPath(n.routeName)
	if !ok {
		return &pongo2.Error{
			Sender:    "tag:url",
			OrigError: fmt.Errorf("route not found: %s", n.routeName),
		}
	}

	for key, expr := range n.params {
		val, err := expr.Evaluate(ctx)
		if err != nil {
			return err
		}
		path = strings.ReplaceAll(path, ":"+key, fmt.Sprintf("%v", val.Interface()))
	}

	_, err := writer.WriteString(path)
	if err != nil {
		return &pongo2.Error{
			Sender:    "tag:url",
			OrigError: err,
		}
	}

	return nil
}
