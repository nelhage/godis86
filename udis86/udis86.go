package udis86

//#cgo CFLAGS: -I..
//#cgo LDFLAGS: libudis86.a
//#include "udis86.h"
//#include "string.h"
import "C"

import "unsafe"

// a Dissassembler encapsulates a single instance of the
// disassembler. Multiple Disassemblers may be created and operated on
// in parallel.
type Disassembler struct {
	u     C.struct_ud
	bytes []byte
}

type Vendor int

const (
	VendorAny   = C.UD_VENDOR_ANY
	VendorIntel = C.UD_VENDOR_INTEL
	VendorAMD   = C.UD_VENDOR_AMD
)

type Syntax int

const (
	SyntaxIntel = iota
	SyntaxATT
)

type Options struct {
	Bits   byte
	PC     uint64
	Vendor Vendor
	Syntax Syntax
}

func New(b []byte) *Disassembler {
	d := &Disassembler{
		bytes: b,
	}
	C.ud_init(&d.u)
	C.ud_set_input_buffer(&d.u, (*C.uint8_t)(unsafe.Pointer(&d.bytes[0])), C.size_t(len(d.bytes)))
	return d
}

func (d *Disassembler) Next() bool {
	return C.ud_disassemble(&d.u) != 0
}

func (d *Disassembler) Len() int {
	return int(C.ud_insn_len(&d.u))
}

func (d *Disassembler) Bytes() []byte {
	b := make([]byte, d.Len())
	C.memcpy(unsafe.Pointer(&b[0]),
		unsafe.Pointer(C.ud_insn_ptr(&d.u)),
		C.size_t(d.Len()))
	return b
}

func (d *Disassembler) String() string {
	return C.GoString(C.ud_insn_asm(&d.u))
}

type Mnemonic uint16

func (m Mnemonic) String() string {
	return C.GoString(C.ud_lookup_mnemonic(uint16(m)))
}

func (d *Disassembler) Mnemonic() Mnemonic {
	return Mnemonic(C.ud_insn_mnemonic(&d.u))
}
