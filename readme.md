# go-phpserialize

![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/trim21/go-phpserialize?style=flat-square)
[![Go Reference](https://pkg.go.dev/badge/github.com/trim21/go-phpserialize#section-readme.svg)](https://pkg.go.dev/github.com/trim21/go-phpserialize#section-readme)

PHP `serialize()` and `unserialize()` for Go.

Support All go type including `map`, `slice`, `struct`, `array`, and simple type like `int`, `uint` ...etc.

Encoding some type from standard library like `time.Time`, `net.IP` are not supported.
If you have any thought about how to support these types, please create an issue.

## Supported and tested go version

- 1.20
- 1.21
- 1.22
- 1.23rc2

## Install

```console
go get github.com/trim21/go-phpserialize
```

## Usage

See [examples](./example_test.go)

### Marshal

Struct and map will be encoded to php array only.

### Unmarshal

Decoding from php serialized array, class and object are both supported.

go `any` type will be decoded as `map[any]any` or `map[string]any`, based on raw input is `array` or `class`,

keys of `map[any]any` maybe `int64` or `string`.

## Note

go `reflect` package allow you to create dynamic struct with [reflect.StructOf](https://pkg.go.dev/reflect#StructOf),
but please use it with caution.

For performance, this package will try to "compile" input type to a static encoder/decoder
at first time and cache it for future use.

So a dynamic struct may cause memory leak.

## Changelog

### v0.1.0 (not released yet)

Add new `Marshaler` to match `json.Marshaler`.

Go 1.23 has decided to [lock down future uses of `//go:linkname`](https://github.com/golang/go/issues/67401),
So we did a major refactoring in v0.1.0.
For simplicity, support for embed struct has been removed,
if you need this feature, send a Pull Request.

## Security

TL;DR: Don't unmarshal content you can't trust.

Attackers may consume large memory with very few bytes.

php serialized array has a length prefix `a:1:{i:0;s:3:"one";}`, when decoding php serialized array into go `slice` or
go `map`,
`go-phpserialize` may call golang's `make()` to create a map or slice with given length.

So a malicious input like `a:100000000:{}` may become `make([]T, 100000000)` and consume high memory.

If you have to decode some un-trusted bytes, make sure only decode them into fixed-length golang array or struct,
never decode them to `interface`, `slice` or `map`.

## License

Heavily inspired by https://github.com/goccy/go-json

MIT License
