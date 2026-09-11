package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	authproto "grpc-auth/proto"

	"google.golang.org/grpc"
)

type servidorAuth struct {
	authproto.UnimplementedAuthServiceServer
}

func (s *servidorAuth) ValidarToken(ctx context.Context, peticion *authproto.PeticionToken) (*authproto.RespuestaValidacion, error) {
	fmt.Printf("🔒 Petición recibida. Validando token: %s\n", peticion.Token)
	if peticion.Token == "token-secreto-123" {
		return &authproto.RespuestaValidacion{EsValido: true, UsuarioId: "usr_999", Rol: "administrador"}, nil
	}
	return &authproto.RespuestaValidacion{EsValido: false, UsuarioId: "", Rol: ""}, nil
}

// NUEVO: Función que mantiene el canal abierto y envía datos constantes
func (s *servidorAuth) TransmitirAlertas(peticion *authproto.PeticionMonitoreo, stream authproto.AuthService_TransmitirAlertasServer) error {
	fmt.Printf("📡 Cliente '%s' conectado al canal de Streaming en tiempo real.\n", peticion.ClienteId)

	alertasSimuladas := []struct {
		nivel   string
		mensaje string
	}{
		{"INFO", "Intento de login desde IP 192.168.1.15"},
		{"WARNING", "Múltiples peticiones fallidas del usuario usr_404"},
		{"CRITICAL", "Intento de inyección SQL detectado en Gateway"},
		{"INFO", "Token de sesión renovado correctamente"},
		{"SUCCESS", "Auditoría de seguridad semanal completada"},
	}

	for i, alerta := range alertasSimuladas {
		// Construir el mensaje de la alerta
		msg := &authproto.AlertaSeguridad{
			Timestamp: time.Now().Format("15:04:05"),
			Nivel:     alerta.nivel,
			Mensaje:   alerta.mensaje,
		}

		// Enviar el mensaje por el stream abierto
		if err := stream.Send(msg); err != nil {
			fmt.Printf("❌ Error al enviar streaming al cliente: %v\n", err)
			return err
		}
		fmt.Printf("  [Stream] Enviada alerta %d/5: %s\n", i+1, alerta.nivel)

		// Esperar 2 segundos antes de enviar la siguiente alerta
		time.Sleep(2 * time.Second)
	}

	fmt.Println("🏁 Transmisión de streaming finalizada con éxito.")
	return nil
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

	fmt.Printf("🚀 Servidor gRPC con Streaming escuchando en el puerto %s\n", puerto)

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
