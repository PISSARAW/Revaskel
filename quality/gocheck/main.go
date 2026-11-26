// Analyse les fonctions Go avec le parseur de la bibliothèque standard.
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

type score struct {
	branches int
	depth int
}

type visitor struct {
	tally *score
	level int
	alternative ast.Node
}

func decision(node ast.Node) bool {
	switch n := node.(type) {
	case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt:
		return true
	case *ast.CaseClause:
		return n.List != nil
	case *ast.CommClause:
		return n.Comm != nil
	case *ast.BinaryExpr:
		return n.Op == token.LAND || n.Op == token.LOR
	}
	return false
}

func block(node ast.Node) bool {
	switch node.(type) {
	case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt, *ast.SwitchStmt, *ast.TypeSwitchStmt, *ast.SelectStmt:
		return true
	}
	return false
}

func (v visitor) Visit(node ast.Node) ast.Visitor {
	if _, nested := node.(*ast.FuncLit); nested {
		return nil
	}
	if decision(node) {
		v.tally.branches++
	}
	if block(node) && node != v.alternative {
		v.level++
	}
	if v.level > v.tally.depth {
		v.tally.depth = v.level
	}
	if branch, ok := node.(*ast.IfStmt); ok {
		v.alternative = branch.Else
	}
	return v
}

func inspect(ft *ast.FuncType, body *ast.BlockStmt, pos token.Position) bool {
	result := score{branches: 1}
	if body != nil {
		ast.Walk(visitor{tally: &result}, body)
	}
	params := ft.Params.NumFields()
	if params > 3 || result.depth > 3 || result.branches > 10 {
		fmt.Printf("%s : paramètres %d/3, profondeur %d/3, complexité %d/10\n", pos, params, result.depth, result.branches)
		return false
	}
	return true
}

func check(path string) bool {
	locations := token.NewFileSet()
	file, err := parser.ParseFile(locations, path, nil, 0)
	if err != nil {
		fmt.Println(err)
		return false
	}
	ok := true
	ast.Inspect(file, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			ok = inspect(x.Type, x.Body, locations.Position(x.Pos())) && ok
		case *ast.FuncLit:
			ok = inspect(x.Type, x.Body, locations.Position(x.Pos())) && ok
		}
		return true
	})
	return ok
}

func main() {
	ok := true
	for _, path := range os.Args[1:] {
		if path != "--" {
			ok = check(path) && ok
		}
	}
	if !ok {
		os.Exit(1)
	}
}
