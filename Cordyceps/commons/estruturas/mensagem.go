package commons

type Mensagem struct {
	AgentID   string
	AgentHost string
	AgentCWD  string
	Comandos  []Commando
}
