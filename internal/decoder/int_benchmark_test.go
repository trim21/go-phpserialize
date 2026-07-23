package decoder

import (
	"fmt"
	"strconv"
	"testing"
)

var (
	benchmarkIntResult   int64
	benchmarkIndexResult int
)

func BenchmarkDigitScan(b *testing.B) {
	inputs := []string{
		"147852369;",
		"9223372036854775807;",
	}

	for _, input := range inputs {
		buf := []byte(input)
		b.Run(fmt.Sprintf("table/%d_digits", len(buf)-1), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				cursor := 0
				for numTable[buf[cursor]] {
					cursor++
				}
				benchmarkIndexResult = cursor
			}
		})

		b.Run(fmt.Sprintf("range/%d_digits", len(buf)-1), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				cursor := 0
				for c := buf[cursor]; c >= '0' && c <= '9'; c = buf[cursor] {
					cursor++
				}
				benchmarkIndexResult = cursor
			}
		})
	}
}

func BenchmarkParseInt(b *testing.B) {
	inputs := []string{
		"7",
		"147852369",
		"-957123587",
		"9223372036854775807",
		"-9223372036854775808",
	}

	for _, input := range inputs {
		buf := []byte(input)
		b.Run("pow10/"+input, func(b *testing.B) {
			b.ReportAllocs()
			var result int64
			var err error
			for b.Loop() {
				result, err = parseIntPow10(buf)
			}
			if err != nil {
				b.Fatal(err)
			}
			benchmarkIntResult = result
		})

		b.Run("strconv/"+input, func(b *testing.B) {
			b.ReportAllocs()
			var result int64
			var err error
			for b.Loop() {
				result, err = strconv.ParseInt(unsafeStr(buf), 10, 64)
			}
			if err != nil {
				b.Fatal(err)
			}
			benchmarkIntResult = result
		})

		b.Run("checked/"+input, func(b *testing.B) {
			b.ReportAllocs()
			var result int64
			var err error
			for b.Loop() {
				result, err = parseIntChecked(buf)
			}
			if err != nil {
				b.Fatal(err)
			}
			benchmarkIntResult = result
		})

		b.Run("pow10_checked/"+input, func(b *testing.B) {
			b.ReportAllocs()
			var result int64
			var err error
			for b.Loop() {
				result, err = parseInt64(buf)
			}
			if err != nil {
				b.Fatal(err)
			}
			benchmarkIntResult = result
		})
	}
}

var benchmarkPow10i64 = [...]int64{
	1e00, 1e01, 1e02, 1e03, 1e04, 1e05, 1e06, 1e07, 1e08, 1e09,
	1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16, 1e17, 1e18,
}

func parseIntPow10(buf []byte) (int64, error) {
	negative := false
	if buf[0] == '-' {
		buf = buf[1:]
		negative = true
	}
	maxDigit := len(buf)
	if maxDigit > len(benchmarkPow10i64) {
		return 0, fmt.Errorf("invalid length of number")
	}
	var sum int64
	for i := range maxDigit {
		digit := int64(buf[i]) - '0'
		sum += digit * benchmarkPow10i64[maxDigit-i-1]
	}
	if negative {
		return -sum, nil
	}
	return sum, nil
}

func parseIntChecked(buf []byte) (int64, error) {
	if len(buf) == 0 {
		return 0, fmt.Errorf("invalid integer")
	}

	negative := buf[0] == '-'
	if negative {
		buf = buf[1:]
		if len(buf) == 0 {
			return 0, fmt.Errorf("invalid integer")
		}
	}

	limit := uint64(^uint64(0) >> 1)
	if negative {
		limit++
	}

	var value uint64
	for _, c := range buf {
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("invalid integer")
		}
		digit := uint64(c - '0')
		if value > (limit-digit)/10 {
			return 0, fmt.Errorf("integer overflow")
		}
		value = value*10 + digit
	}

	if negative {
		if value == uint64(1)<<63 {
			return -1 << 63, nil
		}
		return -int64(value), nil
	}
	return int64(value), nil
}
