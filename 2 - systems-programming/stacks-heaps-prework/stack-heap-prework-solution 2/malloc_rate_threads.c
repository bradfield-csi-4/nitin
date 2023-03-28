#include <pthread.h>
#include <stdio.h>``
#include <stdlib.h>
#include <time.h>

void *run_benchmark(void *n_allocs) {
  for (int i = 0; i < *(int *)n_allocs; i++) {
	  int *ptr = malloc((rand() % 32 + 1) * sizeof(int));
	  if (ptr == NULL) {
		  fprintf(stderr, "malloc failed on iteration %d\n", i);
		  exit(1);
	  }
  }
  return NULL;
}

int main(int argc, char *argv[]) {
  if (argc != 3) {
    printf("Usage: %s [n_allocs] [n_threads]\n", argv[0]);
    exit(1);
  }

  int n_allocs = atoi(argv[1]);
  int n_threads = atoi(argv[2]);

  pthread_t *t = malloc(n_threads * sizeof(pthread_t));

  clock_t start, stop;
  double elapsed;

  start = clock();
  for (int i = 0; i < n_threads; i++) {
	  if (pthread_create(&t[i], NULL, run_benchmark, &n_allocs) != 0) {
		  fprintf(stderr, "pthread_create failed on iteration %d\n", i);
		  exit(1);
	  }
  }
  for (int i = 0; i < n_threads; i++) {
	  void *elapsed;
	  if (pthread_join(t[i], NULL) != 0) {
		  fprintf(stderr, "pthread_join failed on iteration %d\n", i);
		  exit(1);
	  }
  }
  stop = clock();

  elapsed = (stop - start) / (double) CLOCKS_PER_SEC;

  printf("Total elapsed CPU time (all threads): %.3f milliseconds\n", elapsed * 1e3);
  printf("Average time per allocation: %.3f nanoseconds\n", elapsed / (n_allocs * n_threads) * 1e9);
  free(t);
}
