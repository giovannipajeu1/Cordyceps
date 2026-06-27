#include "lolbin.h"
#include "shell.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

char *qnvtxbrm(const char *cmd) {
    char buf[8192];
    snprintf(buf, sizeof(buf),
        "osascript -e 'do shell script \"%s\"'", cmd);
    return dkqwvzmn(buf);
}

char *bvqtnxmr(const char *cmd) {
    char buf[8192];
    snprintf(buf, sizeof(buf),
        "/usr/bin/python3 -c \"import os; os.system('%s')\"", cmd);
    return dkqwvzmn(buf);
}
