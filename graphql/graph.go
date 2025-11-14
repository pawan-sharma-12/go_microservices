package main

import (
"github.com/99designs/gqlgen/graphql"
"github.com/pawan-sharma-12/go_microservices/account"
"github.com/pawan-sharma-12/go_microservices/catalog"
"github.com/pawan-sharma-12/go_microservices/order"
)





type Server struct {
	accountClient *account.Client
	catalogClient *catalog.Client
	orderClient   *order.Client
}
func NewGraphQlServer(accountUrl, catalogUrl, orderUrl string) (*Server, error) {
    accountClient, err := account.NewClient(accountUrl)
    if err != nil {
        return nil, err
    }

    catalogClient, err := catalog.NewClient(catalogUrl)
    if err != nil {
        accountClient.Close() // close already created client
        return nil, err
    }

    orderClient, err := order.NewClient(orderUrl)
    if err != nil {
        accountClient.Close()
        catalogClient.Close()
        return nil, err
    }

    return &Server{
        accountClient: accountClient,
        catalogClient: catalogClient,
        orderClient:   orderClient,
    }, nil
}

func (s *Server) Mutation() MutationResolver {
	return &mutationResolver{
		server: s,
	}
}
func (s *Server) Query() QueryResolver {
	return &queryResolver{
		server: s,
	}
}
// Define the accountResolver type and implement the AccountResolver interface
// type accountResolver struct {
// 	Server *Server
// }

func (s *Server) Account() AccountResolver {
	return &accountResolver{
		server: s,
	}
}

func (s *Server) ToExecutableSchema() graphql.ExecutableSchema {
	return NewExecutableSchema(Config{
		Resolvers: s,
	})
}