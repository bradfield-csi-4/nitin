#include <stdio.h>
#include <stdlib.h>
#include <sys/time.h>

int main(int argc, char *argv[]) {
  if (argc != 2) {
    printf("Usage: %s [n_allocs]\n", argv[0]);
    exit(1);
  }

  int n_allocs = atoi(argv[1]);
  struct timeval start, stop;
  double elapsed;

  gettimeofday(&start, NULL);
  for (int i = 0; i < n_allocs; i++) {
	  int *ptr = malloc((rand() % 32 + 1) * sizeof(int));
	  if (ptr == NULL) {
		  fprintf(stderr, "malloc failed on iteration %d\n", i);
		  exit(1);
	  }
  }
  gettimeofday(&stop, NULL);

  elapsed = (stop.tv_sec - start.tv_sec) * 1000000 + stop.tv_usec - start.tv_usec;
  printf("Total elapsed time: %.3f milliseconds\n", elapsed / 1000);
  printf("Average time per allocation: %.3f nanoseconds\n", elapsed * 1000 / n_allocs);
}
