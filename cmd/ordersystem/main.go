package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/JMKobayashi/GoExpert-ClearArchitecture/configs"
	"github.com/JMKobayashi/GoExpert-ClearArchitecture/internal/event/handler"
	"github.com/JMKobayashi/GoExpert-ClearArchitecture/internal/infra/graph"
	"github.com/JMKobayashi/GoExpert-ClearArchitecture/internal/infra/grpc/pb"
	"github.com/JMKobayashi/GoExpert-ClearArchitecture/internal/infra/grpc/service"
	"github.com/JMKobayashi/GoExpert-ClearArchitecture/internal/infra/web/webserver"
	"github.com/JMKobayashi/GoExpert-ClearArchitecture/pkg/events"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Define the expected config file path inside the container
	configFilePath := "/app/.env" // Path where the .env file is mounted

	configs, err := configs.LoadConfig(configFilePath) // Pass the full path
	if err != nil {
		panic(fmt.Sprintf("Failed to load configuration from %s: %v", configFilePath, err)) // Improved panic message
	}

	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rabbitMQChannel := getRabbitMQChannel()

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	listOrdersUseCase := NewListOrdersUseCase(db)

	webserver := webserver.NewWebServer(":8000")
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)
	listOrdersHandler := NewListOrdersHandler(db)

	webserver.AddHandler("/order", webOrderHandler.Create, "POST")
	webserver.AddHandler("/order", listOrdersHandler.Handle, "GET")

	fmt.Println("Starting web server on port", configs.WebServerPort)
	go webserver.Start()

	grpcServer := grpc.NewServer()
	createOrderService := service.NewOrderService(*createOrderUseCase)
	pb.RegisterOrderServiceServer(grpcServer, createOrderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", configs.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
		ListOrdersUseCase:  *listOrdersUseCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", configs.GraphQLServerPort)
	http.ListenAndServe(":"+configs.GraphQLServerPort, nil)
}

func getRabbitMQChannel() *amqp.Channel {
	rabbitMQHost := os.Getenv("RABBITMQ_HOST")
	if rabbitMQHost == "" {
		rabbitMQHost = "localhost"
	}
	conn, err := amqp.Dial(fmt.Sprintf("amqp://guest:guest@%s:5672/", rabbitMQHost))
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}
