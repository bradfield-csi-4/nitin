#include <stdio.h>
#include <stdlib.h>

#define SHELL_CHAR '>'
#define INPUT_LEN 20

void interruptHandler(int);

int main(void) {

    signal(SIGINT, interruptHandler);

    system("stty intr ^D");
    char line[INPUT_LEN + 1];

    while (1) {
        printf("%c ", SHELL_CHAR);
        fgets(line, sizeof(line), stdin);
        fputs(line, stdout);
    }

}

void interruptHandler(int code) {
    printf("\nHave a good one!\n");
    exit(EXIT_SUCCESS);
}
