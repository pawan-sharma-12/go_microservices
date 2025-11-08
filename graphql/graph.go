package main

import "github.com/99designs/gqlgen/graphql"




type Server struct {
	// accountClient *accont.Client
	// catalogClient *catalog.Client
	// orderClient   *order.Client
}
func NewGraphQlServer(accountUrl, catalogUrl, orderUrl string) (*Server, error) {
	// Implementation goes here
	// accountClient, err := account.NewClient(accountUrl)
	// if err != nil {
	// 	accountClient.Close()
	// 	return nil, err
	// }
	// catalogClient, err := catalog.NewClient(catalogUrl)
	// if err != nil {
	// 	catalogClient.Close()
	// 	return nil, err
	// }
	// orderClient, err := order.NewClient(orderUrl)
	// if err != nil {
	// 	orderClient.Close()
	// 	return nil, err
	// }
	return &Server{
		// accountClient: accountClient,
		// catalogClient: catalogClient,
		// orderClient:   orderClient,
	}, nil

}
// func (s *Server) Mutataion() MutationResolver {
// 	return &mutationResolver{
// 		Server: s,
// 	}
// }
// func (s *Server) Query() QueryResolver {
// 	return &queryResolver{
// 		Server: s,
// 	}
// }
// Define the accountResolver type and implement the AccountResolver interface
type accountResolver struct {
	Server *Server
}

func (s *Server) Account() AccountResolver {
	return &accountResolver{
		Server: s,
	}
}

func (s *Server) ToExecutableSchema() graphql.ExecutableSchema {
	return NewExecutableSchema(Config{
		Resolvers: s,
	})
}