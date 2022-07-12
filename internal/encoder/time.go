package encoder

import (
	"time"
	"unsafe"

	"github.com/trim21/go-phpserialize/internal/runtime"
)

func compileTime(rt *runtime.Type) (encoder, error) {
	if rt != timeType {
		panic("err")
	}

	return encodeTime, nil
}

func encodeTime(ctx *Ctx, b []byte, p uintptr) ([]byte, error) {
	t := **(**time.Time)(unsafe.Pointer(&p))
	buf := t.AppendFormat(ctx.smallBuffer[:0], time.RFC3339Nano)
	b = appendByteString(b, buf)
	ctx.smallBuffer = buf
	return b, nil
}
