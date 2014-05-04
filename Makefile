CPPFLAGS=-DHAVE_STRING_H

OBJECTS = \
	libudis86/itab.o \
	libudis86/decode.o \
	libudis86/syn.o \
	libudis86/syn-intel.o \
	libudis86/syn-att.o \
	libudis86/udis86.o \

libudis86.a: $(OBJECTS)
	$(AR) rc $@ $^

PYTHON  = python
OPTABLE = docs/x86/optable.xml

libudis86/itab.c libudis86/itab.h: $(OPTABLE) \
               scripts/ud_itab.py \
               scripts/ud_opcode.py
	$(PYTHON) scripts/ud_itab.py $(OPTABLE) libudis86

clean:
	rm -f $(OBJECTS) libudis86.a
