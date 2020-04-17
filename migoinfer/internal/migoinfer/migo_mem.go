package migoinfer

import (
	"github.com/jujuyuki/gospal/store"
	"github.com/jujuyuki/migo"
	"strings"
)

func migoNewMem(mem migo.NamedVar) migo.Statement {
	return &migo.NewMem{Name: mem}
}

func migoRead(v *Instruction, local store.Key) migo.Statement {
	if !strings.Contains(v.Get(local).UniqName(), "mem") { // If not a local memory field.
		return &migo.TauStatement{} //&migo.MemWrite{Name: v.Get(local).UniqName()}
	}
	return &migo.MemRead{Name: local.Name()}
}

func migoWrite(v *Instruction, local store.Key) migo.Statement {
	if !strings.Contains(v.Get(local).UniqName(), "mem") { // If not a local memory field.
		return &migo.TauStatement{} //&migo.MemWrite{Name: v.Get(local).UniqName()}
	}
	return &migo.MemWrite{Name: local.Name()}
}

func migoNewMutex(mu migo.NamedVar) migo.Statement {
	return &migo.NewSyncMutex{Name: mu}
}

func migoLock(mu store.Key) migo.Statement {
	return &migo.SyncMutexLock{Name: mu.Name()}
}

func migoUnlock(mu store.Key) migo.Statement {
	return &migo.SyncMutexUnlock{Name: mu.Name()}
}

func migoNewRWMutex(mu migo.NamedVar) migo.Statement {
	return &migo.NewSyncRWMutex{Name: mu}
}

func migoRLock(mu store.Key) migo.Statement {
	return &migo.SyncRWMutexRLock{Name: mu.Name()}
}

func migoRUnlock(mu store.Key) migo.Statement {
	return &migo.SyncRWMutexRUnlock{Name: mu.Name()}
}
