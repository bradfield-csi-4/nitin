#include <stdio.h>
#include <dirent.h>
#include <stdlib.h>
#include <unistd.h>
#include <fcntl.h>
#include <sys/stat.h>

#ifndef MAX_BUF
#define MAX_BUF 200
#endif

int main(void) {

    DIR *dir;

    // Opens current directory
    if ((dir = opendir(".")) == NULL) {
        perror("opendir() error");
        exit(EXIT_FAILURE);
    }

    // Stores directory entries in an array
    int i = 0;
    struct dirent *entry;
    struct dirent *entries[MAX_BUF];

    while ((entry = readdir(dir)) != NULL)
        entries[i++] = entry;

    // Closes directory
    closedir(dir);

    // Prints file info
    printf("%-20s %-10s %-s\n\n", "Filename", "Size", "Last Modified");

    struct stat *sp = malloc(sizeof(struct stat));
    int fd;

    for (int j = 0; j < i; j++) {
        char *filename = entries[j]->d_name;

        // Skips hidden files
        if (filename[0] == '.')
            continue;

        fd = open(filename, O_RDONLY);

        if (fstat(fd, sp) == -1) {
            perror("fstat() error");
            exit(EXIT_FAILURE);
        }

        printf("%-20s %-10lld %-ld\n", filename, sp->st_size, sp->st_mtimespec.tv_sec);

        close(fd);
    }

    free(sp);
}

