#ifndef __TC_WRAP__
#define __TC_WRAP__
#include <stdlib.h>

void TCRequest(char** keys, char** values, size_t map_size,
                char* buf, size_t buf_size);

#endif
