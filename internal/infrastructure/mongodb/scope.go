package mdb

import (
	"pech/es-krake/internal/domain/shared/scope"

	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type findScope struct {
	builder *options.FindOptionsBuilder
}

func FindScope(builder *options.FindOptionsBuilder) scope.Base {
	return &findScope{builder}
}

func (fs *findScope) GetScope() interface{} {
	return fs.builder
}

type findOneScope struct {
	builder *options.FindOneOptionsBuilder
}

func FindOneScope(builder *options.FindOneOptionsBuilder) scope.Base {
	return &findOneScope{builder}
}

func (fos *findOneScope) GetScope() interface{} {
	return fos.builder
}
