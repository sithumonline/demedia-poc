package internal

import (
	"fmt"
	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	_ "github.com/pingcap/tidb/parser/test_driver"
	"github.com/sithumonline/demedia-poc/core/config"
	"strings"
)

type colX struct {
	colNames []string
}

func (v *colX) Enter(in ast.Node) (ast.Node, bool) {
	if name, ok := in.(*ast.TableName); ok {
		v.colNames = append(v.colNames, name.Name.O)
	}
	return in, false
}

func (v *colX) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

func extract(rootNode *ast.StmtNode) []string {
	v := &colX{}
	(*rootNode).Accept(v)
	return v.colNames
}

func parse(sql string) (*ast.StmtNode, error) {
	p := parser.New()

	stmtNodes, _, err := p.Parse(sql, "", "")
	if err != nil {
		return nil, err
	}

	return &stmtNodes[0], nil
}

func CheckQuery(query string) (bool, error) {
	q := strings.ReplaceAll(query, "\"", "")
	astNode, err := parse(q)
	if err != nil {
		return false, fmt.Errorf("sql parse error: %w", err)
	}

	tn := extract(astNode)[0]
	sp := strings.Split(tn, "_")
	return sp[0] == config.HubHostId, nil
}
