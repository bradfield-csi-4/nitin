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
	var result int
	var err error

	ast.Inspect(expr, func(n ast.Node) bool {
		result, err = getResult(n)
		if err != nil {
			log.Fatal(err)
		}
		return false
	})
	return result, nil
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

	//result, err := Evaluate(expr)
	//if err != nil {
	//	return
	//}
	//fmt.Println(result)
}

func getResult(n ast.Node) (int, error) {
	// Recursion base case
	lit, ok := n.(*ast.BasicLit)
	if ok {
		fmt.Println("Found literal: ", lit)
		nodeVal, _ := strconv.Atoi(lit.Value)
		return nodeVal, nil
	}

	// Skip parentheses expressions by assigning to child
	parenExpr, isParen := n.(*ast.ParenExpr)
	if isParen {
		n = parenExpr.X
	}

	binExpr, isBinExpr := n.(*ast.BinaryExpr)
	if !isBinExpr {
		return 0, fmt.Errorf("unsupported expression")
	}

	fmt.Println("Found binary expression: ", binExpr)

	x, err := getResult(binExpr.X)
	if err != nil {
		return 0, err
	}
	y, err := getResult(binExpr.Y)
	if err != nil {
		return 0, err
	}

	if binExpr.Op == token.ADD {
		return x + y, nil
	} else if binExpr.Op == token.SUB {
		return x - y, nil
	} else if binExpr.Op == token.MUL {
		return x * y, nil
	} else {
		return 0, fmt.Errorf("unsupported operation")
	}
}
