package ssa

import (
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

// MainPkgs returns the main packages in the program.
func MainPkgs(prog *ssa.Program) ([]*ssa.Package, error) {
	pkgs := prog.AllPackages()

	mains := ssautil.MainPackages(pkgs)
	if len(mains) == 0 {
		return nil, ErrNoMainPkgs
	}
	return mains, nil
}
