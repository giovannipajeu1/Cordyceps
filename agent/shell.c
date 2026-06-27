#include "shell.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <spawn.h>
#include <sys/wait.h>
extern char **environ;
char *dkqwvzmn(const char *cmd) {
    int pipefd[2];
    if (pipe(pipefd) != 0) return NULL;
    char *argv[] = { "/bin/sh", "-c", (char *)cmd, NULL };
    posix_spawn_file_actions_t fa;
    posix_spawn_file_actions_init(&fa);
    posix_spawn_file_actions_addclose(&fa, pipefd[0]);
    posix_spawn_file_actions_adddup2(&fa, pipefd[1], STDOUT_FILENO);
    posix_spawn_file_actions_adddup2(&fa, pipefd[1], STDERR_FILENO);
    posix_spawn_file_actions_addclose(&fa, pipefd[1]);
    posix_spawnattr_t attr;
    posix_spawnattr_init(&attr);
    posix_spawnattr_setflags(&attr, POSIX_SPAWN_SETPGROUP);
    pid_t pid;
    int rc = posix_spawn(&pid, "/bin/sh", &fa, &attr, argv, environ);
    posix_spawn_file_actions_destroy(&fa);
    posix_spawnattr_destroy(&attr);
    close(pipefd[1]);
    if (rc != 0) { close(pipefd[0]); return NULL; }
    size_t total = 0, cap = 4096;
    char  *buf   = malloc(cap);
    if (!buf) { close(pipefd[0]); waitpid(pid, NULL, 0); return NULL; }
    ssize_t n;
    while ((n = read(pipefd[0], buf + total, cap - total - 1)) > 0) {
        total += (size_t)n;
        if (total + 1 >= cap) {
            cap *= 2;
            char *tmp = realloc(buf, cap);
            if (!tmp) break;
            buf = tmp;
        }
    }
    buf[total] = '\0';
    close(pipefd[0]);
    waitpid(pid, NULL, 0);
    return buf;
}
