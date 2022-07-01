package main

import (
	"math"
	"runtime"
	"sync"
	"time"

	"github.com/trim21/go-phpserialize"
)

func main() {
	type T1 struct {
		F1 uint64 // 这里在怀疑会不会是 string 的问题。
	}
	type T2 struct {
		F1 T1 // 这里在怀疑是不是结构指针的问题。
	}
	type T3 struct {
		F1 []T2    // 这里在怀疑是不是 slice of struct 的问题。
		F2 []int64 // 这里怀疑是不是 slice of string 的问题。
	}
	makeT1 := func(n int) *T1 {
		return &T1{
			F1: 5, // 确保字符串都是新申请的内存。
		}
	}
	makeT2Slice := func(n int) []T2 {
		l := n%5 + 1 // 申请一个有长度但不太长的 slice。
		slice := make([]T2, 0, l)

		for i := 0; i < l; i++ {
			slice = append(slice, T2{*makeT1(n)})
		}

		return slice
	}
	makeT3 := func(n int) *T3 {
		return &T3{
			F1: makeT2Slice(n),
			F2: []int64{1, 2, 3, 4},
		}
	}

	max := runtime.GOMAXPROCS(0) // 尽量占用所有 CPU 资源从而造成竞争环境。
	var wg sync.WaitGroup
	wg.Add(max)

	for i := 0; i < max; i++ {
		go func() {
			defer wg.Done()

			for n := 0; n < math.MaxInt32; n++ {
				phpserialize.Marshal(makeT1(1))
				phpserialize.Marshal(makeT2Slice(10))
				phpserialize.Marshal(makeT3(2))
			}
		}()
	}

	done := make(chan bool, 1)
	ticker := time.NewTicker(time.Second)

	go func() {
		wg.Wait()
		done <- true
	}()

	for {
		select {
		case <-ticker.C:
			// 构造一个 1s 执行一次的定时器，强制触发 GC。
			runtime.GC()

		case <-done:
			return
		}
	}

}
