package gel_test

import (
	"fmt"
	"testing"

	"github.com/Stromberg/gel"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

var _ = Suite(S{})

type S struct{}

func (S) TestEval(c *C) {
	for _, test := range evalList {
		fset := gel.NewFileSet()
		node, err := gel.ParseString(fset, "", test.code)
		c.Assert(err, IsNil)
		scope := gel.NewScope(fset)
		scope.Create("sprintf", sprintfFn)
		scope.Create("list", listFn)
		scope.Create("append", appendFn)
		value, err := scope.Eval(node)
		if e, ok := test.value.(error); ok {
			c.Assert(err, ErrorMatches, e.Error(), Commentf("Code: %s", test.code))
			c.Assert(value, IsNil)
		} else {
			tvalue := test.value
			if i, ok := tvalue.(int); ok {
				tvalue = int64(i)
			}
			c.Assert(err, IsNil, Commentf("Code: %s", test.code))
			c.Assert(value, DeepEquals, tvalue, Commentf("Code: %s", test.code))
		}
	}
}

func sprintfFn(args ...interface{}) (interface{}, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("sprintf takes at least one format argument")
	}
	format, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("sprintf takes format string as first argument")
	}
	return fmt.Sprintf(format, args[1:]...), nil
}

func listFn(args ...interface{}) (interface{}, error) {
	return args, nil
}

func appendFn(args ...interface{}) (interface{}, error) {
	list, ok := args[0].([]interface{})
	if !ok {
		return nil, fmt.Errorf("append takes list as first argument")
	}
	return append(list, args[1:]...), nil
}

func errorf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

var evalList = []struct {
	code  string
	value interface{}
}{
	// some basics
	{
		`1`,
		1,
	}, {
		`1.0`,
		1.0,
	}, {
		`0x10`,
		16,
	}, {
		`010`,
		8,
	}, {
		`"foo\"bar"`,
		`foo"bar`,
	}, {
		`foo`,
		errorf("twik source:1:1: undefined symbol: foo"),
	}, {
		`(1)`,
		errorf(`twik source:1:2: cannot use 1 as a function`),
	}, {
		`true`,
		true,
	}, {
		`false`,
		false,
	}, {
		`nil`,
		nil,
	}, {
		`1 2 3`,
		3,
	},

	// error
	{
		"(\nerror \"error message\")",
		errorf("twik source:2:1: error message"),
	}, {
		`(error)`,
		errorf("twik source:1:2: error function takes a single string argument"),
	}, {
		`(error 1)`,
		errorf("twik source:1:2: error function takes a single string argument"),
	}, {
		`(error "foo" 2)`,
		errorf("twik source:1:2: error function takes a single string argument"),
	},

	// vec
	{
		`(vec)`,
		[]float64{},
	},
	{
		`(vec 12.0)`,
		[]float64{12.0},
	},
	{
		`(vec 12)`,
		[]float64{12.0},
	},
	{
		`(vec 12 3.14)`,
		[]float64{12.0, 3.14},
	},
	{
		`(vec (list 12 3.14))`,
		[]float64{12.0, 3.14},
	},

	// vec?
	{
		`(vec? (vec))`,
		true,
	},
	{
		`(vec? (list))`,
		false,
	},
	{
		`(vec? 1)`,
		false,
	},
	{
		`(vec? (vec) (vec))`,
		errorf(`twik source:1:2: Wrong number of parameters`),
	},

	// list
	{
		`(list)`,
		[]interface{}{},
	},
	{
		`(list 12.0)`,
		[]interface{}{12.0},
	},
	{
		`(list 12)`,
		[]interface{}{int64(12)},
	},
	{
		`(list 12 3.14)`,
		[]interface{}{int64(12), 3.14},
	},
	{
		`(list 12.0 "do")`,
		[]interface{}{12.0, "do"},
	},
	{
		`(list (vec 12 3.14))`,
		[]interface{}{12.0, 3.14},
	},

	// list?
	{
		`(list? (list))`,
		true,
	},
	{
		`(list? (vec))`,
		false,
	},
	{
		`(list? 1)`,
		false,
	},
	{
		`(list? (list) (list))`,
		errorf(`twik source:1:2: Wrong number of parameters`),
	},

	// dict
	{
		`(dict)`,
		map[interface{}]interface{}{},
	},
	{
		`(dict "d" 12.0)`,
		map[interface{}]interface{}{"d": 12.0},
	},
	{
		`(dict "d")`,
		errorf(`twik source:1:2: Wrong number of parameters`),
	},

	// dict?
	{
		`(dict? (dict))`,
		true,
	},
	{
		`(dict? (vec))`,
		false,
	},
	{
		`(dict? 1)`,
		false,
	},
	{
		`(dict? (dict) (dict))`,
		errorf(`twik source:1:2: Wrong number of parameters`),
	},

	// dict-keys
	{
		`(dict-keys (dict))`,
		[]interface{}{},
	},
	{
		`(dict-keys (dict "d" 12.0))`,
		[]interface{}{"d"},
	},
	{
		`(dict-keys)`,
		errorf(`twik source:1:2: Wrong number of parameters`),
	},
	{
		`(dict-keys "d")`,
		errorf(`twik source:1:2: dict-keys expects a dictionary`),
	},

	// get
	{
		`(get (dict))`,
		errorf(`twik source:1:2: Wrong number of parameters`),
	},
	{
		`(get (vec))`,
		errorf(`twik source:1:2: Wrong number of parameters`),
	},
	{
		`(get (list))`,
		errorf(`twik source:1:2: Wrong number of parameters`),
	},
	{
		`(get (dict "d" 12.0) "d")`,
		12.0,
	},
	{
		`(get (vec 12.0) 0)`,
		12.0,
	},
	{
		`(get (list "d") 0)`,
		"d",
	},
	{
		`(get (dict "d" 12.0) "a")`,
		errorf(`twik source:1:2: Key not found`),
	},
	{
		`(get (vec 12.0) 1)`,
		errorf(`twik source:1:2: Key not found`),
	},
	{
		`(get (list "d") 1)`,
		errorf(`twik source:1:2: Key not found`),
	},

	// contains?
	{
		`(contains? (dict))`,
		errorf(`twik source:1:2: Wrong number of parameters`),
	},
	{
		`(contains? (vec))`,
		errorf(`twik source:1:2: Wrong number of parameters`),
	},
	{
		`(contains? (list))`,
		errorf(`twik source:1:2: Wrong number of parameters`),
	},
	{
		`(contains? (dict "d" 12.0) "d")`,
		true,
	},
	{
		`(contains? (vec 12.0) 0)`,
		true,
	},
	{
		`(contains? (list "d") 0)`,
		true,
	},
	{
		`(contains? (dict "d" 12.0) "a")`,
		false,
	},
	{
		`(contains? (vec 12.0) 1)`,
		false,
	},
	{
		`(contains? (list "d") 1)`,
		false,
	},

	// update
	{
		`(update (dict))`,
		errorf(`twik source:1:2: Wrong number of parameters`),
	},
	{
		`(update (vec))`,
		errorf(`twik source:1:2: Wrong number of parameters`),
	},
	{
		`(update (list))`,
		errorf(`twik source:1:2: Wrong number of parameters`),
	},
	{
		`(update (dict "d" 12.0) "a" 13)`,
		map[interface{}]interface{}{"a": int64(13), "d": 12.0},
	},
	{
		`(update (vec 12.0) 0 13)`,
		[]float64{13.0},
	},
	{
		`(update (list "d") 0 45)`,
		[]interface{}{int64(45)},
	},
	{
		`(update (vec 12.0) 1 34)`,
		errorf(`twik source:1:2: Out of range`),
	},
	{
		`(update (list "d") 1 "a")`,
		errorf(`twik source:1:2: Out of range`),
	},

	// append
	{
		`(append (dict))`,
		errorf(`twik source:1:2: Wrong number of parameters`),
	},
	{
		`(append (vec))`,
		errorf(`twik source:1:2: Wrong number of parameters`),
	},
	{
		`(append (list))`,
		errorf(`twik source:1:2: Wrong number of parameters`),
	},
	{
		`(append (vec 12.0) 13)`,
		[]float64{12.0, 13.0},
	},
	{
		`(append (list "d") 0 45)`,
		[]interface{}{"d", int64(0), int64(45)},
	},
	{
		`(append (vec 12.0) 1 34)`,
		[]float64{12.0, 1.0, 34},
	},
	{
		`(append (list "d") 1 "a")`,
		[]interface{}{"d", int64(1), "a"},
	},

	// range
	{
		`(range)`,
		errorf(`twik source:1:2: Wrong number of parameters`),
	},
	{
		`(range 1)`,
		errorf(`twik source:1:2: Wrong number of parameters`),
	},
	{
		`(range 1 2)`,
		errorf(`twik source:1:2: Wrong number of parameters`),
	},
	{
		`(range 1 2 1)`,
		[]interface{}{int64(1)},
	},
	{
		`(range 1 6 2)`,
		[]interface{}{int64(1), int64(3), int64(5)},
	},
	{
		`(range 1.0 6.0 2.5)`,
		[]interface{}{1.0, 3.5},
	},

	// vec-range
	{
		`(vec-range)`,
		errorf(`twik source:1:2: Wrong number of parameters`),
	},
	{
		`(vec-range 1)`,
		errorf(`twik source:1:2: Wrong number of parameters`),
	},
	{
		`(vec-range 1 2)`,
		errorf(`twik source:1:2: Wrong number of parameters`),
	},
	{
		`(vec-range 1 2 1)`,
		[]float64{1},
	},
	{
		`(vec-range 1 6 2)`,
		[]float64{1, 3, 5},
	},
	{
		`(vec-range 1.0 6.0 2.5)`,
		[]float64{1.0, 3.5},
	},

	// repeat
	{
		`(repeat)`,
		errorf(`twik source:1:2: Wrong number of parameters`),
	},
	{
		`(repeat 1)`,
		errorf(`twik source:1:2: Wrong number of parameters`),
	},
	{
		`(repeat 1 2)`,
		[]interface{}{int64(2)},
	},
	{
		`(repeat 3 6)`,
		[]interface{}{int64(6), int64(6), int64(6)},
	},
	{
		`(repeat 2 6.0)`,
		[]interface{}{6.0, 6.0},
	},

	// vec-repeat
	{
		`(vec-repeat)`,
		errorf(`twik source:1:2: Wrong number of parameters`),
	},
	{
		`(vec-repeat 1)`,
		errorf(`twik source:1:2: Wrong number of parameters`),
	},
	{
		`(vec-repeat 1 2)`,
		[]float64{2},
	},
	{
		`(vec-repeat 3 6)`,
		[]float64{6, 6, 6},
	},
	{
		`(vec-repeat 2 6.0)`,
		[]float64{6.0, 6.0},
	},

	// len
	{
		`(len (dict))`,
		int64(0),
	},
	{
		`(len (vec))`,
		int64(0),
	},
	{
		`(len (list))`,
		int64(0),
	},
	{
		`(len (dict "d" 12.0))`,
		int64(1),
	},
	{
		`(len (vec 12.0))`,
		int64(1),
	},
	{
		`(len (list "d"))`,
		int64(1),
	},
	{
		`(len 1)`,
		errorf(`twik source:1:2: Error in parameter type`),
	},
	{
		`(len)`,
		errorf(`twik source:1:2: Wrong number of parameters`),
	},

	// +
	{
		`(+)`,
		0,
	}, {
		`(+ 1)`,
		1,
	}, {
		`(+ 1 2)`,
		3,
	}, {
		`(+ 1 (+ 2 3))`,
		6,
	}, {
		`(+ "123")`,
		errorf(`twik source:1:2: Error in parameter type`),
	}, {
		`(+ 1.5)`,
		1.5,
	}, {
		`(+ 1.5 1.5)`,
		3.0,
	}, {
		`(+ 1.5 1)`,
		2.5,
	}, {
		`(+ 1 1.5)`,
		2.5,
	}, {
		`(+ 3 (vec 2))`,
		[]float64{5},
	}, {
		`(+ 3 (vec 2 3 4))`,
		[]float64{5, 6, 7},
	},
	{
		`(+ (vec 2 4.5 6) (vec 2 3 4))`,
		[]float64{4, 7.5, 10},
	},
	{
		`(+ (vec 2 4 6) (vec 1.5 3.14))`,
		errorf("twik source:1:2: Vectors of different length"),
	},
	{
		`(+ (vec 2 4) (vec 1.5 3.14 6))`,
		errorf("twik source:1:2: Vectors of different length"),
	},

	// -
	{
		`(-)`,
		errorf(`twik source:1:2: function "-" takes one or more arguments`),
	}, {
		`(- 1)`,
		-1,
	}, {
		`(- 10 1)`,
		9,
	}, {
		`(- 10 1 2)`,
		7,
	}, {
		`(- 10 (- 2 1))`,
		9,
	}, {
		`(- "123")`,
		errorf(`twik source:1:2: Error in parameter type`),
	}, {
		`(- 1.5)`,
		-1.5,
	}, {
		`(- 2.0 1.5)`,
		0.5,
	}, {
		`(- 1.5 1)`,
		0.5,
	}, {
		`(- 1 1.5)`,
		-0.5,
	}, {
		`(- (vec 2 4))`,
		[]float64{-2, -4},
	}, {
		`(- 3 (vec 2))`,
		[]float64{1},
	}, {
		`(- 3 (vec 2 3 4))`,
		[]float64{1, 0, -1},
	},
	{
		`(- (vec 2 4.5 6) (vec 2 3 4))`,
		[]float64{0, 1.5, 2},
	},
	{
		`(- (vec 2 4 6) (vec 1.5 3.14))`,
		errorf("twik source:1:2: Vectors of different length"),
	},
	{
		`(- (vec 2 4) (vec 1.5 3.14 6))`,
		errorf("twik source:1:2: Vectors of different length"),
	},

	// *
	{
		`(*)`,
		1,
	}, {
		`(* 1)`,
		1,
	}, {
		`(* 2 3 4)`,
		24,
	}, {
		`(* 2 (* 3 4))`,
		24,
	}, {
		`(* "123")`,
		errorf(`twik source:1:2: Error in parameter type`),
	}, {
		`(* 1.5)`,
		1.5,
	}, {
		`(* 2.0 1.5)`,
		3.0,
	}, {
		`(* 1.5 1)`,
		1.5,
	}, {
		`(* 1 1.5)`,
		1.5,
	}, {
		`(* 1 (vec 1.5))`,
		[]float64{1.5},
	}, {
		`(* 2 (vec 1.5 3.14 5.6))`,
		[]float64{3.0, 6.28, 11.2},
	},
	{
		`(* (vec 2 4 6) (vec 1.5 3.14 5))`,
		[]float64{3.0, 12.56, 30},
	},
	{
		`(* (vec 2 4 6) (vec 1.5 3.14))`,
		errorf("twik source:1:2: Vectors of different length"),
	},
	{
		`(* (vec 2 4) (vec 1.5 3.14 6))`,
		errorf("twik source:1:2: Vectors of different length"),
	},

	// /
	{
		`(/)`,
		errorf(`twik source:1:2: function "/" takes two or more arguments`),
	}, {
		`(/ 1)`,
		errorf(`twik source:1:2: function "/" takes two or more arguments`),
	}, {
		`(/ 10 2)`,
		5,
	}, {
		`(/ 30 3 2)`,
		5,
	}, {
		`(/ 30 (/ 10 2))`,
		6,
	}, {
		`(/ 10 "123")`,
		errorf(`twik source:1:2: Error in parameter type`),
	}, {
		`(/ 10.0 2.0)`,
		5.0,
	}, {
		`(/ 10.0 2)`,
		5.0,
	}, {
		`(/ 10 2.0)`,
		5.0,
	}, {
		`(/ 3 (vec 2))`,
		[]float64{1.5},
	}, {
		`(/ 3 (vec 2 3 4))`,
		[]float64{1.5, 1.0, 0.75},
	},
	{
		`(/ (vec 2 4.5 6) (vec 2 3 4))`,
		[]float64{1.0, 1.5, 1.5},
	},
	{
		`(/ (vec 2 4 6) (vec 1.5 3.14))`,
		errorf("twik source:1:2: Vectors of different length"),
	},
	{
		`(/ (vec 2 4) (vec 1.5 3.14 6))`,
		errorf("twik source:1:2: Vectors of different length"),
	},

	// ==
	{
		`(== "a" "a")`,
		true,
	}, {
		`(== "a" "b")`,
		false,
	}, {
		`(== 42 42)`,
		true,
	}, {
		`(== 42 43)`,
		false,
	}, {
		`(== 42 "a")`,
		false,
	}, {
		`(== 42 42.0)`,
		false,
	}, {
		`(== 1 2 3)`,
		errorf("twik source:1:2: == takes two values"),
	}, {
		`(==)`,
		errorf("twik source:1:2: == takes two values"),
	},

	// !=
	{
		`(!= "a" "a")`,
		false,
	}, {
		`(!= "a" "b")`,
		true,
	}, {
		`(!= 42 42)`,
		false,
	}, {
		`(!= 42 43)`,
		true,
	}, {
		`(!= 42 "a")`,
		true,
	}, {
		`(!= 42 42.0)`,
		true,
	}, {
		`(!= 1 2 3)`,
		errorf("twik source:1:2: != takes two values"),
	}, {
		`(!=)`,
		errorf("twik source:1:2: != takes two values"),
	},

	// <
	{
		`(< 1 2)`,
		true,
	},
	{
		`(< 1 1)`,
		false,
	},
	{
		`(< 1 1.0)`,
		false,
	},
	{
		`(< 1.0 1.0)`,
		false,
	},
	{
		`(< 1.0 1)`,
		false,
	},

	// >

	{
		`(> 1 2)`,
		false,
	},
	{
		`(> 1 1)`,
		false,
	},
	{
		`(> 1 1.0)`,
		false,
	},
	{
		`(> 1.0 1.0)`,
		false,
	},
	{
		`(> 1.0 1)`,
		false,
	},
	{
		`(> 2.0 1)`,
		true,
	},

	// <=

	{
		`(<= 1 2)`,
		true,
	},
	{
		`(<= 2 1)`,
		false,
	},
	{
		`(<= 1 1.0)`,
		true,
	},
	{
		`(<= 1.0 1.0)`,
		true,
	},
	{
		`(<= 1.0 1)`,
		true,
	},

	// >=

	{
		`(>= 1 2)`,
		false,
	},
	{
		`(>= 2 1)`,
		true,
	},
	{
		`(>= 1 1.0)`,
		true,
	},
	{
		`(>= 1.0 1.0)`,
		true,
	},
	{
		`(>= 1.0 1)`,
		true,
	},

	// or
	{
		`(or)`,
		false,
	}, {
		`(or false 1 2 (error "must not get here"))`,
		1,
	}, {
		`(or (error "boom") 1 2 3)`,
		errorf("twik source:1:6: boom"),
	},

	// and
	{
		`(and)`,
		true,
	}, {
		`(and 1 2 3)`,
		3,
	}, {
		`(and false (error "must not get here"))`,
		false,
	}, {
		`(and (error "boom") true)`,
		errorf("twik source:1:7: boom"),
	},

	// var
	{
		`(var x (+ 1 2)) x`,
		3,
	}, {
		`(var x) x`,
		nil,
	}, {
		`(var x 1 2)`,
		errorf("twik source:1:2: var takes one or two arguments"),
	}, {
		`(var)`,
		errorf("twik source:1:2: var takes one or two arguments"),
	}, {
		"(var x)\n(var x)",
		errorf("twik source:2:2: symbol already defined in current scope: x"),
	},

	// set
	{
		`(var x) (set x 2) (+ x 3)`,
		5,
	}, {
		`(set x 1)`,
		errorf("twik source:1:2: cannot set undefined symbol: x"),
	}, {
		`(var x) (set x 1 2)`,
		errorf(`twik source:1:10: function "set" takes two arguments`),
	}, {
		`(var x) (set x)`,
		errorf(`twik source:1:10: function "set" takes two arguments`),
	}, {
		`(var x) (set)`,
		errorf(`twik source:1:10: function "set" takes two arguments`),
	},

	// do
	{
		`(do)`,
		nil,
	}, {
		`(do 1 2 3)`,
		3,
	}, {
		`(var x 1) (do (set x 2) x)`,
		2,
	}, {
		`(var x 1) (do (set x 2)) x`,
		2,
	}, {
		`(var x 1) (do (var x) (set x 2) x)`,
		2,
	}, {
		`(var x 1) (do (var x) (set x 2)) x`,
		1,
	},

	// func
	{
		`((func (a b) (+ a b)) 1 2)`,
		3,
	}, {
		`(var add (do (var x 0) (func (n) (set x (+ x n)) x))) (add 1) (add 2)`,
		3,
	}, {
		`(func add (a b) (+ a b)) (add 1 2)`,
		3,
	}, {
		`(func)`,
		errorf("twik source:1:2: func takes three or more arguments"),
	}, {
		`(func x)`,
		errorf("twik source:1:2: func takes three or more arguments"),
	}, {
		`(func 1 2)`,
		errorf("twik source:1:2: func takes a list of parameters"),
	}, {
		`(func f 2)`,
		errorf("twik source:1:2: func takes a list of parameters"),
	}, {
		`(func f (a)) (f 1 2)`,
		errorf(`twik source:1:2: func takes a body sequence`),
	}, {
		"(var f (func (a) 1))\n(f 1 2)",
		errorf(`twik source:2:2: anonymous function takes one argument`),
	}, {
		"(func f () 1)\n(f 1)",
		errorf(`twik source:2:2: function "f" takes no arguments`),
	}, {
		"(func f (a) 1)\n(f 1 2)",
		errorf(`twik source:2:2: function "f" takes one argument`),
	}, {
		"(func f (a b) 1)\n(f 1)",
		errorf(`twik source:2:2: function "f" takes 2 arguments`),
	},

	// if
	{
		`(if true 1)`,
		1,
	}, {
		`(if 0 1)`,
		1,
	}, {
		`(if false 1)`,
		false,
	}, {
		`(if false 1 2)`,
		2,
	}, {
		`(if)`,
		errorf(`twik source:1:2: function "if" takes two or three arguments`),
	}, {
		`(if 1)`,
		errorf(`twik source:1:2: function "if" takes two or three arguments`),
	},

	// cond
	{
		`(cond true 1)`,
		1,
	}, {
		`(cond 0 1)`,
		1,
	}, {
		`(cond false 1)`,
		false,
	}, {
		`(cond false 1 2)`,
		2,
	}, {
		`(cond true 1 2)`,
		1,
	}, {
		`(cond false 1 true 2)`,
		2,
	}, {
		`(cond false 1 false 2 3)`,
		3,
	}, {
		`(cond)`,
		errorf(`twik source:1:2: function "cond" takes two or more arguments`),
	}, {
		`(cond 1)`,
		errorf(`twik source:1:2: function "cond" takes two or more arguments`),
	},

	// for
	{
		`(for 1 2 3)`,
		errorf("twik source:1:2: for takes four or more arguments"),
	}, {
		`(for (error "init") (error "test") (error "step") (error "code"))`,
		errorf("twik source:1:7: init"),
	}, {
		`(for () (error "test") (error "step") (error "code"))`,
		errorf("twik source:1:10: test"),
	}, {
		`(for () () (error "step") (error "code"))`,
		errorf("twik source:1:28: code"),
	}, {
		`(for () () (error "step") ())`,
		errorf("twik source:1:13: step"),
	}, {
		`(for (var i 0) false () ()) i`,
		errorf("twik source:1:29: undefined symbol: i"),
	}, {
		`(var x 0) (for (var i 0) (!= i 4) (set i (+ i 1)) (set x (+ x i)) (* 2 x))`,
		12,
	},

	// identity
	{
		`(identity 1)`,
		1,
	},

	// calling of custom functions
	{
		`(sprintf "Value: %.02f" 1.0)`,
		"Value: 1.00",
	},
}
