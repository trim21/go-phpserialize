package encoder

import (
	"sync"
	"unsafe"
)

var ctxPool = sync.Pool{
	New: func() any {
		return &Ctx{
			Buf:         make([]byte, 0, 1024),
			smallBuffer: make([]byte, 0, 20),
			Seen:        make(map[unsafe.Pointer]empty, 100),
		}
	},
}

type empty struct{}

type Ctx struct {
	smallBuffer []byte // a small buffer to encode float and time.Time as string
	Buf         []byte
	Seen        map[unsafe.Pointer]empty
	StackDepth  uint
}

func newCtx() *Ctx {
	ctx := ctxPool.Get().(*Ctx)
	ctx.smallBuffer = ctx.smallBuffer[:0]

	return ctx
}

func freeCtx(ctx *Ctx) {
	ctxPool.Put(ctx)
}
