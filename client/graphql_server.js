const { ApolloServer } = require('@apollo/server');
const { startStandaloneServer } = require('@apollo/server/standalone');
const grpc = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');
const path = require('path');

// =========================================================================
// 1. CONFIGURACIÓN DE CLIENTE gRPC (Comunicación interna con Docker)
// =========================================================================
const PROTO_PATH = path.join(__dirname, '../proto/auth.proto');
const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
    keepCase: true, longs: String, enums: String, defaults: true, oneofs: true
});
const authProto = grpc.loadPackageDefinition(packageDefinition);
// Conectamos al Walkie-Talkie del contenedor de Go
const clienteGrpc = new authProto.AuthService('localhost:50051', grpc.credentials.createInsecure());

// Prometificamos la llamada gRPC para poder usar async/await de forma limpia en GraphQL
const validarTokenPromesa = (token) => {
    return new Promise((resolve, reject) => {
        clienteGrpc.ValidarToken({ token: token }, (error, respuesta) => {
            if (error) reject(error);
            else resolve(respuesta);
        });
    });
};

// =========================================================================
// 2. CONFIGURACIÓN DE GRAPHQL (La cara hacia el usuario/frontend)
// =========================================================================

// Definimos el esquema de GraphQL: Lo que el cliente público puede consultar
const typeDefs = `#graphql
  type PerfilUsuario {
    es_valido: Boolean!
    usuario_id: String
    rol: String
  }

  type Query {
    # El usuario enviará un token y recibirá la estructura del perfil
    obtenerPerfil(token: String!): PerfilUsuario
  }
`;

// Los Resolvers: La lógica que decide de dónde sacar los datos consultados
const resolvers = {
    Query: {
        obtenerPerfil: async (_, args) => {
            console.log(`\n📬 GraphQL recibió una petición del usuario. Token enviado: "${args.token}"`);
            
            try {
                // LLAMADA INTERNA VIA gRPC: Le pedimos al servidor de Go en Docker que valide el token
                console.log("⚡ Redireccionando validación internamente a Go por gRPC binario...");
                const resultadoBinario = await validarTokenPromesa(args.token);
                
                // Retornamos la respuesta binaria formateada de vuelta a GraphQL
                return {
                    es_valido: resultadoBinario.es_valido,
                    usuario_id: resultadoBinario.usuario_id,
                    rol: resultadoBinario.rol
                };
            } catch (error) {
                console.error("❌ Error en la comunicación gRPC interna:", error);
                throw new Error("Error interno del servidor al procesar la autenticación");
            }
        }
    }
};

// =========================================================================
// 3. ARRANCAR EL SERVIDOR DE APOLLO GRAPHQL
// =========================================================================
const server = new ApolloServer({ typeDefs, resolvers });

async function encenderServidor() {
    const { url } = await startStandaloneServer(server, { listen: { port: 4000 } });
    console.log(`🚀 Servidor GraphQL (BFF) listo en la web pública: ${url}`);
    console.log(`🔒 gRPC conectado internamente en el puerto 50051 hacia Docker.`);
}

encenderServidor();
