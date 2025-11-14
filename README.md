# Go Microservices with GraphQL Gateway

A complete microservices architecture built with Go, featuring gRPC communication between services and a GraphQL gateway for unified API access.

## 🏗️ Architecture

```
GraphQL Gateway (Port 8000)
    ↓
┌─────────────────┬─────────────────┬─────────────────┐
│   Account       │    Catalog      │     Order       │
│   Service       │    Service      │    Service      │
│  (Port 50051)   │  (Port 50052)   │  (Port 50053)   │
│                 │                 │                 │
│  PostgreSQL     │  Elasticsearch  │  PostgreSQL     │
│  (Port 5432)    │  (Port 9200)    │  (Port 5432)    │
└─────────────────┴─────────────────┴─────────────────┘
```

## 🚀 Services

- **Account Service**: Manages user accounts (gRPC + PostgreSQL)
- **Catalog Service**: Product catalog with search (gRPC + Elasticsearch)
- **Order Service**: Order management with product details (gRPC + PostgreSQL)
- **GraphQL Gateway**: Unified API gateway aggregating all services

## 📋 Prerequisites

- Go 1.19+
- PostgreSQL
- Elasticsearch
- Protocol Buffers compiler (`protoc`)

## 🛠️ Setup

### 1. Install Dependencies

```bash
# Install protoc plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Install project dependencies
go mod tidy
```

### 2. Setup Databases

#### PostgreSQL Setup
```bash
# Create databases
createdb accountdb
createdb orderdb

# Run migrations
psql -d accountdb -f account/migrations/001_create_accounts_table.up.sql
psql -d orderdb -f order/migrations/1_create_orders_tables.up.sql
```

#### Elasticsearch Setup
```bash
# Start Elasticsearch (using Docker)
docker run -d --name elasticsearch -p 9200:9200 -e "discovery.type=single-node" elasticsearch:7.17.0
```

### 3. Environment Configuration

The project uses `.env.local` for local development:

```bash
# ACCOUNT SERVICE
ACCOUNT_DATABASE_URL=postgres://postgres:postgres@localhost:5432/accountdb?sslmode=disable
ACCOUNT_SERVICE_URL=localhost:50051

# CATALOG SERVICE
CATALOG_DATABASE_URL=postgres://postgres:postgres@localhost:5432/catalogdb?sslmode=disable
CATALOG_SERVICE_URL=localhost:50052
CATALOG_ELASTIC_URL=http://localhost:9200

# ORDER SERVICE
ORDER_DATABASE_URL=postgres://postgres:postgres@localhost:5432/orderdb?sslmode=disable
ORDER_SERVICE_URL=localhost:50053
```

## 🏃‍♂️ Running the Services

Start each service in separate terminals:

### 1. Account Service
```bash
go run account/cmd/account/main.go
```

### 2. Catalog Service
```bash
go run catalog/cmd/catalog/main.go
```

### 3. Order Service
```bash
go run order/cmd/order/main.go
```

### 4. GraphQL Gateway
```bash
cd graphql
go run .
```

## 🧪 Testing

Access GraphQL Playground at: `http://localhost:8000/playground`

### Sample Test Flow

1. **Create Account**
2. **Create Products**
3. **Create Orders**
4. **Query Data with Relationships**

See `test_data.json` for complete test queries and mutations.

## 🐳 Docker Support

### Development with Docker Compose
```bash
# Start all services
docker-compose up

# Stop all services
docker-compose down
```

### Manual Docker Commands
```bash
# Build individual service
docker build -f account/app.dockerfile -t account-service .

# Run with environment
docker run -p 50051:8080 --env-file .env.local account-service
```

## 🔧 Development Commands

### Regenerate Protocol Buffers
```bash
# Account service
protoc --go_out=./account/pb --go-grpc_out=./account/pb account/account.proto

# Catalog service  
protoc --go_out=./catalog/pb --go-grpc_out=./catalog/pb catalog/catalog.proto

# Order service
protoc --go_out=./order/pb --go-grpc_out=./order/pb order/order.proto
```

### Database Operations
```bash
# Connect to PostgreSQL
psql -h localhost -U postgres -d accountdb

# Connect to Docker PostgreSQL
docker exec -it account_db psql -U postgres -d account_db

# Check Elasticsearch
curl http://localhost:9200/_cluster/health
```

## 📊 API Endpoints

### GraphQL Mutations
- `createAccount(account: AccountInput!): Account!`
- `createProduct(product: ProductInput!): Product!`
- `createOrder(order: OrderInput!): Order!`

### GraphQL Queries
- `accounts(pagination: PaginationInput!, id: String): [Account!]!`
- `products(pagination: PaginationInput!, query: String, id: String): [Product!]!`

### Nested Resolvers
- `Account.orders: [Order!]!` - Get all orders for an account

## 🔍 Troubleshooting

### Common Issues

1. **Port Conflicts**
   - Ensure no other services are using ports 8000, 50051, 50052, 50053
   - Check with: `lsof -i :8000`

2. **Database Connection Issues**
   - Verify PostgreSQL is running: `pg_isready`
   - Check database exists: `psql -l`

3. **gRPC Connection Timeouts**
   - Ensure all services are running
   - Check service logs for connection errors

4. **Elasticsearch Issues**
   - Verify Elasticsearch is running: `curl http://localhost:9200`
   - Check cluster health: `curl http://localhost:9200/_cluster/health`

### Logs
Each service provides detailed logging:
- ✅ Success operations
- ❌ Error conditions  
- 🔗 Service connections
- 📦 Configuration loading

## 🏗️ Project Structure

```
go_microservices/
├── account/           # Account microservice
│   ├── cmd/account/   # Service entry point
│   ├── pb/           # Generated protobuf files
│   └── migrations/   # Database migrations
├── catalog/          # Catalog microservice
├── order/            # Order microservice
├── graphql/          # GraphQL gateway
├── config/           # Configuration files
├── .env.local        # Local environment variables
├── docker-compose.yaml
└── README.md
```

## 🤝 Contributing

1. Fork the repository
2. Create feature branch
3. Make changes
4. Test thoroughly
5. Submit pull request

## 📝 License

This project is licensed under the MIT License.