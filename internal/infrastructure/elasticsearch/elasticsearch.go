package es

import (
	"context"
	"sync"

	"github.com/dpe27/es-krake/config"
	"github.com/dpe27/es-krake/pkg/log"
	"github.com/elastic/go-elasticsearch/v9"
)

type ElasticSearch struct {
	mu     sync.RWMutex
	logger *log.Logger
	client *elasticsearch.TypedClient
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
	cli, err := elasticsearch.NewTypedClient(elasticsearch.Config{
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

	e.setCli(cli)
	return nil
}

func (e *ElasticSearch) setCli(newCli *elasticsearch.TypedClient) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.client = newCli
}

func (e *ElasticSearch) cli() *elasticsearch.TypedClient {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.client
}

func (e *ElasticSearch) Ping(ctx context.Context) (bool, error) {
	return e.cli().Ping().Do(ctx)
}
