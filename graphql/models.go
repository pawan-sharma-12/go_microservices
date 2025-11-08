package main

type Account struct{
	ID 	  string `json:"id"`
	Name  string `json:"name"`
	Orders	[]Order `json:"orders"`
}
// type Order struct{
// 	ID 	  string   `json:"id"`
// 	CreatedAt string   `json:"createAt"`
// 	TotalPrice float64  `json:"totalPrice"`
// 	Products []OrderedProduct `json:"products"`
// }
// type OrderedProduct struct{
// 	ID 	  string `json:"id"`
// 	Name  string `json:"name"`
// 	Price float64 `json:"price"`
// 	Quantity int    `json:"quantity"`
// 	Description string `json:"description"`
// }