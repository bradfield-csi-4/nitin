#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

int a = 1, b;
static int c = 1, d;

int main () {
  int e = 1, f;
  static int g = 2, h;
  int *i = malloc(10 * sizeof(int));
  int *j = malloc(10 * sizeof(int));


  printf("name    location\n");
  printf("\n");

  printf("TEXT\n");
  printf("main:   %p\n", main);
  printf("\n");

  printf("DATA: initialized globals/statics\n");
  printf("&a:     %p\n", &a);
  printf("&c:     %p\n", &c);
  printf("&g:     %p\n", &g);
  printf("\n");

  printf("BSS: Uninitialized globals/statics\n");
  printf("&d:     %p\n", &d);
  printf("&h:     %p\n", &h);
  printf("&b:     %p\n", &b);
  printf("\n");

  printf("HEAP\n");
  printf("i:      %p\n", i);
  printf("i[1]:   %p\n", &(i[1]));
  printf("i[9]:   %p\n", &(i[9]));
  printf("j:      %p\n", j);
  // printf("j[1]:   %p\n", &(j[1]));
  // printf("j[9]:   %p\n", &(j[9]));
  printf("\n");

  printf("SHARED LIBRARIES\n");
  printf("printf: %p\n", printf);
  printf("\n");

  printf("STACK\n");
  printf("&e:     %p\n", &e);
  printf("&f:     %p\n", &f);
  printf("&i:     %p\n", &i);
  printf("&j:     %p\n", &j);

  printf("\npid:    %d\n", getpid());
  // while (1)
  //   ;
}
