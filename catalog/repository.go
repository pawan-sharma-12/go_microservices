package catalog

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/olivere/elastic/v7"
)

var (
	ErrNotFound = errors.New("entity not found")
)

type Repository interface {
	Close()
	PutProduct(ctx context.Context, product Product) error
	GetProductByID(ctx context.Context, id string) (*Product, error)
	ListProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error)
	ListProductsWithIDs(ctx context.Context, ids []string) ([]Product, error)
	SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error)
}

type elasticRepository struct {
	client *elastic.Client
}

type productDocument struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

func NewElasticRepository(url string) (Repository, error) {
	client, err := elastic.NewClient(elastic.SetURL(url), elastic.SetSniff(false))
	if err != nil {
		return nil, err
	}

	// Ensure the index exists
	exists, err := client.IndexExists("catalog").Do(context.Background())
	if err != nil {
		return nil, err
	}
	if !exists {
		_, err := client.CreateIndex("catalog").BodyString(`{
			"mappings": {
				"properties": {
					"name": {"type": "text"},
					"description": {"type": "text"},
					"price": {"type": "double"}
				}
			}
		}`).Do(context.Background())
		if err != nil {
			return nil, err
		}
	}

	return &elasticRepository{client: client}, nil
}

func (r *elasticRepository) Close() {
	r.client.Stop()
}

func (r *elasticRepository) PutProduct(ctx context.Context, p Product) error {
	_, err := r.client.Index().
		Index("catalog").
		Id(p.ID).
		BodyJson(productDocument{
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		}).
		Do(ctx)
	return err
}

func (r *elasticRepository) GetProductByID(ctx context.Context, id string) (*Product, error) {
	result, err := r.client.Get().
		Index("catalog").
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

func (r *elasticRepository) ListProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error) {
	results, err := r.client.Search().
		Index("catalog").
		Query(elastic.NewMatchAllQuery()).
		From(int(skip)).
		Size(int(take)).
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
	return products, nil
}

func (r *elasticRepository) ListProductsWithIDs(ctx context.Context, ids []string) ([]Product, error) {
	items := []*elastic.MultiGetItem{}
	for _, id := range ids {
		items = append(items, elastic.NewMultiGetItem().Index("catalog").Id(id))
	}

	res, err := r.client.MultiGet().Add(items...).Do(ctx)
	if err != nil {
		return nil, err
	}

	products := []Product{}
	for _, doc := range res.Docs {
		if doc.Found {
			var p productDocument
			if err := json.Unmarshal(doc.Source, &p); err == nil {
				products = append(products, Product{
					ID:          doc.Id,
					Name:        p.Name,
					Description: p.Description,
					Price:       p.Price,
				})
			}
		}
	}
	return products, nil
}

func (r *elasticRepository) SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error) {
	res, err := r.client.Search().
		Index("catalog").
		Query(elastic.NewMultiMatchQuery(query, "name", "description")).
		From(int(skip)).
		Size(int(take)).
		Do(ctx)
	if err != nil {
		log.Println("Error executing search:", err)
		return nil, err
	}

	products := []Product{}
	for _, hit := range res.Hits.Hits {
		var p productDocument
		if err := json.Unmarshal(hit.Source, &p); err == nil {
			products = append(products, Product{
				ID:          hit.Id,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			})
		}
	}

	return products, nil
}
