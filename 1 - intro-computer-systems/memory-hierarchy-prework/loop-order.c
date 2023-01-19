  /*

Two different ways to loop over an array of arrays.

Spotted at:
http://stackoverflow.com/questions/9936132/why-does-the-order-of-the-loops-affect-performance-when-iterating-over-a-2d-arra

*/

#include <time.h>
#include <stdio.h>
#include <unistd.h>

void option_one() {
  int i, j;
  static int x[4000][4000];
  for (i = 0; i < 4000; i++) {
      for (j = 0; j < 4000; j++) {
          x[i][j] = i + j;
      }
  }
  sleep(1);
}

void option_two() {
  int i, j;
  static int x[4000][4000];
  for (i = 0; i < 4000; i++) {
      for (j = 0; j < 4000; j++) {
          x[j][i] = i + j;
      }
  }
}

int main() {
    clock_t option1_start, option1_end, option2_end;
    option1_start = clock();
    option_one();
    option1_end = clock();
    option_two();
    option2_end = clock();

    double option1_elapsed = (option1_end - option1_start) / CLOCKS_PER_SEC;
    double option2_elapsed = (option2_end - option1_end) / CLOCKS_PER_SEC;

    printf("Option 1: %fs\n", option1_elapsed);
    printf("Option 2: %fs\n", option2_elapsed);

    return 0;
}
