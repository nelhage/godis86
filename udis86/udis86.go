package udis86

//#include "udis86.h"
//#include "string.h"
//extern int ud_input_hook(ud_t *u);
import "C"

import (
	"io"
	"unsafe"
)

// a Dissassembler encapsulates a single instance of the
// disassembler. Multiple Disassemblers may be created and operated on
// in parallel.
type Disassembler struct {
	u     C.struct_ud
	bytes []byte
	r     io.Reader
}

type Vendor int

const (
	VendorAny   Vendor = C.UD_VENDOR_ANY
	VendorIntel Vendor = C.UD_VENDOR_INTEL
	VendorAMD   Vendor = C.UD_VENDOR_AMD
)

type Syntax int

const (
	SyntaxNone Syntax = iota
	SyntaxIntel
	SyntaxATT
)

type Config struct {
	Buf    []byte
	Reader io.Reader
	Bits   byte
	PC     uint64
	Vendor Vendor
	Syntax Syntax
}

//export inputHook
func inputHook(p unsafe.Pointer) C.int {
	d := (*Disassembler)(p)
	var b [1]byte
	if n, e := d.r.Read(b[:]); e != nil || n == 0 {
		return C.UD_EOI
	}
	return C.int(b[0])
}

func New(c *Config) *Disassembler {
	if c.Buf != nil && c.Reader != nil {
		panic("New: must pass either bytes or a Reader")
	}
	d := &Disassembler{
		bytes: c.Buf,
		r:     c.Reader,
	}
	C.ud_init(&d.u)
	C.ud_set_user_opaque_data(&d.u, unsafe.Pointer(d))
	if d.bytes != nil {
		C.ud_set_input_buffer(&d.u, (*C.uint8_t)(unsafe.Pointer(&d.bytes[0])), C.size_t(len(d.bytes)))
	} else {
		C.ud_set_input_hook(&d.u, (*[0]byte)(C.ud_input_hook))
	}
	if c.Bits != 0 {
		C.ud_set_mode(&d.u, C.uint8_t(c.Bits))
	}
	if c.PC != 0 {
		C.ud_set_pc(&d.u, C.uint64_t(c.PC))
	}
	// We rely on the fact that the zero value is also the default
	// to udis86, and don't check if it is set.
	C.ud_set_vendor(&d.u, C.unsigned(c.Vendor))
	switch c.Syntax {
	case SyntaxIntel:
		C.ud_set_syntax(&d.u, (*[0]byte)(C.UD_SYN_INTEL))
	case SyntaxATT:
		C.ud_set_syntax(&d.u, (*[0]byte)(C.UD_SYN_ATT))
	}
	return d
}

func (d *Disassembler) Next() bool {
	return C.ud_disassemble(&d.u) != 0
}

func (d *Disassembler) PC() uint64 {
	return uint64(C.ud_insn_off(&d.u))
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
