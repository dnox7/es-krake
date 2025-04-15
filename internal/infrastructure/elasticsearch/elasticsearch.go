package es

import (
	"pech/es-krake/config"
	"pech/es-krake/pkg/log"

	"github.com/elastic/go-elasticsearch/v8"
)

type ElasticSeachService interface {
}

type elasticSeachService struct {
	logger *log.Logger
	client *elasticsearch.Client
}

func NewElasticSearchService(cfg *config.Config) (ElasticSeachService, error) {
	cli, err := elasticsearch.NewClient(elasticsearch.Config{
		// Addresses:         cfg.ES.Addresses,
		// Username:          cfg.ES.Username,
		// Password:          cfg.ES.Password,
		// MaxRetries:        cfg.ES.MaxRetries,
		// EnableMetrics:     cfg.ES.EnableMetrics,
		// EnableDebugLogger: cfg.ES.Debug,
	})
	if err != nil {
		return nil, err
	}

	return &elasticSeachService{
		logger: log.With("service", "elasticsearch"),
		client: cli,
	}, nil
}
