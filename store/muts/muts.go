// Package muts implements store.Value for (rw)mutex locks.
package muts

import (
	"fmt"

	"github.com/jujuyuki/gospal/store"
	"golang.org/x/tools/go/ssa"
)

// Mut is a wrapper for a Mutex lock SSA-value.
type Mut struct {
	ssa.Value

	namespace store.Value
}

// New returns a new mut.
func New(callsite store.Value, val ssa.Value) *Mut {
	return &Mut{
		Value:     val,
		namespace: callsite,
	}
}

func (m *Mut) UniqName() string {
	return fmt.Sprintf("%s.%s_mut.%v", m.namespace.UniqName(), m.Value.Name(), m.Pos())
}
