package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	authproto "grpc-auth/proto"

	"google.golang.org/grpc"
)

type servidorAuth struct {
	authproto.UnimplementedAuthServiceServer
}

func (s *servidorAuth) ValidarToken(ctx context.Context, peticion *authproto.PeticionToken) (*authproto.RespuestaValidacion, error) {
	fmt.Printf("🔒 Petición recibida. Validando token: %s\n", peticion.Token)

	if peticion.Token == "token-secreto-123" {
		return &authproto.RespuestaValidacion{
			EsValido:  true,
			UsuarioId: "usr_999",
			Rol:       "administrador",
		}, nil
	}

	return &authproto.RespuestaValidacion{
		EsValido:  false,
		UsuarioId: "",
		Rol:       "",
	}, nil
}

func main() {
	puerto := ":50051"
	escuchador, err := net.Listen("tcp", puerto)
	if err != nil {
		fmt.Printf("❌ Error al abrir el puerto: %v\n", err)
		return
	}

	servidorGrpc := grpc.NewServer()
	authproto.RegisterAuthServiceServer(servidorGrpc, &servidorAuth{})

	fmt.Printf("🚀 Servidor gRPC de Autenticación escuchando en el puerto %s\n", puerto)

	go func() {
		if err := servidorGrpc.Serve(escuchador); err != nil {
			fmt.Printf("❌ Error en el servidor: %v\n", err)
		}
	}()

	parada := make(chan os.Signal, 1)
	signal.Notify(parada, os.Interrupt, syscall.SIGTERM)
	<-parada

	fmt.Println("🛑 Apagando el servidor gRPC...")
	servidorGrpc.GracefulStop()
}
