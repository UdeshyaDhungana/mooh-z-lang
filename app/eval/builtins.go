package eval

import (
	"fmt"

	"github.com/udeshyadhungana/interprerer/app/object"
)

var builtins = map[string]*object.Builtin{
	//common
	"lambai_muji": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Arr))}
			case *object.HashMap:
				return &object.Integer{Value: int64(len(arg.Pairs))}
			default:
				return newError("argument to `lambai_muji` not supported, got %s", args[0].Type())
			}
		},
	},
	// array operations
	"khaad_muji": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d want=2", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJECT {
				return newError("argument to `khaad_muji` not supported, got %s", args[0].Type())
			}

			a := args[0].(*object.Array)
			a.Arr = append(a.Arr, args[1])
			return object.NULL
		},
	},
	"udaa_muji": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 && len(args) != 2 {
				return newError("wrong number of arguments to `udaa_muji`, got=%d want=1 or 2", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJECT {
				return newError("argument to `udaa_muji` not supported, got %s", args[0].Type())
			}

			a := args[0].(*object.Array)
			idx := int64(len(a.Arr) - 1)
			if len(args) == 2 {
				if args[1].Type() != object.INTEGER_OBJ {
					return newError("args[1] of `udaa_muji` expected to be an integer object")
				}
				idx = args[1].(*object.Integer).Value
				if idx >= int64(len(a.Arr)) {
					return newError("cannot `udaa_muji` using index %d, index out of bounds", len(a.Arr))
				}
			}
			popped := a.Arr[idx]
			a.Arr = append(a.Arr[:idx], a.Arr[idx+1:]...)
			return popped
		},
	},
	"bhan_muji": {
		Fn: func(args ...object.Object) object.Object {
			for _, a := range args {
				fmt.Print(a.Inspect())
			}
			fmt.Println()
			return object.NULL
		},
	},
}
