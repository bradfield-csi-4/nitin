#include <stdlib.h>
#include <stdio.h>

int main(void) {
//    int b[10000000000000];

    int *b = malloc(100000000000000000000000000000000L);
    printf("%d", b);
}

// stack frame size (1345294384) exceeds limit (4294967295)
// 4294967295 --> 2^32 - 1