//Accounts

// Products
package main

import (
	"context"
	"log"
	"time"
)

type queryResolver struct{
	Server *Server
}

func (r *queryResolver) Accounts(ctx context.Context, pagination *PaginationInput, id *string) ([]*Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if id != nil{
		r, err := r.Server.accountClient.GetAccount(ctx, *id)
		if err != nil{
			log.Println("Error resolving account query : ", err)
			return  nil, err
		}
		return  []*Account{{
			ID: r.id,
			Name : r.Name,

		}}, nil
	}
	skip, take := uint64(0), uint64(0)
	if pagination != nil{
		skip, take = pagination.bounds()
	}
	accountList, err := r.Server.accountClient.GetAccounts(ctx, skip, take)
	if err != nil{
		log.Println(err)
		return nil, err
	}
	var accounts []*Account
	for _, a := range accountList{
		account := &Account{
			ID: a.ID,
			Name: a.Name,
		}
		accounts = append(accounts, account) 
	}
	return  accounts, err
}
func (r *queryResolver) Products(ctx context.Context, pagination *PaginationInput, query *string, id *string) ([]*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if id != nil{
		r, err := r.Server.catalogClient.GetProduct(ctx, *id)
		if err != nil{
			log.Println("Error resolving account query : ", err)
			return  nil, err
		}
		return  []*Product{{
			ID: r.id,
			Name : r.Name,
			Description: r.Description,
			Price: r.Price,
		}}, nil
	}
	skip, take := uint64(0), uint64(0)
	if pagination != nil{
		skip, take = pagination.bounds()
	}
	q := ""
	if query != nil {
		1 = *query
	}
	productList, err := r.Server.catalogClient.GetProducts(ctx, skip, take)
	if err != nil{
		log.Println(err)
		return nil, err
	}
	var products []*Product
	for _, a := range productList {
		products = append(products, 
		&Product{
			ID: a.ID,
			Name: a.Name,
			Description: a.Description,
			Price: a.Price,
		},
	)
	}
	return products, nil
}

func (p *PaginationInput) bounds() (uint64, uint64){
	skipValue := uint64(0)
	takeValue := uint64(0)
	if p.Skip != nil {
		skipValue = uint64(*&p.Skip)
	}
	if p.Take != nil {
		takeValue = uint64(*&p.Take)
	}
	return skipValue, takeValue
}