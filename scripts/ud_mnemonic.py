import os
import sys
import subprocess
from ud_opcode import UdOpcodeTable, UdOpcodeTables, UdInsnDef

class UdMnemonicTableGenerator:
    def __init__(self, tables):
        self.tables = tables

    def genTable(self, path):
        with open(path, 'w') as f:
            f.write("package udis86\n\n")
            f.write('//#include "udis86.h"' + "\n")
            f.write('import "C"' + "\n\n")
            f.write("const (\n")
            for mnemonic in self.tables.getMnemonicsList():
                f.write("\tI_%s Mnemonic = C.UD_I%s\n" % (mnemonic, mnemonic))
            f.write("\n)\n")
        subprocess.check_call(["gofmt", "-w", path])

    def generate(self, loc):
        self.genTable(os.path.join(loc, "mnemonics.go"))

def main():

    if len(sys.argv) != 3:
        usage()
        sys.exit(1)

    tables = UdOpcodeTables(xml=sys.argv[1])
    itab   = UdMnemonicTableGenerator(tables)
    itab.generate(sys.argv[2])

if __name__ == '__main__':
    main()
