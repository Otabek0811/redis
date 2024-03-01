package postgres

import (
	"app/api/models"
	"app/pkg/helper"
	"context"
	"database/sql"
	"errors"
	"fmt"

	uuid "github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4/pgxpool"
)

type marketRepo struct {
	db *pgxpool.Pool
}

func NewMarketRepo(db *pgxpool.Pool) *marketRepo {
	return &marketRepo{
		db: db,
	}
}

func (r *marketRepo) Create(ctx context.Context, req *models.CreateMarket) (string, error) {

	trx, err := r.db.Begin(ctx)
	if err != nil {
		return "", nil
	}

	defer func() {
		if err != nil {
			trx.Rollback(ctx)
		} else {
			trx.Commit(ctx)
		}
	}()

	var (
		id    = uuid.New().String()
		query string
	)

	query = `
		INSERT INTO market(id, name, address, phone_number, updated_at)
		VALUES ($1, $2, $3, $4, NOW())
	`

	_, err = trx.Exec(ctx, query,
		id,
		req.Name,
		req.Address,
		req.PhoneNumber,
	)

	if err != nil {
		return "", err
	}

	// Market to Product Relation -> Many to Many
	if len(req.ProductIds) > 0 {
		marketProductInsertQuery := `
			INSERT INTO 
				market_product_relation(market_id, product_id) 
			VALUES`

		marketProductInsertQuery, args := helper.InsertMultiple(marketProductInsertQuery, id, req.ProductIds)
		_, err = trx.Exec(ctx, marketProductInsertQuery, args...)
		if err != nil {
			return "", err
		}
	}

	return id, nil
}

func (r *marketRepo) GetByID(ctx context.Context, req *models.MarketPrimaryKey) (*models.Market, error) {

	var (
		query string

		id          sql.NullString
		name        sql.NullString
		address     sql.NullString
		phoneNumber sql.NullString
		createdAt   sql.NullString
		updatedAt   sql.NullString

		productObj pgtype.JSONB
	)

	query = `
		WITH market_product AS (
			SELECT
				JSON_AGG(
					JSON_BUILD_OBJECT (
						'id', p.id,
						'name', p.name,
						'price', p.price,
						'category_id', p.category_id,
						'created_at', p.created_at,
						'updated_at', p.updated_at
					)
				)  AS products,
				mpr.market_id AS market_id

			FROM product AS p
			JOIN market_product_relation AS mpr ON mpr.product_id = p.id
			WHERE mpr.market_id = $1
			GROUP BY mpr.market_id
		)
		SELECT
			m.id,
			m.name,
			m.address,
			m.phone_number,
			m.created_at,
			m.updated_at,

			mp.products
			
		FROM market AS m
		JOIN market_product AS mp ON mp.market_id = m.id
		WHERE m.id =  $1
	`

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&name,
		&address,
		&phoneNumber,
		&createdAt,
		&updatedAt,
		&productObj,
	)

	if err != nil {
		return nil, err
	}

	products := []*models.Product{}
	productObj.AssignTo(&products)

	return &models.Market{
		Id:          id.String,
		Name:        name.String,
		Address:     address.String,
		PhoneNumber: phoneNumber.String,
		CreatedAt:   createdAt.String,
		UpdatedAt:   updatedAt.String,
		Products:    products,
	}, nil
}

func (r *marketRepo) GetList(ctx context.Context, req *models.MarketGetListRequest) (*models.MarketGetListResponse, error) {

	var (
		resp   = &models.MarketGetListResponse{}
		query  string
		where  = " WHERE TRUE"
		offset = " OFFSET 0"
		limit  = " LIMIT 10"
	)

	query = `
		SELECT
			COUNT(*) OVER(),
			id,
			name,
			address,
			phone_number,
			created_at,
			updated_at
		FROM market
	`

	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	if req.Search != "" {
		where += ` AND name ILIKE '%' || '` + req.Search + `' || '%'`
	}

	// if req.UserID != "" {
	// 	where += fmt.Sprintf(" AND user_id = %d", req.UserId)
	// }

	query += where + offset + limit

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			id          sql.NullString
			name        sql.NullString
			address     sql.NullString
			phoneNumber sql.NullString
			createdAt   sql.NullString
			updatedAt   sql.NullString
		)

		err := rows.Scan(
			&resp.Count,
			&id,
			&name,
			&address,
			&phoneNumber,
			&createdAt,
			&updatedAt,
		)

		if err != nil {
			return nil, err
		}

		resp.Markets = append(resp.Markets, &models.Market{
			Id:          id.String,
			Name:        name.String,
			Address:     address.String,
			PhoneNumber: phoneNumber.String,
			CreatedAt:   createdAt.String,
			UpdatedAt:   updatedAt.String,
		})
	}

	return resp, nil
}

func (r *marketRepo) Update(ctx context.Context, req *models.UpdateMarket) (int64, error) {
	trx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, nil
	}

	defer func() {
		if err != nil {
			trx.Rollback(ctx)
		} else {
			trx.Commit(ctx)
		}
	}()

	var (
		query  string
		params map[string]interface{}
	)

	query = `
		UPDATE
			market
		SET
			name = :name,
			address = :address,
			phone_number = :phone_number,
			updated_at = NOW()
		WHERE id = :id
	`

	params = map[string]interface{}{
		"id":           req.Id,
		"name":         req.Name,
		"address":      req.Address,
		"phone_number": req.PhoneNumber,
	}

	query, args := helper.ReplaceQueryParams(query, params)

	result, err := trx.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	// Market to Product Relation -> Many to Many
	if len(req.ProductIds) > 0 {
		var count int
		marketProductRelationCountQuery := `
			SELECT COUNT(*) FROM market_product_relation WHERE market_id = $1 
		`

		err := trx.QueryRow(ctx, marketProductRelationCountQuery, req.Id).Scan(&count)
		if err != nil {
			return 0, err
		}

		if count > 0 {
			marketProductRelationDeleteQuery := `
				DELETE FROM market_product_relation WHERE market_id = $1 
			`

			_, err := trx.Exec(ctx, marketProductRelationDeleteQuery, req.Id)
			if err != nil {
				return 0, err
			}
		}

		marketProductInsertQuery := `
				INSERT INTO 
					market_product_relation(market_id, product_id) 
				VALUES`

		marketProductInsertQuery, args := helper.InsertMultiple(marketProductInsertQuery, req.Id, req.ProductIds)
		_, err = trx.Exec(ctx, marketProductInsertQuery, args...)
		if err != nil {
			return 0, err
		}
	}

	return result.RowsAffected(), nil
}

func (r *marketRepo) Patch(ctx context.Context, req *models.PatchRequest) (int64, error) {

	var (
		query string
		set   string
	)

	if len(req.Fields) <= 0 {
		return 0, errors.New("no fields")
	}

	for key := range req.Fields {
		set += fmt.Sprintf(" %s = :%s, ", key, key)
	}

	query = `
		UPDATE
			market
		SET ` + set + ` updated_at = now()
		WHERE id = :id
	`

	req.Fields["id"] = req.ID

	fmt.Println(query)

	query, args := helper.ReplaceQueryParams(query, req.Fields)
	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

func (r *marketRepo) Delete(ctx context.Context, req *models.MarketPrimaryKey) error {

	_, err := r.db.Exec(ctx, "DELETE FROM market WHERE id = $1", req.Id)
	if err != nil {
		return err
	}

	return nil
}
