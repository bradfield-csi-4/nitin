package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strconv"
)

// Given an expression containing only int types, evaluate
// the expression and return the result.
func Evaluate(expr ast.Expr) (int, error) {
	// Recursion base case
	lit, ok := expr.(*ast.BasicLit)
	if ok {
		nodeVal, _ := strconv.Atoi(lit.Value)
		return nodeVal, nil
	}

	// Skip parentheses expressions by assigning to child
	parenExpr, isParen := expr.(*ast.ParenExpr)
	if isParen {
		expr = parenExpr.X
	}

	binExpr, isBinExpr := expr.(*ast.BinaryExpr)
	if !isBinExpr {
		return 0, fmt.Errorf("unsupported expression")
	}

	x, err := Evaluate(binExpr.X)
	if err != nil {
		return 0, err
	}
	y, err := Evaluate(binExpr.Y)
	if err != nil {
		return 0, err
	}

	if binExpr.Op == token.ADD {
		return x + y, nil
	} else if binExpr.Op == token.SUB {
		return x - y, nil
	} else if binExpr.Op == token.MUL {
		return x * y, nil
	} else if binExpr.Op == token.QUO {
		return x / y, nil
	} else {
		return 0, fmt.Errorf("unsupported operation: %v", binExpr.Op)
	}
}

func main() {
	expr, err := parser.ParseExpr("2*(2+3)")
	if err != nil {
		log.Fatal(err)
	}
	fset := token.NewFileSet()
	err = ast.Print(fset, expr)
	if err != nil {
		log.Fatal(err)
	}
}
