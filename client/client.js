const grpc = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');
const path = require('path');

const PROTO_PATH = path.join(__dirname, '../proto/auth.proto');
const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
    keepCase: true, longs: String, enums: String, defaults: true, oneofs: true
});
const authProto = grpc.loadPackageDefinition(packageDefinition);
const cliente = new authProto.AuthService('localhost:50051', grpc.credentials.createInsecure());

function escucharAlertas() {
    console.log("📞 Solicitando canal de Streaming en vivo a Docker...");

    // Invocamos la función de streaming pasando nuestro ID
    const canalStream = cliente.TransmitirAlertas({ cliente_id: "consola_monitoreo_node" });

    // Escuchamos el evento 'data'. Cada vez que Go envíe una alerta por el cable, este evento se dispara solo
    canalStream.on('data', (alerta) => {
        console.log(`\n🚨 [${alerta.timestamp}] Nueva Alerta Recibida:`);
        console.log(`   - Nivel:   ${alerta.nivel}`);
        console.log(`   - Mensaje: ${alerta.mensaje}`);
    });

    // Detectar cuando el servidor decide cerrar el canal de streaming
    canalStream.on('end', () => {
        console.log("\n🏁 El servidor de Go ha cerrado el canal de streaming con éxito.");
    });

    canalStream.on('error', (err) => {
        console.error("❌ Error en el stream:", err);
    });
}

escucharAlertas();
