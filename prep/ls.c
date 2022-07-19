#include <stdio.h>
#include <dirent.h>
#include <stdlib.h>
#include <unistd.h>
#include <fcntl.h>
#include <sys/stat.h>

#ifndef MAX_BUF
#define MAX_BUF 100
#endif

int main(void) {

    // Opens current directory
    DIR *dirp = opendir(".");
    if (dirp == NULL) {
        perror("opendir() error");
        exit(EXIT_FAILURE);
    }

    // Stores entries in an array
    int i = 0;
    struct dirent *entryp;
    struct dirent *entries[MAX_BUF];

    while ((entryp = readdir(dirp)) != NULL)
        entries[i++] = entryp;

    // Closes opened directory
    closedir(dirp);

    // Prints file info
    printf("%-20s %-10s %-s\n\n", "Filename", "Size", "Last Modified Time");

    struct stat sp;
    int fd;

    for (int j = 0; j < i; j++) {
        char *filename = entries[j]->d_name;

        // Skips hidden files/commands
        if (filename[0] == '.')
            continue;

        fd = open(filename, O_RDONLY);

        if (fstat(fd, &sp) == -1) {
            perror("fstat() error");
            exit(EXIT_FAILURE);
        }

        printf("%-20s %-10lld %-ld\n", filename, sp.st_size, sp.st_mtimespec.tv_sec);

        close(fd);
    }
}
