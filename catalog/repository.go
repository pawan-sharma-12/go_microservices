package catalog

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	_ "github.com/lib/pq"
	"github.com/olivere/elastic/v7"
)

var (
	ErrNotFound = errors.New("Entitiy not found")
)
type Repository interface {
	Close()
	PutProduct(ctx context.Context, product Product) error
	GetProductByID(ctx context.Context, id string) (*Product, error)
	ListProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error)
	ListProductsWithIDs(ctx context.Context, ids []string) ([]Product, error)
	SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error)
}

type elasticReposityory struct {
	client *elastic.Client

}
type productDocument struct {
	Name 	  string  `json:"name"`
	Description string  `json:"description"`
	Price 	  float64 `json:"price"`
}
func NewElasticRepository(url string) (Repository, error) {
	client, err := elastic.NewClient(elastic.SetURL(url), elastic.SetSniff(false))
	if err != nil {
		return nil, err
	}
	return &elasticReposityory{
		client: client,
	}, nil
}
func (r *elasticReposityory) Close () {
	r.client.Stop()
}
func (r *elasticReposityory) PutProduct(ctx context.Context, p Product) error {
	_, err := r.client.Index().
	Index("catalog").
	Type("product").
	Id(p.ID).BodyJson(productDocument{
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
	}).Do(ctx)
	return err

}
func (r *elasticReposityory) GetProductByID(ctx context.Context, id string) (*Product, error) {
	result, err := r.client.Get().
		Index("catalog").
		Type("product").
		Id(id).
		Do(ctx)
	if err != nil {
		if elastic.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if !result.Found {
		return nil, ErrNotFound
	}
	var doc productDocument
	if err := json.Unmarshal(result.Source, &doc); err != nil {
		return nil, err
	}
	return &Product{
		ID:          result.Id,
		Name:        doc.Name,
		Description: doc.Description,
		Price:       doc.Price,
	}, nil
}
func (r *elasticReposityory) ListProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error) {
	results, err := r.client.Search().
		Index("catalog").
		Type("product").
		Query(elastic.NewMatchAllQuery()).
		From(int(skip)).
		Size(int(take)).
		Sort("id", false).
		Do(ctx)
	if err != nil {
		log.Println("Error executing search:", err)
		return nil, err
	}
	products := make([]Product, 0, len(results.Hits.Hits))
	for _, hit := range results.Hits.Hits {
		var doc productDocument
		if err := json.Unmarshal(hit.Source, &doc); err != nil {
			continue
		}
		products = append(products, Product{
			ID:          hit.Id,
			Name:        doc.Name,
			Description: doc.Description,
			Price:       doc.Price,
		})
	}
	return products, err
}
func (r *elasticReposityory) ListProductsWithIDs(ctx context.Context, ids []string) ([]Product, error) {
	items := []*elastic.MultiGetItem{}
	for _, id := range ids {
		items = append(items, elastic.NewMultiGetItem().Index("catalog").Type("product").Id(id))
	}
	res, err := r.client.MultiGet().Add(items...).Do(ctx)
	if err != nil {
		return nil, err
	}
	products := []Product{}
	for _, doc := range res.Docs{
		p := productDocument{}
		if err := json.Unmarshal(doc.Source, &p); err == nil {
			products = append(products, Product{
				ID:          doc.Id,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			})
		}
	}
	return products, nil
}
func (r *elasticReposityory) SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error) {
	res, err := r.client.Search().Index("catalog").Type("product").Query(elastic.NewMultiMatchQuery(query, "name", "description")).
	From(int(skip)).Size(int(take)).Do(ctx)
	if err != nil {
		log.Println("Error executing search:", err)
		return nil, err
	}
	products := []Product{}
	for _, hit := range res.Hits.Hits {
		p := productDocument{}
		if err := json.Unmarshal(hit.Source, &p); err == nil {
			products = append(products, Product{
				ID:          hit.Id,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			})
		}
	}
	return products, err
}
