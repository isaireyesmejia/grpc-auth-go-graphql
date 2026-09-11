const grpc = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');
const path = require('path');

// 1. Cargar el contrato .proto directamente desde su carpeta
const PROTO_PATH = path.join(__dirname, '../proto/auth.proto');
const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
    keepCase: true,
    longs: String,
    enums: String,
    defaults: true,
    oneofs: true
});

const authProto = grpc.loadPackageDefinition(packageDefinition);

function main() {
    // 2. Conectar al Walkie-Talkie del servidor Go en Docker
    console.log("📞 Conectando al servidor gRPC en Docker (localhost:50051)...");
    const cliente = new authProto.AuthService('localhost:50051', grpc.credentials.createInsecure());

    // ==========================================
    // PRUEBA 1: Enviar Token Correcto
    // ==========================================
    console.log("\n🔒 Enviando token correcto...");
    cliente.ValidarToken({ token: "token-secreto-123" }, (error, respuesta) => {
        if (error) {
            console.error("❌ Error:", error);
            return;
        }
        console.log("📥 Respuesta del servidor Go:");
        console.log(`   - ¿Es válido?: ${respuesta.es_valido}`);
        console.log(`   - Usuario ID: ${respuesta.usuario_id}`);
        console.log(`   - Rol asignado: ${respuesta.rol}`);

        // ==========================================
        // PRUEBA 2: Enviar Token Incorrecto
        // ==========================================
        console.log("\n🔓 Enviando token incorrecto...");
        cliente.ValidarToken({ token: "token-falso-hacker" }, (error, respuestaFalsa) => {
            if (error) return;
            console.log("📥 Respuesta del servidor Go:");
            console.log(`   - ¿Es válido?: ${respuestaFalsa.es_valido}`);
        });
    });
}

main();
