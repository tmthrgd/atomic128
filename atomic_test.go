// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package atomic128

import (
	"fmt"
	"runtime"
	"strings"
	"testing"
	"unsafe"
)

// Tests of correct behavior, without contention.
// (Does the function work as advertised?)
//
// Test that the Add functions add correctly.
// Test that the CompareAndSwap functions actually
// do the comparison and the swap correctly.
//
// The loop over power-of-two values is meant to
// ensure that the operations apply to the full word size.
// The struct fields x.before and x.after check that the
// operations do not extend past the full word size.

var magic128 = Uint128{0xdeddeadbeefbeef, 0xdeddeadbeefbeef}

func TestSwapUint128(t *testing.T) {
	var x struct {
		before Uint128
		i      Uint128
		after  Uint128
	}
	x.before = magic128
	x.after = magic128
	var j Uint128
	for delta := uint64(1); delta+delta > delta; delta += delta {
		d := Uint128{delta, ^delta}
		k := SwapUint128(&x.i, d)
		if x.i != d || k != j {
			t.Fatalf("delta=%d i=%d j=%d k=%d", d, x.i, j, k)
		}
		j = d
	}
	if x.before != magic128 || x.after != magic128 {
		t.Fatalf("wrong magic: %#x _ %#x != %#x _ %#x", x.before, x.after, magic128, magic128)
	}
}

func TestCompareAndSwapUint128(t *testing.T) {
	var x struct {
		before Uint128
		i      Uint128
		after  Uint128
	}
	x.before = magic128
	x.after = magic128
	for val := uint64(1); val+val > val; val += val {
		x.i = Uint128{val, ^val}
		val1 := Uint128{(val + 1), ^(val + 1)}
		if !CompareAndSwapUint128(&x.i, Uint128{val, ^val}, val1) {
			t.Fatalf("should have swapped %#x %#x", Uint128{val, ^val}, val1)
		}
		if x.i != val1 {
			t.Fatalf("wrong x.i after swap: x.i=%#x val+1=%#x", x.i, val1)
		}
		x.i = val1
		if CompareAndSwapUint128(&x.i, Uint128{val, ^val}, Uint128{(val + 2), ^(val + 2)}) {
			t.Fatalf("should not have swapped %#x %#x", Uint128{val, ^val}, Uint128{(val + 2), ^(val + 2)})
		}
		if x.i != val1 {
			t.Fatalf("wrong x.i after swap: x.i=%#x val+1=%#x", x.i, val1)
		}
	}
	if x.before != magic128 || x.after != magic128 {
		t.Fatalf("wrong magic: %#x _ %#x != %#x _ %#x", x.before, x.after, magic128, magic128)
	}
}

func TestLoadUint128(t *testing.T) {
	var x struct {
		before Uint128
		i      Uint128
		after  Uint128
	}
	x.before = magic128
	x.after = magic128
	for delta := uint64(1); delta+delta > delta; delta += delta {
		k := LoadUint128(&x.i)
		if k != x.i {
			t.Fatalf("delta=%d i=%d k=%d", delta, x.i, k)
		}
		x.i[0] += delta
		x.i[1] -= delta
	}
	if x.before != magic128 || x.after != magic128 {
		t.Fatalf("wrong magic: %#x _ %#x != %#x _ %#x", x.before, x.after, magic128, magic128)
	}
}

func TestStoreUint128(t *testing.T) {
	var x struct {
		before Uint128
		i      Uint128
		after  Uint128
	}
	x.before = magic128
	x.after = magic128
	var v Uint128
	for delta := uint64(1); delta+delta > delta; delta += delta {
		StoreUint128(&x.i, v)
		if x.i != v {
			t.Fatalf("delta=%d i=%d v=%d", delta, x.i, v)
		}
		v[0] += delta
		v[1] -= delta
	}
	if x.before != magic128 || x.after != magic128 {
		t.Fatalf("wrong magic: %#x _ %#x != %#x _ %#x", x.before, x.after, magic128, magic128)
	}
}

// Tests of correct behavior, with contention.
// (Is the function atomic?)
//
// For each function, we write a "hammer" function that repeatedly
// uses the atomic operation to add 1 to a value. After running
// multiple hammers in parallel, check that we end with the correct
// total.
// Swap can't add 1, so it uses a different scheme.
// The functions repeatedly generate a pseudo-random number such that
// low bits are equal to high bits, swap, check that the old value
// has low and high bits equal.

var hammer128 = map[string]func(*Uint128, int){
	"SwapUint128":           hammerSwapUint128,
	"CompareAndSwapUint128": hammerCompareAndSwapUint128,
}

func hammerSwapUint128(addr *Uint128, count int) {
	seed := uint64(uintptr(unsafe.Pointer(&count)))
	for i := uint64(0); i < uint64(count); i++ {
		new := Uint128{seed + i, seed + i}
		old := SwapUint128(addr, new)
		if old[0] != old[1] {
			panic(fmt.Sprintf("SwapUint128 is not atomic: %v", old))
		}
	}
}

func hammerCompareAndSwapUint128(addr *Uint128, count int) {
	for i := 0; i < count; i++ {
		for {
			v := LoadUint128(addr)
			if CompareAndSwapUint128(addr, v, Uint128{v[0] + 1, v[1] - 1}) {
				break
			}
		}
	}
}

func TestHammer128(t *testing.T) {
	const p = 4
	n := 100000
	if testing.Short() {
		n = 1000
	}
	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(p))

	for name, testf := range hammer128 {
		c := make(chan int)
		var val Uint128
		for i := 0; i < p; i++ {
			go func() {
				defer func() {
					if err := recover(); err != nil {
						t.Error(err.(string))
					}
					c <- 1
				}()
				testf(&val, n)
			}()
		}
		for i := 0; i < p; i++ {
			<-c
		}
		exp := Uint128{uint64(n) * p, -uint64(n) * p}
		if !strings.HasPrefix(name, "Swap") && val != exp {
			t.Fatalf("%s: val=%v want %v", name, val, exp)
		}
	}
}

func hammerStoreLoadUint128(t *testing.T, paddr unsafe.Pointer) {
	addr := (*Uint128)(paddr)
	v := LoadUint128(addr)
	if v[0] != v[1] {
		t.Fatalf("Uint128: %#x != %#x", v[0], v[1])
	}
	new := Uint128{v[0] + 1, v[1] + 1}
	StoreUint128(addr, new)
}

func TestHammerStoreLoad(t *testing.T) {
	var tests []func(*testing.T, unsafe.Pointer)
	tests = append(tests, hammerStoreLoadUint128)
	n := int(1e6)
	if testing.Short() {
		n = int(1e4)
	}
	const procs = 8
	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(procs))
	for _, tt := range tests {
		c := make(chan int)
		var val uint64
		for p := 0; p < procs; p++ {
			go func() {
				for i := 0; i < n; i++ {
					tt(t, unsafe.Pointer(&val))
				}
				c <- 1
			}()
		}
		for p := 0; p < procs; p++ {
			<-c
		}
	}
}

func TestStoreLoadSeqCst64(t *testing.T) {
	if runtime.NumCPU() == 1 {
		t.Skipf("Skipping test on %v processor machine", runtime.NumCPU())
	}
	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(4))
	N := int64(1e3)
	if testing.Short() {
		N = int64(1e2)
	}
	c := make(chan bool, 2)
	X := [2]Uint128{}
	var m = Uint128{^uint64(0), ^uint64(0)}
	ack := [2][3]Uint128{{m, m, m}, {m, m, m}}
	for p := 0; p < 2; p++ {
		go func(me int) {
			he := 1 - me
			for i := uint64(1); i < uint64(N); i++ {
				ii := Uint128{i, i}
				StoreUint128(&X[me], ii)
				my := LoadUint128(&X[he])
				StoreUint128(&ack[me][i%3], my)
				for w := 1; LoadUint128(&ack[he][i%3]) == m; w++ {
					if w%1000 == 0 {
						runtime.Gosched()
					}
				}
				his := LoadUint128(&ack[he][i%3])
				ii1 := Uint128{i - 1, i - 1}
				if (my != ii && my != ii1) || (his != ii && his != ii1) {
					t.Errorf("invalid values: %d/%d (%d)", my, his, i)
					break
				}
				if my != ii && his != ii {
					t.Errorf("store/load are not sequentially consistent: %d/%d (%d)", my, his, i)
					break
				}
				StoreUint128(&ack[me][(i-1)%3], m)
			}
			c <- true
		}(p)
	}
	<-c
	<-c
}

func TestNilDeref(t *testing.T) {
	funcs := [...]func(){
		func() { CompareAndSwapUint128(nil, Uint128{}, Uint128{}) },
		func() { SwapUint128(nil, Uint128{}) },
		func() { LoadUint128(nil) },
		func() { StoreUint128(nil, Uint128{}) },
	}
	for _, f := range funcs {
		func() {
			defer func() {
				runtime.GC()
				recover()
			}()
			f()
		}()
	}
}
