// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License that can be found in
// the LICENSE file.

//go:generate go run asm_gen.go

package atomic128

import (
	"reflect"
	"unsafe"
)

// Uint128 represents a 128-bit unsigned integer.
type Uint128 [2]uint64

// DoublePointer represents two pointers.
type DoublePointer [2]unsafe.Pointer

// SwapUint128 atomically stores new into *addr and returns the previous *addr value.
func SwapUint128(addr *Uint128, new Uint128) (old Uint128)

// CompareAndSwapUint128 executes the compare-and-swap operation for a Uint128 value.
func CompareAndSwapUint128(addr *Uint128, old, new Uint128) (swapped bool)

// LoadUint128 atomically loads *addr.
func LoadUint128(addr *Uint128) (val Uint128)

// StoreUint128 atomically stores val into *addr.
func StoreUint128(addr *Uint128, val Uint128)

// SwapDoublePointer atomically stores new into *addr and returns the previous *addr value.
func SwapDoublePointer(addr *DoublePointer, new DoublePointer) (old DoublePointer)

// CompareAndSwapDoublePointer executes the compare-and-swap operation for a DoublePointer value.
func CompareAndSwapDoublePointer(addr *DoublePointer, old, new DoublePointer) (swapped bool)

// LoadDoublePointer atomically loads *addr.
func LoadDoublePointer(addr *DoublePointer) (val DoublePointer)

// StoreDoublePointer atomically stores val into *addr.
func StoreDoublePointer(addr *DoublePointer, val DoublePointer)

// SwapStringHeader atomically stores new into *addr and returns the previous *addr value.
func SwapStringHeader(addr *reflect.StringHeader, new reflect.StringHeader) (old reflect.StringHeader)

// CompareAndSwapStringHeader executes the compare-and-swap operation for a reflect.StringHeader value.
func CompareAndSwapStringHeader(addr *reflect.StringHeader, old, new reflect.StringHeader) (swapped bool)

// LoadStringHeader atomically loads *addr.
func LoadStringHeader(addr *reflect.StringHeader) (val reflect.StringHeader)

// StoreStringHeader atomically stores val into *addr.
func StoreStringHeader(addr *reflect.StringHeader, val reflect.StringHeader)

// SwapInterface atomically stores new into *addr and returns the previous *addr value.
func SwapInterface(addr *interface{}, new interface{}) (old interface{})

// CompareAndSwapInterface executes the compare-and-swap operation for an interface{} value.
func CompareAndSwapInterface(addr *interface{}, old, new interface{}) (swapped bool)

// LoadInterface atomically loads *addr.
func LoadInterface(addr *interface{}) (val interface{})

// StoreInterface atomically stores val into *addr.
func StoreInterface(addr *interface{}, val interface{})
