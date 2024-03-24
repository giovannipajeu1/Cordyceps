# Comando e Controle Golang - Cordyceps

Um Servidor de Comando e Controle feito em Golang.

## 🚀 Começando

Essas instruções permitirão que você obtenha uma cópia do projeto em operação na sua máquina local para fins de desenvolvimento e teste.

### 📋 Pré-requisitos

Você precisara da Linguagem Golang instalada na sua maquina


Como instalar o Golang? 
<br>
<a href="https://go.dev/doc/install">Passo a Passo para instalar o Golang</a>


### 🔧 Instalação

Apos o Golang Instalado, precisamos buildar o agente e o servidor
###

Agente para Windows   <img align="center" height="20" width="20" src="https://raw.githubusercontent.com/devicons/devicon/master/icons/windows11/windows11-original.svg"> :
```
GOOS=windows go build  -ldflags -H=windowsgui agente.go
```
Agente para Linux <img align="center" height="20" width="20" src="https://raw.githubusercontent.com/devicons/devicon/master/icons/linux/linux-original.svg"> :
```
GOOS=linux go build agente.go
```
Agente para MacOs <img align="center" height="20" width="20" src="https://raw.githubusercontent.com/devicons/devicon/master/icons/apple/apple-original.svg">  : 
```
GOOS=dawrin go build agente.go
```

## 📦 Implantação

Adicione notas adicionais sobre como implantar isso em um sistema ativo

## 🛠️ Construído com

*  <img align="center" height="20" width="20" src="https://raw.githubusercontent.com/devicons/devicon/master/icons/go/go-original.svg"> - Linguagem Usada