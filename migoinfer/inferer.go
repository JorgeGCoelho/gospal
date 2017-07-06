package migoinfer

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"github.com/nickng/gospal/callctx"
	"github.com/nickng/gospal/funcs"
	"github.com/nickng/gospal/migoinfer/internal/migoinfer"
	"github.com/nickng/gospal/ssa"
	"github.com/nickng/gospal/store"
	"github.com/nickng/migo"
)

// Inferer is the main MiGo inference entry point.
type Inferer struct {
	Env  migoinfer.Environment // Program environment.
	Info *ssa.Info             // SSA IR.
	MiGo *migo.Program         // MiGo program.

	Raw bool

	outWriter io.Writer // Output stream.
	errWriter io.Writer // Error stream.
}

// New returns a new Inferer, and uses w for logging messages.
func New(info *ssa.Info, w io.Writer) *Inferer {
	inferer := Inferer{
		Env:       migoinfer.NewEnvironment(info),
		Info:      info,
		MiGo:      migo.NewProgram(),
		Raw:       false,
		outWriter: ioutil.Discard,
		errWriter: ioutil.Discard,
	}
	if w != nil {
		inferer.errWriter = w
	}
	return &inferer
}

func (i *Inferer) Analyse() {

	pkg := migoinfer.NewPackage(&i.Env)
	// Package/global variables initialisation.
	for _, p := range i.Info.Prog.AllPackages() {
		pkg.InitGlobals(p)
		pkg.VisitInit(p)
	}
	// Find main packages to start analysis.
	mains, err := ssa.MainPkgs(i.Info.Prog, false)
	if err != nil {
		log.Fatal("Cannot find main package:", err)
	}
	// Call context
	ctx := callctx.Toplevel()
	if l, ok := ctx.(store.Logger); ok {
		l.SetLog(i.errWriter)
	}
	for _, main := range mains {
		if mainFn := main.Func("main"); mainFn != nil {
			mainDef := funcs.MakeCall(funcs.MakeDefinition(mainFn), nil, nil)
			mainFnAnalyser := migoinfer.NewFunction(mainDef, &i.Env, ctx)
			mainFnAnalyser.EnterFunc(mainDef.Function())
		}
	}
	if !i.Raw {
		i.Env.Prog.CleanUp()
	}
	// Print main.main first.
	for _, f := range i.Env.Prog.Funcs {
		if f.SimpleName() == "main.main" {
			fmt.Fprintf(i.outWriter, f.String())
		}
	}
	for _, f := range i.Env.Prog.Funcs {
		if f.SimpleName() != "main.main" {
			fmt.Fprintf(i.outWriter, f.String())
		}
	}
}
