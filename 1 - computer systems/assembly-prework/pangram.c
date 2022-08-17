//
// Created by Nitin Savant on 8/17/22.
//
#include <stdio.h>

int pangram(char *s) {

    int result = 0;

    while (*s) {
        if (*s > 0b1100000 && *s < 0b1111011) {
            result |= 1 << (*s-0b1100001);
        }
        s++;
    }
    return result == 0b11111111111111111111111111;
}

int main(void) {
    printf("%d\n", pangram("abc"));
    printf("%d\n", pangram("the quick brown fox jumps over the lazy dog"));
}
