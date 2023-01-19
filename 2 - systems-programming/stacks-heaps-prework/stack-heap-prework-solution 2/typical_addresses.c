#include <pthread.h>
#include <stdio.h>
#include <stdlib.h>

char *s = "hello, world";
int x, y, z;

int factorial(int n) {
	printf("\t%p (n = %d)\n", &n, n);
	if (n < 2) {
		return 1;
	} else {
		return n * factorial(n - 1);
	}
}

void *f(void *ptr) {
	int a, b, c;
	printf("local variables in another thread:\n\t%p (int a)\n\t%p (int b)\n\t%p (int c)\n", &a, &b, &c);
	printf("parameters for recursive function calls in another thread:\n");
	factorial(5);

	// Block waiting on input, so we can inspect the program's address space
	// using command line tools
	scanf("%d\n", &a);
	return NULL;
}

int main() {
	char *s2 = "hello, world";
	char *s3 = "foo bar";
	int a, b, c;

	printf("string constants:\n\t%p (s = \"hello, world\")\n\t%p (s2 = \"hello, world\")\n\t%p (s3 = \"foo bar\")\n", s, s2, s3);
	printf("global variables:\n\t%p (char *s)\n\t%p (int x)\n\t%p (int y)\n\t%p (int z)\n", &s, &x, &y, &z);
	printf("local variables:\n\t%p (int a)\n\t%p (int b)\n\t%p (int c)\n", &a, &b, &c);
	printf("parameters for recursive function calls:\n");
	factorial(5);

	pthread_t t;
	pthread_create(&t, NULL, f, NULL);
	pthread_join(t, NULL);
}
