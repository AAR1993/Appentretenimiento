# appAlejo

Aplicativo web para la gestión de lugares y productos.

## Requisitos mínimos

- Go 1.16 o superior
- Navegador web moderno (Chrome 80+, Firefox 75+, Edge 80+, Safari 13+)
- (Opcional) Node.js para herramientas de desarrollo frontend

## Instalación

1. Clona el repositorio:
   ```bash
   git clone https://github.com/usuario/appAlejo.git
   cd appAlejo
   ```

2. Instala las dependencias del backend:
   ```bash
   cd backend
   go mod tidy
   ```

3. Ejecuta el backend:
   ```bash
   go run main.go
   ```

4. Abre el frontend:
   - Puedes abrir `frontend/index.html` directamente en tu navegador, o
   - Servir la carpeta `frontend` con un servidor estático (por ejemplo, usando `Live Server` en VSCode).

## Despliegue

- **Backend:** Puede desplegarse en cualquier servidor compatible con Go.
- **Frontend:** Puede alojarse en servicios de hosting estático.
- **Base de datos:** Actualmente se usa un archivo JSON, pero se recomienda migrar a una base de datos relacional para producción.

## Protocolos de transmisión

- HTTP/HTTPS para la comunicación web.
- API REST con intercambio de datos en formato JSON.

## Licencia

Este proyecto está licenciado bajo la [MIT License](LICENSE).

## Enlace del repositorio

[https://github.com/usuario/appAlejo](https://github.com/usuario/appAlejo)
