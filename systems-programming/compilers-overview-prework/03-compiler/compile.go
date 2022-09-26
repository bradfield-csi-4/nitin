package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"log"
	"strconv"
)

var memMap = map[string]int{"x": 1, "y": 2}

// Given an AST node corresponding to a function (guaranteed to be
// of the form `func f(x, y byte) byte`), compile it into assembly
// code.
//
// Recall from the README that the input parameters `x` and `y` should
// be read from memory addresses `1` and `2`, and the return value
// should be written to memory address `0`.
func compile(node *ast.FuncDecl) (string, error) {
	var asm string

	statements := node.Body.List

	for i := 0; i < len(statements); i++ {
		switch expr := statements[i].(type) {
		case *ast.ReturnStmt:
			evaluateExpr(&expr.Results[0], &asm)
			return asm + "pop 0\nhalt", nil
		case *ast.AssignStmt:
			evaluateExpr(&expr.Rhs[0], &asm)
			identName := expr.Lhs[0].(*ast.Ident).Name
			asm += fmt.Sprintf("pop %v\n", memMap[identName])
		default:
			return "", fmt.Errorf("unsupported statement type: %T", expr)
		}
	}
	return "", fmt.Errorf("probably missing return statement")
}

func evaluateExpr(stmt *ast.Expr, asm *string) {
	ast.Inspect(*stmt, func(n ast.Node) bool {
		err := getAsm(asm, n)
		if err != nil {
			log.Fatal(err)
		}
		return false
	})
}

func getAsm(asm *string, n ast.Node) error {
	var err error

	// Base case
	switch expr := n.(type) {
	case *ast.BasicLit:
		value, _ := strconv.Atoi(expr.Value)
		*asm += fmt.Sprintf("pushi %v\n", value)
		return nil
	case *ast.Ident:
		*asm += fmt.Sprintf("push %v\n", memMap[expr.Name])
		return nil
	}

	n = stepIntoParens(n)

	binExpr, isBinExpr := n.(*ast.BinaryExpr)
	if !isBinExpr {
		return fmt.Errorf("unsupported expression")
	}

	err = getAsm(asm, binExpr.X)
	if err != nil {
		return err
	}
	err = getAsm(asm, binExpr.Y)
	if err != nil {
		return err
	}

	err = appendBinaryOp(asm, binExpr)
	if err != nil {
		return err
	}

	return nil
}

func appendBinaryOp(asm *string, binExpr *ast.BinaryExpr) error {
	switch binExpr.Op {
	case token.ADD:
		*asm += "add\n"
	case token.SUB:
		*asm += "sub\n"
	case token.MUL:
		*asm += "mul\n"
	case token.QUO:
		*asm += "div\n"
	case token.EQL:
		*asm += "eq\n"
	case token.LSS:
		*asm += "lt\n"
	case token.GTR:
		*asm += "gt\n"
	case token.NEQ:
		*asm += "neq\n"
	case token.LEQ:
		*asm += "leq\n"
	case token.GEQ:
		*asm += "geq\n"
	default:
		return fmt.Errorf("unsupported operation")
	}
	return nil
}

// Skip parentheses expressions by elevating child node (i.e. X)
func stepIntoParens(n ast.Node) ast.Node {
	parenExpr, isParen := n.(*ast.ParenExpr)
	if isParen {
		n = parenExpr.X
	}
	return n
}
