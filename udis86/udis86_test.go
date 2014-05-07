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
	bytes := d.Bytes()
	if !reflect.DeepEqual(bytes, []byte{0x90}) {
		t.Errorf("Bytes() returned %v, wanted %v", bytes, []byte{0x90})
	}

	m := d.Mnemonic()
	if m != I_nop {
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
