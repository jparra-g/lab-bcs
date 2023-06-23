Integrantes del Grupo:
- Pamela Latoja
- Rodrigo Henriquez
- Pablo Saavedra
- Javier Parra

Consideraciones: 

- Hay archivos build.sh en los proyectos custody-service y frontend-api para generar los docker
- En 01-kubernetes se agregan los yaml para generar los deployments y services
- Se crea un frontend-api para poder probar la correcta implementación
- Las validaciones se realizan solo en el backend, para mostrar la correcta implementación en GO
- Se corrige el DOCKERFILE con errores
- Si se usa la configuración actual los endpoints son los siguientes para hacer la prueba con Port Forwarding sobre el servicio de frontend

Endpoint Crear Custodia (POST): http://localhost:8080/custody/create

{
    "period":   "2023-02",
    "stock":    "stock",
    "client_id": "1-9",
    "quantity": 2
}


Endpoint Obtener Custodias (POST): http://localhost:8080/custody/getCustody

{
    "period":   "2023-02",
    "stock":    "stock",
    "client_id": "1-9"
}
