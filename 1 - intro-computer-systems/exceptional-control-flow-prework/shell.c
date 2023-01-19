#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <string.h>
#include <unistd.h>
#include <errno.h>

#define SHELL_PROMPT "> "
#define MAX_INPUT 20
#define MAX_ARGS 5
#define MAX_CMDS 2

/* Suppresses "unused parameter" warning */
#define UNUSED(x) (void)(x)

/* Function Prototypes */
void sig_int_handler(int);
void unix_error(char *msg);
pid_t Fork(void);
int split(char *input, char **argv, char *delim);
void eval(char *cmdline);
void builtin_command(char **argv);
void exit_program();
void repl();

int main(int argc, char *argv[]) {
    /* Register Signal Handlers (for parent process only) */
    signal(SIGINT, sig_int_handler);
//    signal(SIGTSTP, NULL);

    /* Pass command directly using "-c" flag (for testing) */
    if (argc == 3 && strcmp(argv[1], "-c") == 0) {
        char *command = argv[2];
        eval(command);
        exit(EXIT_SUCCESS);
    }

    repl();
}

void repl() {
    char cmdline[MAX_INPUT + 1];

    /* LOOP */
    while (true) {
        printf(SHELL_PROMPT);

        /* READ */
        fgets(cmdline, MAX_INPUT, stdin);
        if (feof(stdin))
            exit_program();

        /* EVALUATE / PRINT */
        eval(cmdline);
    }
}

void eval(char *cmdline) {
    pid_t pid;

    char *cmds[MAX_CMDS];
    split(cmdline, cmds, "&&");

    for (int i = 0; i < MAX_CMDS; i++) {

    }

    char *argv[MAX_ARGS];
    split(cmdline, argv, " ");

    builtin_command(argv);

    if ((pid = Fork()) == 0) {
        // Executed by child
        // Exec replaces the program running in the child process with the program in argv
        if (execvp(argv[0], argv) < 0) {
            fprintf(stderr, "%s: Command not found. Errno: %d\n", argv[0], errno);
            _Exit(EXIT_FAILURE);
        }
    }

    // Executed by parent (since child process is always terminated or replaced by exec with a new program)
    int status;
    // Parent waits for the child to terminate
    if (waitpid(pid, &status, 0) < 0)
        unix_error("wait: waitpid error");
}

void builtin_command(char **argv) {
    // strcmp actually returns 0 (i.e. "false") when the strings are equal
    if (!strcmp(argv[0], "exit")) {
        exit_program();
    }
}

void exit_program() {
    printf("\nHave a good one!\n");
    exit(EXIT_SUCCESS);
}

/*
 * In this program, this SIGINT handler is only registered (see above) for the parent process.
 * So the program running in the child process will likely terminate (unless SIGINT is handled differently in the exec'd program),
 * while the parent process will simply print a newline, wait to receive find out that the child process terminated, and then loop.
 */
void sig_int_handler(int code) {
    UNUSED(code);
    write(STDOUT_FILENO, "\n", 1);
}

// Reference: CS:APP 8.3
pid_t Fork(void) {
    pid_t pid;
    if ((pid = fork()) < 0)
        unix_error("Fork error");
    return pid;
}

// Reference: CS:APP 8.3
void unix_error(char *msg) {
    fprintf(stderr, "%s: %s\n", msg, strerror(errno));
    exit(EXIT_FAILURE);
}

int split(char *input, char **argv, char *delim) {

    // Chop off trailing '\n'
    if (input[strlen(input)-1] == '\n')
        input[strlen(input)-1] = '\0';

    int i = 0;
    while (true) {
        // Point the ith element in the list to the input string
        argv[i++] = input;

        // Search input for the next delimiter character
        char* delim_p = strstr(input, delim);
        if (delim_p == NULL) {
            // Exit loop when we can't find another delimiter string in the remaining input
            break;
        }
        // Replace the start of the delimiter string with the null character
        *delim_p = '\0';
        // Increment input to right after the delimiter string
        input = delim_p + strlen(delim);
    }

    // Follow argv convention in setting argv[argc] = NULL
    argv[i] = NULL;

    return 0;
}
