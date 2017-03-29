// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License that can be found in
// the LICENSE file.

package atomic128

import "unsafe"

// A Value provides an atomic load and store of an interface{}.
// Values can be created as part of other data structures.
// The zero value for a Value returns nil from Load.
// Once Store has been called, a Value must not be copied.
//
// A Value must not be copied after first use.
type Value struct {
	noCopy noCopy
	v      interface{}
}

// Load returns the value set by the most recent Store.
// It returns nil if there has been no call to Store for this Value.
func (v *Value) Load() (x interface{}) {
	ip := (*interface{})(unsafe.Pointer(v))
	return LoadInterface(ip)
}

// Store sets the value of the Value to x.
func (v *Value) Store(x interface{}) {
	ip := (*interface{})(unsafe.Pointer(v))
	StoreInterface(ip, x)
}

// noCopy may be embedded into structs which must not be copied
// after the first use.
//
// See https://github.com/golang/go/issues/8005#issuecomment-190753527
// for details.
type noCopy struct{}

// Lock is a no-op used by -copylocks checker from `go vet`.
func (*noCopy) Lock() {}
