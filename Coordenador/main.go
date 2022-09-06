// Vinícius Martins e João Bragantino

package main

import (
	"fmt"
	"log"
	"net/http"

	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

func main() {
	var donoRecurso string
	filaRecurso := make([]string, 0)

	quadroDeAvisos := make([]string, 0)

	server := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())

	server.On(gosocketio.OnConnection, func(c *gosocketio.Channel) { OnClienteConectado(c) })

	server.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) { OnClienteDesconectado(c, &donoRecurso, &filaRecurso) })

	server.On("travarRecurso", func(c *gosocketio.Channel) { OnTravarRecurso(c, &donoRecurso, &filaRecurso) })

	server.On("liberarRecurso", func(c *gosocketio.Channel) { OnLiberarRecurso(c, &donoRecurso, &filaRecurso) })

	server.On("lerQuadro", func(c *gosocketio.Channel) { OnLerQuadro(c, donoRecurso, quadroDeAvisos) })

	server.On("adicionarAviso", func(c *gosocketio.Channel, aviso string) { OnAdicionarAviso(c, donoRecurso, &quadroDeAvisos, aviso) })

	serveMux := http.NewServeMux()
	serveMux.Handle("/socket.io/", server)
	log.Println("Iniciando servidor...")
	log.Println("Servidor iniciado")
	log.Panic(http.ListenAndServe(":8000", serveMux))
}

func OnClienteConectado(c *gosocketio.Channel) {
	log.Println("Cliente conectado: ", c.Id())
	c.Join(c.Id())
}

func OnClienteDesconectado(c *gosocketio.Channel, donoRecurso *string, filaRecurso *[]string) {
	log.Println("Cliente desconectado: ", c.Id())
	c.Leave(c.Id())

	*filaRecurso = remover(*filaRecurso, c.Id())

	OnLiberarRecurso(c, donoRecurso, filaRecurso)
}

func remover(lista []string, valor string) []string {
	for i, other := range lista {
		if other == valor {
			return append(lista[:i], lista[i+1:]...)
		}
	}
	return lista
}

func OnTravarRecurso(c *gosocketio.Channel, donoRecurso *string, filaRecurso *[]string) {
	if *donoRecurso == "" {
		fmt.Println("Recurso travado para: ", c.Id())
		*donoRecurso = c.Id()
		c.Emit("confirmarTrava", nil)
	} else {
		*filaRecurso = append(*filaRecurso, c.Id())
		fmt.Println(c.Id(), " adicionado na fila")
		fmt.Println("Fila de espera: ", *filaRecurso)
	}
}

func OnLiberarRecurso(c *gosocketio.Channel, donoRecurso *string, filaRecurso *[]string) {
	if *donoRecurso != c.Id() {
		return
	}

	fmt.Println("Recurso liberado por: ", c.Id())
	if len(*filaRecurso) > 0 {
		*donoRecurso = (*filaRecurso)[0]
		*filaRecurso = (*filaRecurso)[1:]
		c.BroadcastTo(*donoRecurso, "confirmarTrava", nil)
		fmt.Println("Recurso travado para: ", *donoRecurso)
	} else {
		*donoRecurso = ""
	}
}

func OnLerQuadro(c *gosocketio.Channel, donoRecurso string, quadroDeAvisos []string) {
	if donoRecurso != c.Id() {
		return
	}

	log.Println("Ler quadro: ", c.Id())
	c.Emit("quadroDeAvisos", quadroDeAvisos)
}

func OnAdicionarAviso(c *gosocketio.Channel, donoRecurso string, quadroDeAvisos *[]string, aviso string) {
	if donoRecurso != c.Id() {
		return
	}

	log.Println("Adicionar aviso: ", c.Id())
	*quadroDeAvisos = append(*quadroDeAvisos, aviso)
}
