# Postago

Postago es un servicio ligero de cola de correos electrónicos desarrollado en Go. Permite procesar y enviar correos electrónicos de manera asíncrona a través de una cola basada en Redis. Está diseñado para integrarse como una dependencia en otros proyectos que requieran una gestión eficiente del envío de correos electrónicos.

## Características

- **Procesamiento asíncrono:** Obtiene y envía correos desde una cola en Redis sin bloquear la ejecución del sistema.  
- **Integración sencilla:** Puede incluirse como una dependencia en otros proyectos.  
- **Despliegue con Docker:** Se puede ejecutar fácilmente como un servicio independiente mediante Docker y Docker Compose.  
- **Configuración flexible:** Se adapta a diferentes entornos mediante un archivo de configuración en YAML.  

## Requisitos previos

- **Redis:** Instalado localmente o en un contenedor Docker.  
- **Go:** Para desarrollo y ejecución directa.  
- **Docker y Docker Compose:** Opcional, para despliegue en contenedores.  

## Instalación

Para utilizar Postago en tu proyecto, clona el repositorio y construye la imagen de Docker:  

```sh
git clone https://github.com/your-username/postago.git
cd postago
docker build -t nombre/postago .
```

## Configuración

Postago obtiene su configuración desde un archivo YAML (`conf.yaml`). A continuación, un ejemplo de configuración:  

```yaml
redis: 
  address: redis
  port: 6379
  passwd: ""
  db: 0
mail:
  account:
    address: "tu-correo@example.com"
    passwd: "tu-contraseña"
  server:
    domain: smtp.gmail.com
    port: 587
queue: 
  name: welcome_email_queue
```

### Uso de plantillas de correo  
Postago utiliza `html/template` para procesar los correos electrónicos. Para cada mensaje en la cola, es necesario especificar un `templateName` en Redis junto con los datos del correo.  

Cada correo enviado debe contener la siguiente información en formato JSON:

```json
{
  "toEmail": "destinatario@example.com",
  "subject": "Asunto del correo",
  "templateName": "nombre_de_la_plantilla",
  "data": {
    "clave1": "valor1",
    "clave2": "valor2"
  }
}
```

Donde:
- **toEmail**: Dirección de correo del destinatario.
- **subject**: Asunto del correo.
- **templateName**: Nombre de la plantilla de correo utilizada.
- **data**: Objeto con los datos dinámicos que se usarán en la plantilla.

Postago extrae esta información de la cola en Redis y la procesa con la función `SendEmail`:

```go
emailService.SendEmail(
    infoMap["toEmail"].(string),
    infoMap["subject"].(string),
    infoMap["templateName"].(string),
    infoMap["data"].(map[string]interface{}),
)
```

## Ejecución del servicio

Para ejecutar Postago, primero construye la imagen de Docker:

```sh
docker build -t nombre/postago .
```

### Ejecutar con Docker Compose  
Para desplegar Postago junto con Redis, usa `docker-compose`:

```sh
docker-compose up -d
```

## Contribuir  

Las contribuciones son bienvenidas. Si encuentras algún problema o tienes una sugerencia, no dudes en abrir un *issue* o enviar un *pull request*.

