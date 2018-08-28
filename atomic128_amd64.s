// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License that can be found in
// the LICENSE file.

// +build amd64,!gccgo,!appengine

#include "textflag.h"

TEXT ·SwapUint128(SB),NOSPLIT,$0
	MOVQ addr+0(FP), BP
	XORQ AX, AX
	XORQ DX, DX
	MOVQ new+8(FP), BX
	MOVQ new+16(FP), CX
loop:
	LOCK
	CMPXCHG16B (BP)
	JNE loop
	MOVQ AX, old+24(FP)
	MOVQ DX, old+32(FP)
	RET

TEXT ·CompareAndSwapUint128(SB),NOSPLIT,$0
	MOVQ addr+0(FP), BP
	MOVQ old+8(FP), AX
	MOVQ old+16(FP), DX
	MOVQ new+24(FP), BX
	MOVQ new+32(FP), CX
	LOCK
	CMPXCHG16B (BP)
	SETEQ swapped+40(FP)
	RET

TEXT ·LoadUint128(SB),NOSPLIT,$0
	MOVQ addr+0(FP), BP
	XORQ AX, AX
	XORQ DX, DX
	XORQ BX, BX
	XORQ CX, CX
	LOCK
	CMPXCHG16B (BP)
	MOVQ AX, val+8(FP)
	MOVQ DX, val+16(FP)
	RET

TEXT ·StoreUint128(SB),NOSPLIT,$0
	MOVQ addr+0(FP), BP
	XORQ AX, AX
	XORQ DX, DX
	MOVQ new+8(FP), BX
	MOVQ new+16(FP), CX
loop:
	LOCK
	CMPXCHG16B (BP)
	JNE loop
	RET

TEXT ·SwapDoublePointer(SB),NOSPLIT,$0
	JMP ·SwapUint128(SB)

TEXT ·CompareAndSwapDoublePointer(SB),NOSPLIT,$0
	JMP ·CompareAndSwapUint128(SB)

TEXT ·LoadDoublePointer(SB),NOSPLIT,$0
	JMP ·LoadUint128(SB)

TEXT ·StoreDoublePointer(SB),NOSPLIT,$0
	JMP ·StoreUint128(SB)

TEXT ·SwapStringHeader(SB),NOSPLIT,$0
	JMP ·SwapUint128(SB)

TEXT ·CompareAndSwapStringHeader(SB),NOSPLIT,$0
	JMP ·CompareAndSwapUint128(SB)

TEXT ·LoadStringHeader(SB),NOSPLIT,$0
	JMP ·LoadUint128(SB)

TEXT ·StoreStringHeader(SB),NOSPLIT,$0
	JMP ·StoreUint128(SB)

TEXT ·SwapInterface(SB),NOSPLIT,$0
	JMP ·SwapUint128(SB)

TEXT ·CompareAndSwapInterface(SB),NOSPLIT,$0
	JMP ·CompareAndSwapUint128(SB)

TEXT ·LoadInterface(SB),NOSPLIT,$0
	JMP ·LoadUint128(SB)

TEXT ·StoreInterface(SB),NOSPLIT,$0
	JMP ·StoreUint128(SB)
