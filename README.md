# Clean Architecture Challenge

Este projeto implementa um sistema de pedidos seguindo os princípios da Clean Architecture.

## Tecnologias Utilizadas

- Go 1.23.5
- MySQL
- Docker
- gRPC
- GraphQL
- REST API

## Serviços e Portas

- REST API: `http://localhost:8000`
- gRPC: `localhost:50051`
- GraphQL: `http://localhost:8080/query`
- MySQL: `localhost:3306`

## Como Executar

1. Clone o repositório:
```bash
git clone https://github.com/JMKobayashi/GoExpert-ClearArchitecture.git
cd GoExpert-ClearArchitecture
```

2. Inicie os containers:
```bash
docker compose up -d
```

3. Execute as migrações do banco de dados:
```bash
go run cmd/migrate/main.go
```

4. Inicie a aplicação:
```bash
go run cmd/api/main.go
```

## Endpoints

### REST API

- Criar Pedido:
```http
POST http://localhost:8000/order
Content-Type: application/json

{
    "id": "a",
    "price": 100.5,
    "tax": 0.5
}
```

- Listar Pedidos:
```http
GET http://localhost:8000/order
```

### GraphQL

- Query para listar pedidos:
```graphql
query {
  orders {
    id
    price
    tax
    finalPrice
  }
}
```

### gRPC

- Serviço de listagem de pedidos disponível na porta 50051 