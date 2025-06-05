package eval

import "github.com/udeshyadhungana/interprerer/app/object"

var builtins = map[string]*object.Builtin{
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
}
