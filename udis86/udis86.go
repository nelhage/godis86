package udis86

//#include "udis86.h"
//#include "string.h"
//extern int ud_input_hook(ud_t *u);
import "C"

import (
	"io"
	"unsafe"
)

// A Disassembler encapsulates a single instance of the
// disassembler. Multiple Disassemblers may be created and operated on
// independently. A given Disassembler should only be used by one
// goroutine at once.
type Disassembler struct {
	u     C.struct_ud
	bytes []byte
	r     io.Reader
}

// Vendor describes which vendor's extensions should be decoded. This
// matters primarily for the VMX/SVM virtualization extensions.
type Vendor int

const (
	// VendorAny means "decode extensions from either vendor"
	VendorAny Vendor = C.UD_VENDOR_ANY
	// VendorIntel means "only decode extensions recognized by Intel CPUs"
	VendorIntel Vendor = C.UD_VENDOR_INTEL
	// VendorAMD means "only decode extensions recognized by AMD CPUs"
	VendorAMD Vendor = C.UD_VENDOR_AMD
)

// Syntax describes how to render the disassembled instructions to a
// textual format.
type Syntax int

const (
	// SyntaxNone means do not produce a string version
	SyntaxNone Syntax = iota
	// SyntaxIntel uses Intel syntax (e.g. no % prefixes,
	// [base+scale*off] syntax)
	SyntaxIntel
	// SyntaxATT uses AT&T assembler syntax. %- prefixes on
	// register, (%base,%index,scale) addressing, etc.
	SyntaxATT
)

// Config defines how to create a new Diasassembler. At a minimum,
// exactly one of Buf and Reader must be populated to define where to
// take input from. All other fields will be populated with defaults
// by udis86.
type Config struct {
	// Buf defines a byte slice to read assembly from.
	Buf []byte
	// Reader defines an io.Reader to read input from. godis86
	// always performs one-byte reads, so expensive Readers should
	// be wrapped in a bufio.Reader if possible.
	Reader io.Reader
	// Bits defines the disassembly mode -- 16-, 32-, or 64-bit.
	Bits byte
	// PC specifies the instruction pointer value (%ip/%eip/%rip)
	// of the start of the input. This affects rendering of
	// relative offsets.
	PC uint64
	// Vendor determines which vendor's extensions to accept
	Vendor Vendor
	// Syntax determines the syntax of output returned by
	// String().
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
