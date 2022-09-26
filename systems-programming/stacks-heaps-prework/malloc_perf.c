#include <time.h>
#include <stdio.h>
#include <stdlib.h>

#define N 100000000

int main(void) {
    clock_t start, end, time_elapsed;

    start = clock();
    for (int i = 0; i < N; i++) {
        void *ptr = malloc(1);
        free(ptr);
    }
    end = clock();

    time_elapsed = (end - start) / CLOCKS_PER_SEC;

    printf("%.2lus to run %d iterations (%lu iterations / sec)\n", time_elapsed, N,
           N / time_elapsed);
}

// 03s to run 100000000 iterations (33333333 iterations / sec) -- sizeof(int)
