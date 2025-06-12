package es

import (
	"context"
	"fmt"

	"github.com/dpe27/es-krake/pkg/utils"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/bulk"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/count"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/index"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/create"
)

type ElasticSearchRepository interface {
	Search(ctx context.Context, index string, req *search.Request) (*search.Response, error)
	IndexDocument(ctx context.Context, index string, document interface{}) (*index.Response, error)
	Bulk(ctx context.Context, req *bulk.Request) (*bulk.Response, error)
	Count(ctx context.Context, req *count.Request) (*count.Response, error)
	CreateIndex(ctx context.Context, index string, req *create.Request) (*create.Response, error)
	CheckExistIndex(ctx context.Context, index string) (bool, error)
	DeleteIndex(ctx context.Context, index string) error
}

func (e *ElasticSearch) Search(ctx context.Context, index string, req *search.Request) (*search.Response, error) {
	return e.cli().
		Search().
		Index(index).
		Request(req).
		Do(ctx)
}

func (e *ElasticSearch) IndexDocument(ctx context.Context, index string, document interface{}) (*index.Response, error) {
	return e.cli().
		Index(index).
		Request(document).
		Do(ctx)
}

func (e *ElasticSearch) Bulk(ctx context.Context, req *bulk.Request) (*bulk.Response, error) {
	return e.cli().
		Bulk().
		Request(req).
		Do(ctx)
}

func (e *ElasticSearch) Count(ctx context.Context, req *count.Request) (*count.Response, error) {
	return e.cli().
		Count().
		Request(req).
		Do(ctx)
}

func (e *ElasticSearch) CreateIndex(ctx context.Context, index string, req *create.Request) (*create.Response, error) {
	return e.cli().
		Indices.
		Create(index).
		Request(req).
		Do(ctx)
}

func (e *ElasticSearch) CheckExistIndex(ctx context.Context, index string) (bool, error) {
	req := esapi.IndicesExistsRequest{
		Index: []string{index},
	}
	res, err := req.Do(ctx, e.cli().Transport)
	if err != nil {
		return false, nil
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			e.logger.Error(ctx, utils.ErrorCloseResponseBody, "error", err)
		}
	}()

	if res.IsError() {
		return false, fmt.Errorf("failed to check if index %q exists: %s", index, res.String())
	}
	return true, nil
}

func (e *ElasticSearch) DeleteIndex(ctx context.Context, index string) error {
	req := esapi.IndicesDeleteRequest{
		Index: []string{index},
	}
	res, err := req.Do(ctx, e.cli().Transport)
	if err != nil {
		return err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			e.logger.Error(ctx, utils.ErrorCloseResponseBody, "error", err)
		}
	}()

	if res.IsError() {
		return fmt.Errorf("failed to delete index %q: %s", index, res.String())
	}
	return nil
}
