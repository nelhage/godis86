#include "_cgo_export.h"
#include "udis86.h"

int ud_input_hook(ud_t *u) {
    return inputHook(ud_get_user_opaque_data(u));
}
