package migoinfer

// Utility helper functions.

import (
	"go/token"
	"go/types"
	"golang.org/x/tools/go/pointer"
	"golang.org/x/tools/go/ssa"

	"github.com/JorgeGCoelho/gospal/v2/store"
)

func isChan(k store.Key) bool {
	switch t := k.Type().Underlying().(type) {
	case *types.Chan:
		return true
	case *types.Pointer:
		switch t.Elem().Underlying().(type) {
		case *types.Chan:
			return true
		}
	}
	return false
}

// isPtrBasic returns true if k is a pointer to
// primitive type (defined as "go/types".Basic).
func isPtrBasic(k store.Key) bool {
	switch t := k.Type().Underlying().(type) {
	case *types.Pointer:
		switch t.Elem().Underlying().(type) {
		case *types.Basic:
			return true
		}
	}
	return false
}

func isStruct(k store.Key) bool {
	switch t := k.Type().Underlying().(type) {
	case *types.Struct:
		return true
	case *types.Pointer:
		// only pointer to struct is 'Struct'
		switch t.Elem().Underlying().(type) {
		case *types.Struct:
			return true
		}
	}
	return false
}
func followCall(sendPos *map[token.Pos]struct{}, call *ssa.Call, recv *ssa.Value, seen map[*ssa.Value]struct{}, pointerConf *pointer.Config) []token.Pos {
	res := []token.Pos{}
	for i, arg := range call.Call.Args {
		if arg == *recv {
			switch value := call.Call.Value.(type) {
			case *ssa.Function:
				v := ssa.Value(value.Params[i])
				followReferrers(sendPos, &v, &v, seen, pointerConf)
				// ref.Call.Method.
			case *ssa.MakeClosure:
				v := ssa.Value(value.Fn.(*ssa.Function).Params[i])
				followReferrers(sendPos, &v, &v, seen, pointerConf)
			}
		}
	}
	return res
}

func followReferrers(sendPos *map[token.Pos]struct{}, instr *ssa.Value, recv *ssa.Value, seen map[*ssa.Value]struct{}, pointerConf *pointer.Config) {
	if _, present := seen[instr]; present {
		return
	}
	seen[instr] = struct{}{}
	for _, ref := range *(*instr).Referrers() {
		switch ref := ref.(type) {
		case *ssa.Send:
			(*sendPos)[ref.Pos()] = struct{}{}
		case *ssa.Call:
			followCall(sendPos, ref, recv, seen, pointerConf)
		case *ssa.Store:
			followReferrers(sendPos, &ref.Addr, &ref.Addr, seen, pointerConf)
			pointerConf.AddQuery(ref.Addr)
		case *ssa.Alloc:
			followReferrers(sendPos, instr, instr, seen, pointerConf)
		case *ssa.MakeClosure:
			for i, binding := range ref.Bindings {
				if binding == *recv {
					instr := ssa.Value(ref.Fn.(*ssa.Function).FreeVars[i])
					followReferrers(sendPos, &instr, &instr, seen, pointerConf)
				}
			}
		case *ssa.UnOp:
			i := ssa.Value(ref)
			followReferrers(sendPos, &i, instr, seen, pointerConf)

		}
	}
}
