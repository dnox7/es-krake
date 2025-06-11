package es

import (
	"context"
	"fmt"
	"sync"

	"github.com/dpe27/es-krake/config"
	"github.com/dpe27/es-krake/pkg/log"
	"github.com/dpe27/es-krake/pkg/utils"
	"github.com/elastic/go-elasticsearch/v8"
)

type ElasticSearch struct {
	mu     sync.RWMutex
	logger *log.Logger
	client *elasticsearch.Client
	params esParams
}

type esParams struct {
	addresses         []string
	maxRetries        int
	enableMetrics     bool
	enableDebugLogger bool
}

var (
	esInstance *ElasticSearch
	once       sync.Once
)

func NewElasticSearchRepository(cfg *config.Config, cred *config.ElasticSearchCredentials) *ElasticSearch {
	once.Do(func() {
		esRepo, err := initElasticSearch(cfg, cred)
		if err != nil {
			panic(err)
		}
		esInstance = esRepo
	})
	return esInstance
}

func initElasticSearch(cfg *config.Config, cred *config.ElasticSearchCredentials) (*ElasticSearch, error) {
	esRepo := &ElasticSearch{
		logger: log.With("service", "elasticsearch"),
		params: esParams{
			addresses:         cfg.ES.Addresses,
			maxRetries:        cfg.ES.MaxRetries,
			enableMetrics:     cfg.ES.EnableMetrics,
			enableDebugLogger: cfg.ES.Debug,
		},
	}

	if err := esRepo.Reconn(cred); err != nil {
		return nil, err
	}
	return esRepo, nil
}

func (e *ElasticSearch) Reconn(cred *config.ElasticSearchCredentials) error {
	newCli, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses:         e.params.addresses,
		Username:          cred.Username,
		Password:          cred.Password,
		MaxRetries:        e.params.maxRetries,
		EnableMetrics:     e.params.enableMetrics,
		EnableDebugLogger: e.params.enableDebugLogger,
	})
	if err != nil {
		return err
	}

	e.setCli(newCli)
	return nil
}

func (e *ElasticSearch) setCli(newCli *elasticsearch.Client) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.client = newCli
}

func (e *ElasticSearch) cli() *elasticsearch.Client {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.client
}

func (e *ElasticSearch) Ping(ctx context.Context) error {
	resp, err := e.cli().Ping()
	if err != nil {
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			e.logger.Error(ctx, utils.ErrorCloseResponseBody, "error", err.Error())
		}
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
