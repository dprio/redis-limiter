# credis-limiter

Este projeto é uma aplicação que possibilita uma simples criação de order e consulta das orders criadas através de diferentes servidores (web, graphql e grpc)

## Pré-requisitos
- Docker e Docker Compose
- Go 1.21 ou superior
- evans
- Make

## Como subir a aplicação

`docker-compose up`

## Funcionamento

### Criação de order
#### web 
```
curl --location 'http://localhost:8080/orders' \
--header 'Content-Type: application/json' \
--data '{
    "price": 101.01,
    "tax": 11.5
}'
```

#### GraphQL
Acessar a url http://localhost:8081/ e rora a mutation:

```
mutation createOrder {
  createOrder(
    input: {
      price: 1.13, 
      tax: 10.51
    }
  ){
    id
    price
    tax
    final_price
  }
}
```

#### gRPC
No terminal inciar o o evans 
```
evans -r repl
```

Selecionar o serviço de criação de order

### Consulta de orders
#### web 
```
curl --location 'http://localhost:8080/orders'
```

#### GraphQL
Acessar a url http://localhost:8081/ e rodar a query:

```
query getOrders {
  getOrders{
      id
      price
      tax
      final_price
  }
}
```

#### gRPC
No terminal inciar o evans 
```
evans -r repl
```

Selecionar o serviço de busca de orders
