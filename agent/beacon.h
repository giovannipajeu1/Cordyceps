#ifndef BEACON_H
#define BEACON_H
typedef struct {
    char cmd[1024];
    char arg[2048];
} BXQVTNMR;
int  nqvxbrtm(const char *host, int port,
                    const char *agent_id, BXQVTNMR *out);
void kxvqtnbr(const char *host, int port,
                    const char *agent_id, const char *response);
#endif
