# go-phpserialize

![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/trim21/go-phpserialize?style=flat-square)
[![Go Reference](https://pkg.go.dev/badge/github.com/trim21/go-phpserialize#section-readme.svg)](https://pkg.go.dev/github.com/trim21/go-phpserialize#section-readme)

PHP `serialize()` and `unserialize()` for Go.

Support All go type including `map`, `slice`, `struct`, `array`, and simple type like `int`, `uint` ...etc.

Encoding some type from standard library like `time.Time`, `net.IP` are not supported.
If you have any thought about how to support these types, please create an issue.

## supported and tested go version

- 1.17
- 1.18 (1.18+ are tested with extra generic type tests)
- 1.19
- 1.20

You may see [compile error](https://github.com/trim21/go-phpserialize/issues/31) in very new golang version, but it's expected.

Due to the usage of unsafe, new version of golang may break this package, so it use go build flags to make it only compile on tested go versions.

## Use case:

You serialize all data into php array only. 

Decoding from php serialized array or class are both supported.

## Install

```console
go get github.com/trim21/go-phpserialize
```

### Advantage:

Low memory allocation and fast, see [benchmark](./docs/benchmark.md)

### Disadvantage:

heavy usage of `unsafe`.

## Usage

## Unmarshal

See [examples](./example_test.go)
`any` type will be decoded to `map[any]any` or `map[string]any`, depends on raw input is `array` or `class`,

map `any` key maybe `int64` or `string`.

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
