package repository

import (
	"context"
	"pech/es-krake/internal/domain/product/entity"
	domainRepo "pech/es-krake/internal/domain/product/repository"
	"pech/es-krake/internal/infrastructure/db"
	"pech/es-krake/pkg/log"
	"pech/es-krake/pkg/utils"
)

type productAttributeValueRepository struct {
	logger *log.Logger
	pg     *db.PostgreSQL
}

func NewProductAttributeValueRepository(pg *db.PostgreSQL) domainRepo.ProductAttributeValueRepository {
	return &productAttributeValueRepository{
		logger: log.With("repo", "product_attribute_value_repo"),
		pg:     pg,
	}
}

// Create implements repository.ProductAttributeValueRepository.
func (p *productAttributeValueRepository) Create(ctx context.Context, attributes map[string]interface{}) (entity.ProductAttributeValue, error) {
	panic("unimplemented")
}

// CreateBatch implements repository.ProductAttributeValueRepository.
func (p *productAttributeValueRepository) CreateBatch(ctx context.Context, attributeValues []map[string]interface{}, batchSize int) ([]entity.ProductAttributeValue, error) {
	panic("unimplemented")
}

// CreateWithTx implements repository.ProductAttributeValueRepository.
func (p *productAttributeValueRepository) CreateWithTx(ctx context.Context, attributes map[string]interface{}) (entity.ProductAttributeValue, error) {
	panic("unimplemented")
}

// FindByConditions implements repository.ProductAttributeValueRepository.
func (p *productAttributeValueRepository) FindByConditions(ctx context.Context, conditions map[string]interface{}, scopes ...utils.Scope) ([]entity.ProductAttributeValue, error) {
	panic("unimplemented")
}

// TakeByConditions implements repository.ProductAttributeValueRepository.
func (p *productAttributeValueRepository) TakeByConditions(ctx context.Context, conditions map[string]interface{}, scopes ...utils.Scope) (entity.ProductAttributeValue, error) {
	panic("unimplemented")
}

// UpdateWithTx implements repository.ProductAttributeValueRepository.
func (p *productAttributeValueRepository) UpdateWithTx(ctx context.Context, attributeValue entity.ProductAttributeValue, attributesToUpdate map[string]interface{}) (entity.ProductAttributeValue, error) {
	panic("unimplemented")
}

// func (r *productAttributeValueRepository) TakeByConditions(ctx context.Context, conditions map[string]interface{}, scopes ...utils.Scope) (entity.ProductAttributeValue, error) {
// 	var prodAttrVal entity.ProductAttributeValue
//
// 	query, args, err := r.buildSelectQuery(conditions)
// 	if err != nil {
// 		r.logger.Error(ctx, utils.ErrQueryBuilderFailedMsg)
// 		return prodAttrVal, err
// 	}
//
// 	row := r.pg.DB.QueryRowxContext(ctx, query, args...)
// 	err = row.Err()
// 	if err == sql.ErrNoRows {
// 		err = domain.ErrRecordNotFound
// 	}
//
// 	if err != nil {
// 		return prodAttrVal, err
// 	}
//
// 	prodAttrVal.Attribute = &entity.Attribute{}
// 	prodAttrVal.Type = &entity.AttributeType{}
//
// 	err = row.Scan(
// 		&prodAttrVal.ID,
// 		&prodAttrVal.ProductID,
// 		&prodAttrVal.AttributeID,
// 		&prodAttrVal.Value,
// 		&prodAttrVal.CreatedAt,
// 		&prodAttrVal.UpdatedAt,
// 		&prodAttrVal.Attribute.ID,
// 		&prodAttrVal.Attribute.Name,
// 		&prodAttrVal.Type.ID,
// 		&prodAttrVal.Type.Name,
// 	)
// 	return prodAttrVal, err
//
// }
//
// func (r *productAttributeValueRepository) FindByConditions(
// 	ctx context.Context,
// 	conditions map[string]interface{},
// 	scopes ...utils.Scope,
// ) ([]entity.ProductAttributeValue, error) {
// 	var pavList []entity.ProductAttributeValue
//
// 	query, args, err := r.buildSelectQuery(conditions)
// 	if err != nil {
// 		r.logger.Error(ctx, utils.ErrQueryBuilderFailedMsg)
// 		return nil, err
// 	}
//
// 	rows, err := r.pg.DB.QueryxContext(ctx, query, args...)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	for rows.Next() {
// 		pav := entity.ProductAttributeValue{
// 			Attribute: &entity.Attribute{},
// 			Type:      &entity.AttributeType{},
// 		}
//
// 		if err := rows.Scan(
// 			&pav.ID,
// 			&pav.ProductID,
// 			&pav.AttributeID,
// 			&pav.Value,
// 			&pav.CreatedAt,
// 			&pav.UpdatedAt,
// 			&pav.Attribute.ID,
// 			&pav.Attribute.Name,
// 			&pav.Type.ID,
// 			&pav.Type.Name,
// 		); err != nil {
// 			return nil, err
// 		}
//
// 		pavList = append(pavList, pav)
// 	}
// 	return pavList, nil
// }
//
// func (r *productAttributeValueRepository) buildSelectQuery(
// 	conditions map[string]interface{},
// ) (string, []interface{}, error) {
// 	return r.pg.Builder.
// 		Select(
// 			"pav.id",
// 			"pav.product_id",
// 			"pav.attribute_id",
// 			"pav.value",
// 			"pav.created_at",
// 			"pav.updated_at",
// 			"attr.id",
// 			"attr.name",
// 			"aty.id",
// 			"aty.name",
// 		).
// 		From(domainRepo.ProductAttributeValueTableName + " AS pav").
// 		InnerJoin(domainRepo.AttributeTableName + " AS attr ON attr.id = pav.attribute_id").
// 		InnerJoin(domainRepo.AttributeTypeTableName + " AS aty ON aty.id = attr.attribute_type_id").
// 		Where(sq.Eq(conditions)).
// 		ToSql()
// }
//
// func (r *productAttributeValueRepository) Create(
// 	ctx context.Context,
// 	attributes map[string]interface{},
// ) (entity.ProductAttributeValue, error) {
// 	var pav entity.ProductAttributeValue
//
// 	if err := utils.MapToStruct(attributes, &pav); err != nil {
// 		return pav, err
// 	}
//
// 	sql, args, err := r.buildInsertSingleQuery(&pav)
// 	if err != nil {
// 		r.logger.Error(ctx, utils.ErrQueryBuilderFailedMsg)
// 		return pav, err
// 	}
//
// 	err = r.pg.DB.QueryRowxContext(ctx, sql, args...).StructScan(&pav)
// 	return pav, err
// }
//
// func (r *productAttributeValueRepository) CreateWithTx(
// 	ctx context.Context,
// 	attributes map[string]interface{},
// ) (entity.ProductAttributeValue, error) {
// 	var pav entity.ProductAttributeValue
//
// 	if err := utils.MapToStruct(attributes, &pav); err != nil {
// 		return pav, err
// 	}
//
// 	sql, args, err := r.buildInsertSingleQuery(&pav)
// 	if err != nil {
// 		r.logger.Error(ctx, utils.ErrQueryBuilderFailedMsg)
// 		return pav, err
// 	}
//
// 	err = utils.SqlTransaction(ctx, r.logger, r.pg.DB, nil, func(tx *sqlx.Tx) error {
// 		return tx.QueryRowxContext(ctx, sql, args...).StructScan(&pav)
// 	})
// 	return pav, err
// }
//
// func (r *productAttributeValueRepository) buildInsertSingleQuery(att *entity.ProductAttributeValue) (string, []interface{}, error) {
// 	return r.pg.Builder.
// 		Insert(domainRepo.ProductAttributeValueTableName).
// 		Columns("product_id", "attribute_id", "value").
// 		Values(att.ProductID, att.AttributeID, att.Value).
// 		Suffix("RETURNING *").
// 		ToSql()
// }
//
// func (r *productAttributeValueRepository) CreateBatch(
// 	ctx context.Context,
// 	attributeValues []map[string]interface{},
// 	batchSize int,
// ) ([]entity.ProductAttributeValue, error) {
// 	if batchSize == 0 {
// 		batchSize = len(attributeValues)
// 	}
// 	k := int(math.Ceil(float64(len(attributeValues)) / float64(batchSize)))
//
// 	var (
// 		wg      sync.WaitGroup
// 		mu      sync.Mutex
// 		results []entity.ProductAttributeValue
// 		errs    []error
// 		ch      = make(chan []entity.ProductAttributeValue, k)
// 		errCh   = make(chan error, k)
// 	)
//
// 	for i := 0; i < k; i++ {
// 		start := i * batchSize
// 		end := start + batchSize
// 		if end > len(attributeValues) {
// 			end = len(attributeValues)
// 		}
//
// 		batch := attributeValues[start:end]
//
// 		wg.Add(1)
// 		go func(batch []map[string]interface{}) {
// 			defer wg.Done()
//
// 			insertedValues, err := r.insertBatch(ctx, batch)
// 			if err != nil {
// 				errCh <- err
// 				return
// 			}
// 			ch <- insertedValues
// 		}(batch)
// 	}
//
// 	go func() {
// 		wg.Wait()
// 		close(ch)
// 		close(errCh)
// 	}()
//
// 	for res := range ch {
// 		mu.Lock()
// 		results = append(results, res...)
// 		mu.Unlock()
// 	}
//
// 	for err := range errCh {
// 		mu.Lock()
// 		errs = append(errs, err)
// 		mu.Unlock()
// 	}
//
// 	if len(errs) > 0 {
// 		r.logger.Error(ctx, "batch insert failed", "detail", errs)
// 		return nil, errs[0]
// 	}
// 	return results, nil
// }
//
// func (r *productAttributeValueRepository) insertBatch(
// 	ctx context.Context,
// 	batch []map[string]interface{},
// ) ([]entity.ProductAttributeValue, error) {
// 	builder := r.pg.Builder.
// 		Insert(domainRepo.ProductAttributeValueTableName).
// 		Columns("product_id", "attribute_id", "product_option_id", "value")
//
// 	for _, input := range batch {
// 		var pav entity.ProductAttributeValue
// 		if err := utils.MapToStruct(input, &pav); err != nil {
// 			return nil, err
// 		}
//
// 		builder = builder.Values(pav.ProductID, pav.AttributeID, pav.ProductOptionID, pav.Value)
// 	}
//
// 	sql, args, err := builder.Suffix("RETURNING *").ToSql()
// 	if err != nil {
// 		r.logger.Error(ctx, utils.ErrQueryBuilderFailedMsg)
// 		return nil, err
// 	}
//
// 	rows, err := r.pg.DB.QueryxContext(ctx, sql, args...)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	var results []entity.ProductAttributeValue
// 	for rows.Next() {
// 		pav := entity.ProductAttributeValue{}
// 		err = rows.Scan(&pav)
// 		if err != nil {
// 			return nil, err
// 		}
// 		results = append(results, pav)
// 	}
//
// 	return results, nil
// }
//
// func (r *productAttributeValueRepository) UpdateWithTx(
// 	ctx context.Context,
// 	pav entity.ProductAttributeValue,
// 	attributesToUpdate map[string]interface{},
// ) (entity.ProductAttributeValue, error) {
// 	if err := utils.MapToStruct(attributesToUpdate, &pav); err != nil {
// 		return pav, err
// 	}
//
// 	sql, args, err := r.pg.Builder.
// 		Update(domainRepo.ProductAttributeValueTableName).
// 		SetMap(map[string]interface{}{
// 			"product_id":   pav.ProductID,
// 			"attribute_id": pav.AttributeID,
// 			"value":        pav.Value,
// 		}).
// 		Where(sq.Eq{"id": pav.ID}).
// 		Suffix("RETURNING *").
// 		ToSql()
// 	if err != nil {
// 		r.logger.Error(ctx, utils.ErrQueryBuilderFailedMsg)
// 		return pav, err
// 	}
//
// 	err = utils.SqlTransaction(ctx, r.logger, r.pg.DB, nil, func(tx *sqlx.Tx) error {
// 		return tx.QueryRowxContext(ctx, sql, args...).StructScan(&pav)
// 	})
//
// 	return pav, err
// }
