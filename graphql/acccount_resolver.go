package main

import "context"

// AccountResolver is defined as a struct
type AccountResolverImpl struct{}

func (r *AccountResolverImpl) Orders(ctx context.Context, obj *Account) ([]*Order, error) {
	// Implementation goes here

	return []*Order{}, nil
}