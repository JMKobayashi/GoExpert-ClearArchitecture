# Clean Architecture Challenge

Este projeto implementa um sistema de pedidos seguindo os princípios da Clean Architecture.

## Tecnologias Utilizadas

- Go 1.23.5
- MySQL
- Docker
- gRPC
- GraphQL
- REST API

## Pré-requisitos

- Go 1.23.5 ou superior
- Docker e Docker Compose
- Wire (para injeção de dependência)
- Protocol Buffers (para gRPC)
- MySQL Client (opcional, para debug)

## Estrutura do Projeto

```
.
├── cmd/
│   ├── ordersystem/    # Aplicação principal
│   └── migrate/        # Migrações do banco de dados
├── internal/
│   ├── domain/         # Entidades e regras de negócio
│   ├── usecase/        # Casos de uso
│   ├── infra/          # Implementações concretas
│   │   ├── database/   # Repositórios
│   │   ├── grpc/       # Serviços gRPC
│   │   ├── graph/      # Resolvers GraphQL
│   │   └── web/        # Handlers HTTP
│   └── event/          # Eventos do sistema
└── pkg/                # Pacotes compartilhados
```

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

2. Instale as dependências:
```bash
go mod download
```

3. Inicie os containers:
```bash
docker compose up -d
```

4. Execute as migrações do banco de dados:
```bash
go run cmd/migrate/main.go
```

5. Gere o arquivo de injeção de dependência:
```bash
cd cmd/ordersystem
wire
```

6. Inicie a aplicação:
```bash
go run main.go wire_gen.go
```

## Endpoints

### REST API

#### Criar Pedido
```http
POST http://localhost:8000/order
Content-Type: application/json

{
    "id": "order-1",
    "price": 100.5,
    "tax": 0.5
}
```

Resposta:
```json
{
    "id": "order-1",
    "price": 100.5,
    "tax": 0.5,
    "final_price": 101.0
}
```

#### Listar Pedidos
```http
GET http://localhost:8000/order
```

Resposta:
```json
[
    {
        "id": "order-1",
        "price": 100.5,
        "tax": 0.5,
        "final_price": 101.0
    }
]
```

### GraphQL

#### Criar Pedido
```graphql
mutation {
  createOrder(input: {
    id: "order-1"
    price: 100.5
    tax: 0.5
  }) {
    id
    price
    tax
    finalPrice
  }
}
```

#### Listar Pedidos
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

#### Criar Pedido
```protobuf
rpc CreateOrder(CreateOrderRequest) returns (CreateOrderResponse)

message CreateOrderRequest {
  string id = 1;
  float price = 2;
  float tax = 3;
}

message CreateOrderResponse {
  string id = 1;
  float price = 2;
  float tax = 3;
  float final_price = 4;
}
```

#### Listar Pedidos
```protobuf
rpc ListOrders(ListOrdersRequest) returns (ListOrdersResponse)

message ListOrdersRequest {}

message ListOrdersResponse {
  repeated CreateOrderResponse orders = 1;
}
```

## Testando os Endpoints

### REST API
```bash
# Criar um pedido
curl -X POST http://localhost:8000/order \
  -H "Content-Type: application/json" \
  -d '{"id":"order-1","price":100.5,"tax":0.5}'

# Listar pedidos
curl http://localhost:8000/order
```

### GraphQL
Acesse o playground GraphQL em `http://localhost:8080` para testar as queries e mutations.

### gRPC
Use o grpcurl para testar os endpoints gRPC:

```bash
# Instalar grpcurl
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# Listar serviços disponíveis
grpcurl -plaintext localhost:50051 list

# Criar um pedido
grpcurl -plaintext -d '{
  "id": "order-1",
  "price": 100.5,
  "tax": 0.5
}' localhost:50051 pb.OrderService/CreateOrder

# Listar pedidos
grpcurl -plaintext localhost:50051 pb.OrderService/ListOrders
```

## Validações

O sistema realiza as seguintes validações:

- ID do pedido não pode ser vazio
- Preço deve ser maior que zero
- Taxa deve ser maior que zero
- Preço final é calculado automaticamente (preço + taxa)

## Eventos

O sistema emite eventos quando um pedido é criado:

- Evento: `OrderCreated`
- Payload: Dados do pedido criado
- Consumido por: RabbitMQ
