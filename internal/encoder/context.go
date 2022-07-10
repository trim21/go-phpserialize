package encoder

import (
	"sync"
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/runtime"
)

var ctxPool = sync.Pool{
	New: func() any {
		return &Ctx{
			Buf:         make([]byte, 0, 1024),
			KeepRefs:    make([]unsafe.Pointer, 0, 8),
			floatBuffer: make([]byte, 0, 20),
		}
	},
}

type Ctx struct {
	floatBuffer []byte
	KeepRefs    []unsafe.Pointer
	Buf         []byte
}

func newCtx() *Ctx {
	ctx := ctxPool.Get().(*Ctx)
	ctx.KeepRefs = ctx.KeepRefs[:0]
	ctx.floatBuffer = ctx.floatBuffer[:0]

	return ctx
}

func freeCtx(ctx *Ctx) {
	ctx.KeepRefs = ctx.KeepRefs[:0]

	ctxPool.Put(ctx)
}

// A mapIter is an iterator for ranging over a map.
// See ValueUnsafeAddress.MapRange.
type mapIter struct {
	Iter runtime.HashIter
}

var mapCtxPool = sync.Pool{
	New: func() any {
		return &mapIter{}
	},
}

func newMapCtx() *mapIter {
	ctx := mapCtxPool.Get().(*mapIter)
	return ctx
}

func freeMapCtx(ctx *mapIter) {
	ctx.Iter = runtime.HashIter{}
	mapCtxPool.Put(ctx)
}

type structCtx struct {
	b            []byte
	writtenField int64
}

var structCtxPool = sync.Pool{New: func() any {
	return &structCtx{
		b: make([]byte, 0, 512),
	}
}}

func newStructCtx() *structCtx {
	return structCtxPool.Get().(*structCtx)
}

func freeStructCtx(ctx *structCtx) {
	ctx.b = ctx.b[:]
	ctx.writtenField = 0
	structCtxPool.Put(ctx)
}
