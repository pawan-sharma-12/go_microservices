package main

import (
	"context"
	"log"
	"time"
)

// AccountResolver is defined as a struct
type accountResolver struct{
	server *Server
}

func (r *accountResolver) Orders(ctx context.Context, obj *Account) ([]*Order, error) {
	// Implementation goes here
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	orderList , err := r.server.orderClient.GetOrdersForAccount(ctx, obj.ID)
	if err != nil{
		log.Println("Error resolving Orders in account_resolver  : ", err)
		return  nil, err
	}
	var orders []*Order
	for _, o := range orderList {
		var products []*OrderProduct
		for _, p := range o.Products{
			products = append(products, &OrderProduct{
				ID: p.ID,
				Name: p.Name,
				Description: p.Description,
				Price: p.Price,
				Quantity: int(p.Quantity),
			})
		}
		orders = append(orders, &Order{
			ID : o.ID,
			CreatedAt : o.CreatedAt,
			TotalPrice : o.TotalPrice,
			Products : products,
		})
	}
	return orders, nil 
}