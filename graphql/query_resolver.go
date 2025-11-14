package main

import (
	"context"
	"log"
	"time"
)

type queryResolver struct {
	server *Server
}

// Accounts resolver
func (r *queryResolver) Accounts(ctx context.Context, pagination PaginationInput, id *string) ([]*Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if id != nil {
		a, err := r.server.accountClient.GetAccount(ctx, *id)
		if err != nil {
			log.Println("Error resolving account query:", err)
			return nil, err
		}
		return []*Account{{
			ID:   a.ID,
			Name: a.Name,
		}}, nil
	}

	// Use pagination directly; no nil check needed
	skip, take := pagination.bounds()

	accountList, err := r.server.accountClient.GetAccounts(ctx, skip, take)
	if err != nil {
		log.Println("Error fetching accounts:", err)
		return nil, err
	}

	accounts := make([]*Account, 0, len(accountList))
	for _, a := range accountList {
		accounts = append(accounts, &Account{
			ID:   a.ID,
			Name: a.Name,
		})
	}

	return accounts, nil
}

// Products resolver
func (r *queryResolver) Products(ctx context.Context, pagination PaginationInput, query *string, id *string) ([]*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if id != nil {
		p, err := r.server.catalogClient.GetProduct(ctx, *id)
		if err != nil {
			log.Println("Error resolving product query:", err)
			return nil, err
		}
		return []*Product{{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		}}, nil
	}

	skip, take := pagination.bounds()

	searchQuery := ""
	if query != nil {
		searchQuery = *query
	}

	productList, err := r.server.catalogClient.GetProducts(ctx, skip, take, []string{}, searchQuery)
	if err != nil {
		log.Println("Error fetching products:", err)
		return nil, err
	}

	products := make([]*Product, 0, len(productList))
	for _, p := range productList {
		products = append(products, &Product{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		})
	}

	return products, nil
}

// Pagination helper (no nil check needed because PaginationInput is a value)
func (p PaginationInput) bounds() (uint64, uint64) {
	return uint64(p.Skip), uint64(p.Take)
}
