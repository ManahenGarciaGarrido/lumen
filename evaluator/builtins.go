package evaluator

import (
	"fmt"
	"strings"

	"github.com/ManahenGarciaGarrido/lumen/object"
)

var builtins = map[string]*object.Builtin{
	"print": {
		Fn: func(args ...object.Object) object.Object {
			parts := make([]string, len(args))
			for i, a := range args {
				parts[i] = a.Inspect()
			}
			fmt.Println(strings.Join(parts, " "))
			return NULL
		},
	},
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("len() takes 1 argument, got %d", len(args))
			}
			switch a := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(a.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(a.Elements))}
			}
			return newError("len() not supported for %s", args[0].Type())
		},
	},
	"first": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("first() takes 1 argument")
			}
			arr, ok := args[0].(*object.Array)
			if !ok {
				return newError("first() requires an array, got %s", args[0].Type())
			}
			if len(arr.Elements) == 0 {
				return NULL
			}
			return arr.Elements[0]
		},
	},
	"last": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("last() takes 1 argument")
			}
			arr, ok := args[0].(*object.Array)
			if !ok {
				return newError("last() requires an array, got %s", args[0].Type())
			}
			n := len(arr.Elements)
			if n == 0 {
				return NULL
			}
			return arr.Elements[n-1]
		},
	},
	"rest": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("rest() takes 1 argument")
			}
			arr, ok := args[0].(*object.Array)
			if !ok {
				return newError("rest() requires an array, got %s", args[0].Type())
			}
			if len(arr.Elements) == 0 {
				return NULL
			}
			newElems := make([]object.Object, len(arr.Elements)-1)
			copy(newElems, arr.Elements[1:])
			return &object.Array{Elements: newElems}
		},
	},
	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("push() takes 2 arguments, got %d", len(args))
			}
			arr, ok := args[0].(*object.Array)
			if !ok {
				return newError("push() first argument must be an array, got %s", args[0].Type())
			}
			newElems := make([]object.Object, len(arr.Elements)+1)
			copy(newElems, arr.Elements)
			newElems[len(arr.Elements)] = args[1]
			return &object.Array{Elements: newElems}
		},
	},
	"type": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("type() takes 1 argument, got %d", len(args))
			}
			return &object.String{Value: strings.ToLower(string(args[0].Type()))}
		},
	},
	"str": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("str() takes 1 argument, got %d", len(args))
			}
			return &object.String{Value: args[0].Inspect()}
		},
	},
	"int": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("int() takes 1 argument, got %d", len(args))
			}
			switch a := args[0].(type) {
			case *object.Integer:
				return a
			case *object.String:
				var n int64
				if _, err := fmt.Sscanf(a.Value, "%d", &n); err != nil {
					return newError("int() cannot convert %q to integer", a.Value)
				}
				return &object.Integer{Value: n}
			case *object.Boolean:
				if a.Value {
					return &object.Integer{Value: 1}
				}
				return &object.Integer{Value: 0}
			}
			return newError("int() not supported for %s", args[0].Type())
		},
	},
}
