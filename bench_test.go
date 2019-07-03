package gel_test

import (
	"testing"

	"github.com/Stromberg/gel"
)

func BenchmarkParse0(b *testing.B) {
	fset := gel.NewFileSet()
	for i := 0; i < b.N; i++ {
		_, _ = gel.ParseString(fset, "", "0")
	}
}

func BenchmarkEval0(b *testing.B) {
	fset := gel.NewFileSet()
	node, _ := gel.ParseString(fset, "", "0")
	for i := 0; i < b.N; i++ {
		scope, _ := gel.NewScope(fset)
		_, _ = scope.Eval(node)
	}
}

func BenchmarkParseFib(b *testing.B) {
	fset := gel.NewFileSet()
	for i := 0; i < b.N; i++ {
		_, _ = gel.ParseString(fset, "", "(func fib (n) (if (== n 0) 0 (if (== n 1) 1 (+ (fib (- n 1)) (fib (- n 2))))))")
	}
}

func BenchmarkEvalFib10(b *testing.B) {
	fset := gel.NewFileSet()
	node, _ := gel.ParseString(fset, "", "(func fib (n) (if (== n 0) 0 (if (== n 1) 1 (+ (fib (- n 1)) (fib (- n 2)))))) (fib 10)")
	for i := 0; i < b.N; i++ {
		scope, _ := gel.NewScope(fset)
		_, _ = scope.Eval(node)
	}
}

func BenchmarkEvalFib10ExistingScope(b *testing.B) {
	fset := gel.NewFileSet()
	node, _ := gel.ParseString(fset, "", "(func fib (n) (if (== n 0) 0 (if (== n 1) 1 (+ (fib (- n 1)) (fib (- n 2)))))) (fib 10)")
	scope, _ := gel.NewScope(fset)
	for i := 0; i < b.N; i++ {
		_, _ = scope.Eval(node)
	}
}

func BenchmarkParseMacroEval(b *testing.B) {
	fset := gel.NewFileSet()
	node, _ := gel.ParseString(fset, "", "(repeatedly 100 (# true))")
	scope, _ := gel.NewScope(fset)
	for i := 0; i < b.N; i++ {
		_, _ = scope.Eval(node)
	}
}

func BenchmarkParseFuncEval(b *testing.B) {
	fset := gel.NewFileSet()
	node, _ := gel.ParseString(fset, "", "(repeatedly 10000 (func [] true))")
	scope, _ := gel.NewScope(fset)
	for i := 0; i < b.N; i++ {
		_, _ = scope.Eval(node)
	}
}

func BenchmarkParseMapMacroEval(b *testing.B) {
	fset := gel.NewFileSet()
	node, _ := gel.ParseString(fset, "", "(map (# (if (and (>= %1 0) (<= %1 100)) %1 3.14)) (range 0 10000 1))")
	scope, _ := gel.NewScope(fset)
	for i := 0; i < b.N; i++ {
		_, _ = scope.Eval(node)
	}
}

func BenchmarkParseMapFuncEval(b *testing.B) {
	fset := gel.NewFileSet()
	node, _ := gel.ParseString(fset, "", "(map (func [x] (if (and (>= x 0) (<= x 100)) x 3.14)) (range 0 10000 1))")
	scope, _ := gel.NewScope(fset)
	for i := 0; i < b.N; i++ {
		_, _ = scope.Eval(node)
	}
}
