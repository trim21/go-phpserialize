package encoder

import "strconv"

func appendArrayBegin(ctx *Ctx, fieldNum int64) {
	ctx.b = append(ctx.b, 'a', ':')
	ctx.b = strconv.AppendInt(ctx.b, fieldNum, 10)
	ctx.b = append(ctx.b, ':', '{')
}

func appendString(ctx *Ctx, s string) {
	ctx.b = append(ctx.b, 's', ':')
	ctx.b = strconv.AppendInt(ctx.b, int64(len(s)), 10)
	ctx.b = append(ctx.b, ':')
	ctx.b = strconv.AppendQuote(ctx.b, s)
	ctx.b = append(ctx.b, ';')
}

func appendNil(ctx *Ctx) {
	ctx.b = append(ctx.b, 'N', ';')
}

func appendStringHead(ctx *Ctx, length int64) {
	ctx.b = append(ctx.b, 's', ':')
	ctx.b = strconv.AppendInt(ctx.b, length, 10)
	ctx.b = append(ctx.b, ':')
}
