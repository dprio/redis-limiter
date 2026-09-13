# redis-limiter

Este projeto é uma exemplificação da utilização do redis como um middleware de rate-limit. 

O middleware está aplicado em dois endpoints, de listagem de orders e criação de order. O controle de rate-limit é feito por IP ou por um token passado como header (`API_KEY`)

## Pré-requisitos
- Docker e Docker Compose
- Go 1.21 ou superior
- evans
- Make

## Como subir a aplicação

`docker-compose up`

## Funcionamento

### Criação de order
```
curl --location 'http://localhost:8080/orders' \
--header 'Content-Type: application/json' \
--data '{
    "price": 101.01,
    "tax": 11.5
}'
```

### Consulta de orders
```
curl --location 'http://localhost:8080/orders'
```

### Utilizando o token
#### Primeiro é necessário criar o token configurando o número de requests possíveis no período avaliado.

```
curl --location 'http://localhost:8080/tokens' \
--header 'Content-Type: application/json' \
--data '{
    "total_requests": 5
}'
```

Nas resposta, obtemos um token que pode ser passado no header nas requests de order

>[!WARNING]
> O token não foi pensado para ser uma implementação real produtiva. A implementação é apenas para exemplificar a associação de uma configuração de ratelimit a um entidade de identificação.

```
curl --location 'http://localhost:8080/orders' \
--header 'Content-Type: application/json' \
--header 'API_KEY': {token}'\
--data '{
    "price": 101.01,
    "tax": 11.5
}'
```

```
curl --location 'http://localhost:8080/orders' \
--header 'API_KEY': {token}'
```
 <br>

> [!NOTE]
> **É possível executar as requests pelo arquivo de endpoint.http**
