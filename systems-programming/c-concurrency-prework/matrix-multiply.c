#include <stdio.h>
#include <stdlib.h>
#include <pthread.h>
#include "csapp.h"

#define NUM_THREADS 8

void *thread(void *vargp);  /* Thread routine prototype */

void print(double **c, int row, int col) {
    int i, j;
    for (i = 0; i < row; i++) {
        for (j = 0; j < col; j++) {
            c[i][j] = '.';
            printf("%f ", c[i][j]);
        }
        printf("\n");
    }
    printf("\n");
}

typedef struct matrix_params {
    int num;
    double **C;
    double **A;
    double **B;
    int a_rows;
    int a_cols;
    int b_cols;
} matrix_params;

/*
  A naive implementation of matrix multiplication.

  DO NOT MODIFY THIS FUNCTION, the tests assume it works correctly, which it
  currently does
*/
void matrix_multiply(double **C, double **A, double **B, int a_rows, int a_cols,int b_cols) {
  for (int i = 0; i < a_rows; i++) {
    for (int j = 0; j < b_cols; j++) {
      C[i][j] = 0;
      for (int k = 0; k < a_cols; k++)
        C[i][j] += A[i][k] * B[k][j];
    }
  }
}

void parallel_matrix_multiply(double **c, double **a, double **b, int a_rows, int a_cols, int b_cols) {

    pthread_t threads[NUM_THREADS];
    matrix_params params_list[NUM_THREADS];

    for (int i = 0; i < NUM_THREADS; i++) {
        matrix_params params = {i, c, a, b, a_rows, a_cols, b_cols};
        params_list[i] = params;
        Pthread_create(&threads[i], NULL, thread, &params_list[i]);
    }

    for (int i = 0; i < NUM_THREADS; i++)
        Pthread_join(threads[i], NULL);
}

void *thread(void *vargp){
    matrix_params *m = (matrix_params *) vargp;

    int block_size = m->a_rows / NUM_THREADS;
    int start = m->num * block_size;
    int end = start + block_size;

    for (int i = start; i < end; i++) {
        for (int j = 0; j < m->b_cols; j++) {
            m->C[i][j] = 0;
            for (int k = 0; k < m->a_cols; k++)
                m->C[i][j] += m->A[i][k] * m->B[k][j];
        }
    }
    return NULL;
}