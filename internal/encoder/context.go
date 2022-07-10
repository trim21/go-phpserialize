package encoder

import (
	"sync"
	"unsafe"
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
	Iter hiter
}

func (iter *mapIter) reset() {
	iter.Iter = hiter{}
}

var mapCtxPool = sync.Pool{
	New: func() any {
		return &mapIter{}
	},
}

func newMapCtx() *mapIter {
	ctx := mapCtxPool.Get().(*mapIter)
	ctx.Iter = hiter{}

	return ctx
}

func freeMapCtx(ctx *mapIter) {
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
