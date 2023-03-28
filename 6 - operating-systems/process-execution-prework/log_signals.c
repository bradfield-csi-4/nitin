#include <stdio.h>
#include <unistd.h>
#include <signal.h>

/* Function Prototypes */
void handler(int);

int main() {
    int pid = (int) getpid();
    printf("Process started w/ PID: %d\n", pid);
    signal(SIGHUP, handler);

    for (int i = 1; i <= 30; i++) {
        signal(i, handler);
    }

    for (int i = 1; i <= 30; i++) {
        if (i == SIGKILL || i == SIGSTOP || i == SIGCONT) {
            continue;
        }
        kill(pid, i);
    }

    return;
}

void handler(int signum) {
    int n;
    char buf[28];
    n = sprintf(buf, "Received signal number: %d\n", signum);
    write(STDOUT_FILENO, buf, n);
    return;
}