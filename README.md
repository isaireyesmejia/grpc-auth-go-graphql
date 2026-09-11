# Microservicio de Autenticación Híbrido (GraphQL + gRPC + Go)

Este proyecto implementa el patrón arquitectónico **BFF (Backend for Frontend)** utilizando las tecnologías más eficientes del mercado actual para demostrar una comunicación desacoplada y de alta velocidad.

## 🏗️ Arquitectura del Sistema

1. **Frontend / Cliente Público:** Expuesto mediante un Gateway de **GraphQL (Apollo Server)** en Node.js (Puerto 4000).
2. **Comunicación Interna:** Traducción de consultas de texto a formato binario ultra ligero usando **gRPC** bajo HTTP/2.
3. **Microservicio Backend:** Core de autenticación e infraestructura escrito en **Go (Golang)** corriendo de forma aislada dentro de un contenedor de **Docker** (Puerto 50051).

## 🚀 Cómo ejecutar el proyecto

### Requisitos previos
* Docker Desktop instalado.
* Node.js instalado localmente.

### Pasos para encender
1. Clonar el repositorio.
2. Encender el servidor backend de Go en Docker:
   ```bash
   docker compose up --build -d
   ```
3. Instalar dependencias del cliente e iniciar el servidor GraphQL:
   ```bash
   npm install
   node ./client/graphql_server.js
   ```
4. Abrir en el navegador `http://localhost:4000` para interactuar con la API.
