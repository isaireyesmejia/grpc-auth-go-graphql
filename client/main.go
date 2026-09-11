package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	// Importamos el contrato proto de nuestro proyecto
	authproto "grpc-auth/proto"
)

func main() {
	// 1. Establecer la conexión con el servidor en Docker (localhost:50051)
	// Usamos "WithTransportCredentials(insecure.NewCredentials())" porque es una prueba local sin certificados SSL
	fmt.Println("📞 Conectando al servidor gRPC en Docker...")
	conexion, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("❌ No se pudo conectar al servidor: %v", err)
	}
	// Nos aseguramos de cerrar la llamada al terminar el programa
	defer conexion.Close()

	// 2. Crear el "control remoto" (Cliente gRPC) utilizando la conexión abierta
	cliente := authproto.NewAuthServiceClient(conexion)

	// 3. Crear un contexto con tiempo límite (si el servidor tarda más de 5 segundos, la llamada se cancela)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// ==========================================
	// PRUEBA 1: Enviar el Token Correcto
	// ==========================================
	tokenCorrecto := "token-secreto-123"
	fmt.Printf("\n🔒 Enviando token correcto: '%s'\n", tokenCorrecto)

	// Invocamos la función remota directamente en el código
	respuesta1, err := cliente.ValidarToken(ctx, &authproto.PeticionToken{Token: tokenCorrecto})
	if err != nil {
		log.Fatalf("❌ Error al invocar la función: %v", err)
	}

	// Imprimimos el resultado que nos devolvió el servidor Go desde Docker
	fmt.Println("📥 Respuesta del servidor:")
	fmt.Printf("   - ¿Es válido?: %v\n", respuesta1.EsValido)
	fmt.Printf("   - Usuario ID: %s\n", respuesta1.UsuarioId)
	fmt.Printf("   - Rol asignado: %s\n", respuesta1.Rol)

	// ==========================================
	// PRUEBA 2: Enviar un Token Incorrecto
	// ==========================================
	tokenFalso := "token-hacker-456"
	fmt.Printf("\n🔓 Enviando token incorrecto: '%s'\n", tokenFalso)

	respuesta2, err := cliente.ValidarToken(ctx, &authproto.PeticionToken{Token: tokenFalso})
	if err != nil {
		log.Fatalf("❌ Error al invocar la función: %v", err)
	}

	fmt.Println("📥 Respuesta del servidor:")
	fmt.Printf("   - ¿Es válido?: %v\n", respuesta2.EsValido)
}
