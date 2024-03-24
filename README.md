# Cordyceps
 
  Esse projeto, conta com um C2, um servidor de Comando em Controle. 

  Esse projeto esta em  constante Evolução. 

  <img align="center" height="40" width="50" src="https://raw.githubusercontent.com/devicons/devicon/master/icons/go/go-original.svg">


  Metodo de Uso:

 1 -  Ip Server: 
    No arquivo agente.go, linha 27, defina o IP publico do seu servidor onde ocorrerá a comunição 
   
2 - Build do Agente. 
   <h1>Windows:</h1>
    <p>GOOS=windows go build -ldflags -H=windowsgui &lt;arquivo.go&gt;</p>
    
    <h2>Linux:</h2>
    <p>GOOS=linux go build -ldflags -H=windowsgui &lt;arquivo.go&gt;</p>
    
    <h3>MacOS:</h3>
    <p>GOOS=darwin go build -ldflags -H=windowsgui &lt;arquivo.go&gt;</p>
