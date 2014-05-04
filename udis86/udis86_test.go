package udis86

import (
	"reflect"
	"testing"
)

func TestBasic(t *testing.T) {
	d := New([]byte{0x90})
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
