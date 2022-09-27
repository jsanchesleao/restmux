package parser

import (
	"fmt"
	"restmux/model"

	g "github.com/jsanchesleao/grammatic"
	gm "github.com/jsanchesleao/grammatic/model"
)

const RestmuxGrammar = `
Query := QueryBlock+

QueryBlock := FromClause
              WithClause? as MaybeWithClause

FromClause := FromKW Key
WithClause := WithKW
              WithClauseItem+ as WithClauseItems

WithClauseItem := Key Equals Value

Value := Number | String

WithKW := /with/
FromKW := /from/

Equals := /=/
Number := $NumberFormat
Key := /\w+/
String := $DoubleQuotedStringFormat

Space := $EmptySpaceFormat (ignore)
`

var grammar = g.Compile(RestmuxGrammar)

func ParseString(query string) (model.RestmuxQuery, error) {
	result := model.RestmuxQuery{}

	tree, err := grammar.Parse("Query", query)
	if err != nil {
		return result, err
	}

	fmt.Printf(tree.PrettyPrint())

	return buildQuery(tree.GetNodeWithType("Query")), nil
}

func buildQuery(tree *gm.Node) model.RestmuxQuery {
	blockNodes := tree.GetNodesWithType("QueryBlock")

	blocks := []model.RestmuxQueryBlock{}
	for _, block := range blockNodes {
		blocks = append(blocks, buildQueryBlock(block))
	}

	return model.RestmuxQuery{
		Name:      "",
		Endpoints: []model.RestmuxEndpoint{},
		Blocks:    blocks,
	}
}

func buildQueryBlock(tree *gm.Node) model.RestmuxQueryBlock {
	from := tree.GetNodeWithType("FromClause").
		GetNodeWithType("Key").Token.Value

	with := []model.RestmuxQueryWithBlock{}
	maybeWithNode := tree.GetNodeWithType("MaybeWithClause")
	if maybeWithNode != nil {
		withNodes := maybeWithNode.GetNodeWithType("WithClause").
			GetNodesWithType("WithClauseItems")
		for _, node := range withNodes {
			itemNode := node.GetNodeWithType("WithClauseItem")
			with = append(with, model.RestmuxQueryWithBlock{
				Field: itemNode.GetNodeWithType("Key").Token.Value,
				Value: buildValue(itemNode.GetNodeWithType("Value")),
			})
		}
	}

	return model.RestmuxQueryBlock{
		From: from,
		With: with,
	}
}

func buildValue(tree *gm.Node) string {
	token := tree.GetNodeByIndex(0).Token
	if token == nil {
		return ""
	}

	return token.Value

}
