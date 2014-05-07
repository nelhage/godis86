package udis86

import (
	"bytes"
	"reflect"
	"testing"
)

func TestBasic(t *testing.T) {
	d := New(&Config{Buf: []byte{0x90}})
	if !d.Next() {
		t.Fatal("Disassemble() failed!")
	}
	if d.Len() != 1 {
		t.Errorf("Len() returned %d, not 1", d.Len())
	}

	if bytes := d.Bytes(); !reflect.DeepEqual(bytes, []byte{0x90}) {
		t.Errorf("Bytes() returned %v, wanted %v", bytes, []byte{0x90})
	}

	if m := d.Mnemonic(); m != I_nop {
		t.Errorf("mnemonic was %d (%s), want I_nop!", m, m.String())
	}
}

func TestReader(t *testing.T) {
	r := bytes.NewBuffer([]byte{0x90})
	d := New(&Config{Reader: r})
	if !d.Next() {
		t.Fatal("Disassemble() failed!")
	}
	if d.Len() != 1 {
		t.Errorf("Len() returned %d, not 1", d.Len())
	}
	if d.Next() {
		t.Errorf("Next() returned true on a second call!")
	}
}

func TestSyntax(t *testing.T) {
	d := New(&Config{Buf: []byte{0x48, 0xff, 0xc0}, Syntax: SyntaxATT, Bits: 64})
	d.Next()
	if asm := d.String(); asm != "inc %rax" {
		t.Errorf("bad ATT-syntax disassembly: %s", asm)
	}

	d = New(&Config{Buf: []byte{0x48, 0xff, 0xc0}, Syntax: SyntaxIntel, Bits: 64})
	d.Next()
	if asm := d.String(); asm != "inc rax" {
		t.Errorf("bad Intel-syntax disassembly: %s", asm)
	}
}
