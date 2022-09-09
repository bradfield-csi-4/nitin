//
// Created by Nitin Savant on 8/16/22.
//
#include <stdio.h>

int binary_convert(char *s) {
    int len = 0;

    while (*s) {
        s++;
        len++;
    }

    s--;
    int mult = 1;
    int result = 0;

    while (len > 0) {
        if (*s == '1') {
            result += mult;
        }
        mult *= 2;
        s--;
        len--;
    }

    return result;
}

int main(void) {
    printf("%d\n", binary_convert("111"));
}
