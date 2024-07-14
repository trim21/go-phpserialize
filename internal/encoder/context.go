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
			smallBuffer: make([]byte, 0, 20),
		}
	},
}

type Ctx struct {
	smallBuffer []byte // a small buffer to encode float and time.Time as string
	KeepRefs    []unsafe.Pointer
	Buf         []byte
}

func newCtx() *Ctx {
	ctx := ctxPool.Get().(*Ctx)
	ctx.KeepRefs = ctx.KeepRefs[:0]
	ctx.smallBuffer = ctx.smallBuffer[:0]

	return ctx
}

func freeCtx(ctx *Ctx) {
	ctx.KeepRefs = ctx.KeepRefs[:0]

	ctxPool.Put(ctx)
}
