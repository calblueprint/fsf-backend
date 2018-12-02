#include "tc_wrap.h"

#include "tclink.h"

void TCRequest(char** keys, char** values, size_t map_size,
        char* buf, size_t buf_size) {
    TCLinkHandle handle = TCLinkCreate();

    for (size_t i = 0; i < map_size; i++) {
        char* key = keys[i];
        char* value = values[i];
        TCLinkPushParam(handle, key, value);
    }

    TCLinkSend(handle);

    TCLinkGetEntireResponse(handle, buf, buf_size);
}
