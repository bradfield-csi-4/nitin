#include <stdio.h>
#include <stdlib.h>
#include <pthread.h>

/*
  A naive implementation of matrix multiplication.

  DO NOT MODIFY THIS FUNCTION, the tests assume it works correctly, which it
  currently does
*/
void matrix_multiply(double **C, double **A, double **B, int a_rows, int a_cols,
                     int b_cols) {
  for (int i = 0; i < a_rows; i++) {
    for (int j = 0; j < b_cols; j++) {
      C[i][j] = 0;
      for (int k = 0; k < a_cols; k++)
        C[i][j] += A[i][k] * B[k][j];
    }
  }
}

typedef struct multiply_params {
	double **A;
	double **B;
	double **C;
	int a_rows;
	int a_cols;
	int b_cols;
	int i0;
	int i1;
} multiply_params;

void *multiply_row(void *ptr) {
	multiply_params *p = (multiply_params *) ptr;
	for (int i = p->i0; i < p->i1; i++) {
		for (int j = 0; j < p->b_cols; j++) {
			p->C[i][j] = 0;
			for (int k = 0; k < p->a_cols; k++) {
				p->C[i][j] += p->A[i][k] * p->B[k][j];
			}
		}
	}
	return NULL;
}

const int NUM_THREADS = 4;

// Launch NUM_THREADS threads to compute the answer in parallel; each thread
// computes the answer for (approximately) `a_rows / NUM_THREADS` rows of the
// final result.
void parallel_matrix_multiply(double **c, double **a, double **b, int a_rows,
                          int a_cols, int b_cols) {
	multiply_params common_params = {
		a, b, c, a_rows, a_cols, b_cols
	};
	multiply_params params[NUM_THREADS];
	pthread_t threads[NUM_THREADS];
	for (int t = 0; t < NUM_THREADS; t++) {
		params[t] = common_params;
		params[t].i0 = a_rows * t / NUM_THREADS;
		params[t].i1 = a_rows * (t + 1) / NUM_THREADS;
		pthread_create(&threads[t], NULL, multiply_row, &params[t]);
    }
	for (int t = 0; t < NUM_THREADS; t++) {
		pthread_join(threads[t], NULL);
	}
}
