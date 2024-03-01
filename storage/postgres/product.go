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

type productRepo struct {
	db *pgxpool.Pool
}

func NewProductRepo(db *pgxpool.Pool) *productRepo {
	return &productRepo{
		db: db,
	}
}

func (r *productRepo) Create(ctx context.Context, req *models.CreateProduct) (string, error) {

	var (
		id    = uuid.New().String()
		query string
	)

	query = `
		INSERT INTO product(id, name, price, category_id, updated_at)
		VALUES ($1, $2, $3, $4, NOW())
	`

	_, err := r.db.Exec(ctx, query,
		id,
		req.Name,
		req.Price,
		helper.NewNullString(req.CategoryId),
	)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *productRepo) GetByID(ctx context.Context, req *models.ProductPrimaryKey) (*models.Product, error) {

	var (
		query string

		id         sql.NullString
		name       sql.NullString
		price      sql.NullFloat64
		categoryID sql.NullString
		createdAt  sql.NullString
		updatedAt  sql.NullString

		categoryObj pgtype.JSONB
	)

	query = `
		SELECT
			p.id,
			p.name,
			p.price,
			p.category_id,
			p.created_at,
			p.updated_at,

			JSON_BUILD_OBJECT (
				'id', c.id,
				'title', c.title,
				'parent_id', c.parent_id,
				'updated_at', c.updated_at,
				'created_at', c.created_at
			) AS category

		FROM product AS p
		LEFT JOIN category AS c ON c.id = p.category_id
		WHERE p.id = $1
	`

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&name,
		&price,
		&categoryID,
		&createdAt,
		&updatedAt,
		&categoryObj,
	)

	if err != nil {
		return nil, err
	}

	category := models.Category{}
	categoryObj.AssignTo(&category)

	return &models.Product{
		Id:           id.String,
		Name:         name.String,
		Price:        price.Float64,
		CategoryId:   categoryID.String,
		CreatedAt:    createdAt.String,
		UpdatedAt:    updatedAt.String,
		CategoryData: &category,
	}, nil
}

func (r *productRepo) GetList(ctx context.Context, req *models.ProductGetListRequest) (*models.ProductGetListResponse, error) {

	var (
		resp   = &models.ProductGetListResponse{}
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
			price,
			category_id,
			created_at,
			updated_at
		FROM product
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
			id         sql.NullString
			name       sql.NullString
			price      sql.NullFloat64
			categoryID sql.NullString
			createdAt  sql.NullString
			updatedAt  sql.NullString
		)

		err := rows.Scan(
			&resp.Count,
			&id,
			&name,
			&price,
			&categoryID,
			&createdAt,
			&updatedAt,
		)

		if err != nil {
			return nil, err
		}

		resp.Products = append(resp.Products, &models.Product{
			Id:         id.String,
			Name:       name.String,
			Price:      price.Float64,
			CategoryId: categoryID.String,
			CreatedAt:  createdAt.String,
			UpdatedAt:  updatedAt.String,
		})
	}

	return resp, nil
}

func (r *productRepo) Update(ctx context.Context, req *models.UpdateProduct) (int64, error) {

	var (
		query  string
		params map[string]interface{}
	)

	query = `
		UPDATE
			product
		SET
			name = :name,
			price = :price,
			category_id = :category_id,
			updated_at = NOW()
		WHERE id = :id
	`

	params = map[string]interface{}{
		"id":          req.Id,
		"name":        req.Name,
		"price":       req.Price,
		"category_id": helper.NewNullString(req.CategoryId),
	}

	query, args := helper.ReplaceQueryParams(query, params)

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

func (r *productRepo) Patch(ctx context.Context, req *models.PatchRequest) (int64, error) {

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
			product
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

func (r *productRepo) Delete(ctx context.Context, req *models.ProductPrimaryKey) error {

	_, err := r.db.Exec(ctx, "DELETE FROM product WHERE id = $1", req.Id)
	if err != nil {
		return err
	}

	return nil
}
