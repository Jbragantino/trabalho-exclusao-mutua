// Vinícius Martins e João Bragantino

package main

import (
	"fmt"
	"log"

	"github.com/graarh/golang-socketio/transport"

	gosocketio "github.com/graarh/golang-socketio"
)

func main() {
	socket, err := gosocketio.Dial(
		gosocketio.GetUrl("localhost", 8000, false),
		transport.GetDefaultWebsocketTransport(),
	)
	if err != nil {
		log.Fatal(err)
	}

	ehDonoRecurso := false
	aguardandoLiberacao := false

	socket.On("quadroDeAvisos", func(c *gosocketio.Channel, quadroDeAvisos []string) { OnQuadroDeAvisos(c, quadroDeAvisos) })
	socket.On("confirmarTrava", func(c *gosocketio.Channel) {
		OnConfirmarTrava(c, &aguardandoLiberacao, &ehDonoRecurso)
	})

	for {
		var opcao string
		fmt.Scanln(&opcao)

		executar(opcao, socket, &aguardandoLiberacao, &ehDonoRecurso)
	}
}

func executar(opcao string, socket *gosocketio.Client, aguardandoLiberacao *bool, ehDonoRecurso *bool) {
	switch opcao {
	case "1":
		solicitarTrava(socket, aguardandoLiberacao, ehDonoRecurso)
	case "2":
		liberarRecurso(socket, aguardandoLiberacao, ehDonoRecurso)
	case "3":
		lerQuadro(socket, ehDonoRecurso)
	case "4":
		adicionarAviso(socket, ehDonoRecurso)
	default:
		fmt.Println("Opção Invalida")
	}
}

func solicitarTrava(socket *gosocketio.Client, aguardandoLiberacao *bool, ehDonoRecurso *bool) {
	if *aguardandoLiberacao {
		fmt.Println("Recurso já foi solicitado")
	} else if *ehDonoRecurso {
		fmt.Println("Você já é dono do recurso")
	} else {
		fmt.Println("Solicitada trava do recurso")
		socket.Emit("travarRecurso", nil)
		*aguardandoLiberacao = true
	}
}

func liberarRecurso(socket *gosocketio.Client, aguardandoLiberacao *bool, ehDonoRecurso *bool) {
	if *ehDonoRecurso {
		fmt.Println("Recurso liberado")
		socket.Emit("liberarRecurso", nil)
		*ehDonoRecurso = false
	} else {
		fmt.Println("Você não é o dono do recurso")
	}
}

func lerQuadro(socket *gosocketio.Client, ehDonoRecurso *bool) {
	if *ehDonoRecurso {
		socket.Emit("lerQuadro", nil)
	} else {
		fmt.Println("Você não é o dono do recurso")
	}
}

func adicionarAviso(socket *gosocketio.Client, ehDonoRecurso *bool) {
	if *ehDonoRecurso {
		var aviso string
		fmt.Print("Digite o aviso: ")
		fmt.Scanln(&aviso)
		socket.Emit("adicionarAviso", aviso)
		fmt.Println("Aviso adicionado com sucesso")
	} else {
		fmt.Println("Você não é o dono do recurso")
	}
}

func OnQuadroDeAvisos(c *gosocketio.Channel, quadroDeAvisos []string) {
	fmt.Println("Quadro: ", quadroDeAvisos)
}

func OnConfirmarTrava(c *gosocketio.Channel, aguardandoLiberacao *bool, ehDonoRecurso *bool) {
	fmt.Println("Recurso travado")
	*aguardandoLiberacao = false
	*ehDonoRecurso = true
}
