package main

import (
	"github.com/djumanoff/amqp"
	"github.com/kirigaikabuto/common-lib31"
	products "github.com/kirigaikabuto/products31"
	"log"
)

func main() {
	productsMongoStore, err := products.NewProductStore(common.MongoConfig{
		Host:           "localhost",
		Port:           "27017",
		Database:       "ivi",
		CollectionName: "products",
	})
	if err != nil {
		log.Fatal(err)
	}
	productsAmqpEndpoints := products.NewProductAmqpEndpoints(productsMongoStore)
	rabbitConfig := amqp.Config{
		Host:     "localhost",
		Port:     5672,
		LogLevel: 5,
	}
	serverConfig := amqp.ServerConfig{
		ResponseX: "response",
		RequestX:  "request",
	}

	sess := amqp.NewSession(rabbitConfig)
	err = sess.Connect()
	if err != nil {
		panic(err)
		return
	}
	srv, err := sess.Server(serverConfig)
	if err != nil {
		panic(err)
		return
	}
	srv.Endpoint("products.create", productsAmqpEndpoints.CreateProductAmqpEndpoint())
	srv.Endpoint("products.list", productsAmqpEndpoints.ListProductAmqpEndpoint())
	err = srv.Start()
	if err != nil {
		panic(err)
		return
	}
}
