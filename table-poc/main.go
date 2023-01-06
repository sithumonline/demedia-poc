package main

import (
	"fmt"
	"strings"

	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	_ "github.com/pingcap/tidb/parser/test_driver"
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

func main() {
	pk := "16Uiu2HAmP44YB5WWWdYccDYRzByum6fWDma13csdVUcySzwPMqYx"
	q := fmt.Sprintf("select * from %s_user_items where user_id=1 order by created_at limit 3 offset 10", pk)

	astNode, err := parse(q)
	if err != nil {
		fmt.Printf("parse error: %v\n", err.Error())
		return
	}

	tn := extract(astNode)[0]
	fmt.Printf("astNode tableName = %v\n", tn)

	sp := strings.Split(tn, "_")
	fmt.Printf("is valid table = %t\n", sp[0] == pk)
}
