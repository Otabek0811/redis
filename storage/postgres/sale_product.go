package postgres

import (
	"app/api/models"
	"app/pkg/helper"
	"context"
	"database/sql"
	"fmt"

	uuid "github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type saleProduct struct {
	db *pgxpool.Pool
}

func NewSaleProductProductRepo(db *pgxpool.Pool) *saleProduct {
	return &saleProduct{
		db: db,
	}
}

func (r *saleProduct) Create(ctx context.Context, req *models.CreateSaleProduct) (string, error) {

	var (
		id       = uuid.New().String()
		// query    string
		// name     string
		// price    int
		// total    int
		// discount int
		// typeD    int
	)

	// query := `
	// 	INSERT INTO sale_product(id,product_id, product_price, discount, discount_type, price_with_discount, dicsount_price, count, total_price, updated_at)
	// 	VALUES ($1, $2, $3, $4,$5, $6,$7,$8,$9, NOW())
	// `

	// _, err := r.db.Exec(ctx, query,
	// 	id,
	// 	req.UserID,
	// 	req.Count,
	// )

	// if err != nil {
	// 	return "", err
	// }

	return id, nil
}

func (r *saleProduct) GetByID(ctx context.Context, req *models.SaleProductPrimaryKey) (*models.SaleProduct, error) {

	var (
		query string

		id        sql.NullString
		userID    sql.NullString
		total     int
		count     int
		createdAt sql.NullString
		updatedAt sql.NullString
	)

	query = `
		SELECT
			user_id,
			total,
			count,
			created_at,
			updated_at
		FROM sale_product 
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&userID,
		&total,
		&count,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &models.SaleProduct{
		Id:        req.Id,
		// UserID:    userID.String,
		Total:     total,
		Count:     count,
		CreatedAt: createdAt.String,
		UpdatedAt: updatedAt.String,
	}, nil
}

func (r *saleProduct) GetList(ctx context.Context, req *models.SaleProductGetListRequest) (*models.SaleProductGetListResponse, error) {

	var (
		resp   = &models.SaleProductGetListResponse{}
		query  string
		where  = " WHERE TRUE"
		offset = " OFFSET 0"
		limit  = " LIMIT 10"
	)

	query = `
		SELECT
			COUNT(*) OVER(),
			id,
			product_id,
			total,
			count,
			created_at,
			updated_at
		FROM saleProducts
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

	query += where + offset + limit

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id        sql.NullString
			userID    sql.NullString
			total     int
			count     int
			createdAt sql.NullString
			updatedAt sql.NullString
		)

		err := rows.Scan(
			&resp.Count,
			&id,
			&total,
			&count,
			&userID,
			&createdAt,
			&updatedAt,
		)

		if err != nil {
			return nil, err
		}

		resp.SaleProducts = append(resp.SaleProducts, &models.SaleProduct{
			Id:        id.String,
			// UserID:    userID.String,
			Total:     total,
			Count:     count,
			CreatedAt: createdAt.String,
			UpdatedAt: updatedAt.String,
		})
	}

	return resp, nil
}

func (r *saleProduct) Update(ctx context.Context, req *models.UpdateSaleProduct) (int64, error) {

	var (
		query  string
		params map[string]interface{}
	)

	query = `
		UPDATE
			saleProducts
		SET
			
			updated_at = NOW()
		WHERE id = :id
	`

	params = map[string]interface{}{
		"id": req.Id,
	}

	query, args := helper.ReplaceQueryParams(query, params)

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

func (r *saleProduct) Delete(ctx context.Context, req *models.SaleProductPrimaryKey) error {

	_, err := r.db.Exec(ctx, "DELETE FROM saleProducts WHERE id = $1", req.Id)
	if err != nil {
		return err
	}

	return nil
}
