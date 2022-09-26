#include <limits.h>
#include <stdio.h>
#include <stdlib.h>

// 10^8 works, but 10^9 causes a crash
int global_array[100000000];

int main() {
	// 10^6 works, but 10^7 causes a crash
	int stack_array[1000000];

	// Even 10^9 works
	int *heap_array = malloc(1000000000UL * sizeof(int));
	if (heap_array == NULL) {
		fprintf(stderr, "malloc failed\n");
	}
}
