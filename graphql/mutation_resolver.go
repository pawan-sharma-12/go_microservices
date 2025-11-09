//Create Account

//  Create Product

//  Create Order

package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/pawan-sharma-12/go_microservices/order"
)

var (
	ErrInvalidParameter  = errors.New("invalid parameter")
)
type mutationResolver struct{ 
	server *Server
}
func (r *mutationResolver) createAccount(ctx context.Context, in AccountInput) (*Account, error) {

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	a, err := r.server.accountClient.PostAccount(ctx, in.Name, in.Email)
	if err != nil{
		log.Println(err)
		return  nil, err
	}
	return &Account{
		ID:a.ID,
		Name: a.Name,
	}, nil
}

func (r *mutationResolver) createProduct(ctx context.Context, in ProductInput) (*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	p, err := r.server.catalogClient.PostProduct(ctx, in.Name, in.Description, in.Price)
	if err != nil{
		log.Println(err)
		return  nil, err
	}
	return &Product{
		ID: p.ID,
		Name : p.Name,
		Description: p.Description,
		Price : p.Price,
	}, nil
}
func (r *mutationResolver) createOrder(ctx context.Context, in OrderInput) (*Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	var products []Order.OrderedProduct 
	for _, p := range in.Products{
		if p.Quantity <= 0{
			return  nil, ErrInvalidParameter
		}
		products = append(products, Order.OrderProduct{
			ID : p.ID, 
			Quantity : p.Quantity,
		})
	}
	o, err := r.server.orderClient.PostOrder(ctx, in.AccountID, products)
	if err != nil{
		log.Println(err)
		return nil, err
	}
	return &Order{
		ID : o.ID,
		CreatedAt: o.CreatedAt,
		TotalPrice: o.TotalPrice,
	},nil
}