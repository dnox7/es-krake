package mdb

import (
	"fmt"

	"github.com/dpe27/es-krake/internal/domain/shared/specification"
	"github.com/dpe27/es-krake/pkg/utils"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const ErrGetMongoCursor = "Failed to get Mongo Cursor"

type findOpts struct {
	builder *options.FindOptionsBuilder
}

func FindOpts(builder *options.FindOptionsBuilder) specification.Base {
	return &findOpts{builder}
}

func (fs *findOpts) GetSpec() interface{} {
	return fs.builder
}

type findOneOpts struct {
	builder *options.FindOneOptionsBuilder
}

func FindOneOpts(builder *options.FindOneOptionsBuilder) specification.Base {
	return &findOneOpts{builder}
}

func (fos *findOneOpts) GetSpec() interface{} {
	return fos.builder
}

type OptionBuilder interface {
	~*options.FindOptionsBuilder | ~*options.FindOneOptionsBuilder
}

func ToOptsBuilder[T OptionBuilder](spec specification.Base) (T, error) {
	if spec == nil {
		return nil, nil
	}
	opts, ok := spec.GetSpec().(T)
	if !ok {
		return nil, fmt.Errorf(utils.ErrorGetSpec)
	}
	return opts, nil
}
