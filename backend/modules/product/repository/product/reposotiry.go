package product

import (
	"context"
	"database/sql"
	"diploma/modules/product/model"
	"diploma/modules/product/repository/product/converter"
	repoModel "diploma/modules/product/repository/product/model"
	"diploma/pkg/client/db"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
)

const (
	// ======== products table ========
	productsTbl = "products"

	pIdCol               = "id"
	pNameCol             = "name"
	pImageUrlCol         = "image_url"
	pGTINCol             = "gtin"
	pLowestSupplierIDCol = "lowest_supplier_id"
	pCreatedAtCol        = "created_at"
	pUpdatedAtCol        = "updated_at"

	// ======== product-supplier table ========
	productsSupplierTbl = "products_supplier"
	psProductIDCol      = "product_id"
	psSupplierIDCol     = "supplier_id"
	psPriceCol          = "price"
	psSellAmountCol     = "sell_amount"

	// ======== supplier table ========
	supplierTbl             = "suppliers"
	sIDCol                  = "user_id"
	sNameCol                = "name"
	sOrderAmountCol         = "order_amount"
	sDeliveryConditionIDCol = "condition_id"

	// ======== delivery conditions ========
	deliveryConditionTbl    = "delivery_conditions"
	dcIDCol                 = "condition_id"
	dcFreeDeliveryAmountCol = "minimum_free_delivery_amount"
	dcDeliveryFeeCol        = "delivery_fee"
	usersTbl                = "users"
	uIDCol                  = "id"
	uNameCol                = "name"
)

const (
	ratingAverageExpr = "COALESCE((SELECT ROUND(AVG(r.rating)::numeric, 1) FROM product_reviews r WHERE r.product_id = p.id), 0)"
	ratingCountExpr   = "COALESCE((SELECT COUNT(*) FROM product_reviews r WHERE r.product_id = p.id), 0)"
	emptySupplierName = "COALESCE(s.name, '')"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) *repo {
	return &repo{db: db}
}

// GetProduct retrieves a product by its ID
func (r *repo) GetProduct(ctx context.Context, id int64) (*model.Product, error) {
	builder := sq.
		Select(
			"p."+pIdCol,
			"p."+pNameCol,
			"p."+pImageUrlCol,
			"p."+pGTINCol,
			"p."+pCreatedAtCol,
			"p."+pUpdatedAtCol,
			ratingAverageExpr+" AS rating_average",
			ratingCountExpr+" AS rating_count",
		).
		From(productsTbl + " AS p").
		Where(sq.Eq{"p." + pIdCol: id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build get product query: %w", err)
	}

	q := db.Query{
		Name:     "product_repository.GetProduct",
		QueryRaw: query,
	}

	var product repoModel.Product
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(
		&product.ID,
		&product.Name,
		&product.ImageUrl,
		&product.GTIN,
		&product.CreatedAt,
		&product.UpdatedAt,
		&product.RatingAverage,
		&product.RatingCount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, model.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return converter.ToProductFromRepo(product), nil
}

// GetProductListByIDList retrieves a list of products by their IDs
func (r *repo) GetProductListByIDList(ctx context.Context, idList []int64) ([]*model.Product, error) {
	builder := sq.
		Select(
			"p."+pIdCol,
			"p."+pNameCol,
			"p."+pImageUrlCol,
			"p."+pGTINCol,
			"p."+pCreatedAtCol,
			"p."+pUpdatedAtCol,
			ratingAverageExpr+" AS rating_average",
			ratingCountExpr+" AS rating_count",
		).
		From(productsTbl + " AS p").
		Where(sq.Eq{"p." + pIdCol: idList}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build get product list query: %w", err)
	}

	q := db.Query{
		Name:     "product_repository.GetProductListByIDList",
		QueryRaw: query,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get product list: %w", err)
	}
	defer rows.Close()

	var products []repoModel.Product
	for rows.Next() {
		var product repoModel.Product
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.ImageUrl,
			&product.GTIN,
			&product.CreatedAt,
			&product.UpdatedAt,
			&product.RatingAverage,
			&product.RatingCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating products: %w", err)
	}

	var result []*model.Product
	for _, product := range products {
		result = append(result, converter.ToProductFromRepo(product))
	}

	return result, nil
}

// GetSupplierProductListByProduct retrieves a list of suppliers for a specific product
func (r *repo) GetSupplierProductListByProduct(ctx context.Context, id int64) ([]model.ProductSupplier, error) {
	builder := sq.
		Select(
			"ps."+psSupplierIDCol+" AS supplier_id",
			"ps."+psPriceCol+" AS price",
			"ps."+psSellAmountCol+" AS sell_amount",
			"s."+sNameCol+" AS supplier_name",
			"s."+sOrderAmountCol+" AS order_amount",
			"dc."+dcFreeDeliveryAmountCol+" AS minimum_free_delivery_amount",
			"dc."+dcDeliveryFeeCol+" AS delivery_fee",
		).
		From(productsSupplierTbl + " AS ps").
		Join(supplierTbl + " AS s ON s." + sIDCol + " = ps." + psSupplierIDCol).
		LeftJoin(deliveryConditionTbl + " AS dc ON dc." + dcIDCol + " = s." + sDeliveryConditionIDCol).
		Where(sq.Eq{psProductIDCol: id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build get supplier product list query: %w", err)
	}

	q := db.Query{
		Name:     "product_repository.GetSupplierProductListByProduct",
		QueryRaw: query,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get supplier product list: %w", err)
	}
	defer rows.Close()

	var results []model.ProductSupplier
	for rows.Next() {
		var ps model.ProductSupplier
		var s model.Supplier
		if err := rows.Scan(
			&s.ID,
			&ps.Price,
			&ps.SellAmount,
			&s.Name,
			&s.OrderAmount,
			&s.FreeDeliveryAmount,
			&s.DeliveryFee,
		); err != nil {
			return nil, fmt.Errorf("failed to scan supplier product: %w", err)
		}
		ps.Supplier = s
		results = append(results, ps)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating supplier products: %w", err)
	}

	return results, nil
}

// GetProductList retrieves a list of products with their lowest supplier info
func (r *repo) GetProductList(ctx context.Context, queryParam *model.ProductListQuery) ([]model.Product, error) {
	builder := sq.Select(
		"p."+pIdCol+" AS product_id",
		"p."+pNameCol+" AS product_name",
		"p."+pImageUrlCol+" AS product_image_url",
		"p."+pGTINCol+" AS product_gtin",
		ratingAverageExpr + " AS rating_average",
		ratingCountExpr + " AS rating_count",
		"COALESCE(ps."+psPriceCol+", 0) AS ps_price",
		"COALESCE(ps."+psSellAmountCol+", 0) AS ps_sell_amount",
		"COALESCE(ps."+psSupplierIDCol+", 0) AS ps_supplier_id",
		emptySupplierName + " AS supplier_name",
		"COALESCE(s."+sOrderAmountCol+", 0) AS supplier_order_amount",
		"COALESCE(dc."+dcFreeDeliveryAmountCol+", 0) AS dc_min_free_delivery_amount",
		"COALESCE(dc."+dcDeliveryFeeCol+", 0) AS dc_delivery_fee",
	).
		From(productsTbl + " AS p").
		LeftJoin(productsSupplierTbl + " AS ps ON ps." + psProductIDCol + " = p." + pIdCol +
			" AND ps." + psSupplierIDCol + " = p." + pLowestSupplierIDCol).
		LeftJoin(supplierTbl + " AS s ON s." + sIDCol + " = ps." + psSupplierIDCol).
		LeftJoin(deliveryConditionTbl + " AS dc ON dc." + dcIDCol + " = s." + sDeliveryConditionIDCol).
		PlaceholderFormat(sq.Dollar)

	builder = applyProductListFilters(builder, queryParam)
	builder = applyProductListSorting(builder, queryParam)

	if queryParam.Limit > 0 {
		builder = builder.Limit(uint64(queryParam.Limit))
	} else {
		builder = builder.Limit(30)
	}
	if queryParam.Offset > 0 {
		builder = builder.Offset(uint64(queryParam.Offset))
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build get product list query: %w", err)
	}

	q := db.Query{
		Name:     "product_repository.GetProductList",
		QueryRaw: query,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get product list: %w", err)
	}
	defer rows.Close()

	var productList []repoModel.Product
	for rows.Next() {
		var (
			p  repoModel.Product
			ps repoModel.ProductSupplier
			s  repoModel.Supplier
		)

		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.ImageUrl,
			&p.GTIN,
			&p.RatingAverage,
			&p.RatingCount,
			&ps.Price,
			&ps.SellAmount,
			&s.ID,
			&s.Name,
			&s.OrderAmount,
			&s.FreeDeliveryAmount,
			&s.DeliveryFee,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}

		ps.Supplier = s
		p.LowestSupplier = ps
		productList = append(productList, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating products: %w", err)
	}

	return converter.ToProductListFromRepo(productList), nil
}

func (r *repo) GetReviewsByProduct(ctx context.Context, productID int64) ([]model.Review, error) {
	builder := sq.
		Select(
			"pr.id",
			"pr.user_id",
			"u."+uNameCol,
			"pr.rating",
			"pr.comment",
			"pr.created_at",
		).
		From("product_reviews AS pr").
		Join(usersTbl + " AS u ON u." + uIDCol + " = pr.user_id").
		Where(sq.Eq{"pr.product_id": productID}).
		OrderBy("pr.created_at DESC").
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build get reviews query: %w", err)
	}

	q := db.Query{Name: "product_repository.GetReviewsByProduct", QueryRaw: query}
	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get reviews: %w", err)
	}
	defer rows.Close()

	var reviews []model.Review
	for rows.Next() {
		var review model.Review
		if err := rows.Scan(&review.ID, &review.UserID, &review.UserName, &review.Rating, &review.Comment, &review.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan review: %w", err)
		}
		reviews = append(reviews, review)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating reviews: %w", err)
	}

	return reviews, nil
}

func (r *repo) UpsertReview(ctx context.Context, input *model.ReviewInput) error {
	builder := sq.
		Insert("product_reviews").
		Columns("user_id", "product_id", "rating", "comment").
		Values(input.UserID, input.ProductID, input.Rating, input.Comment).
		Suffix("ON CONFLICT (user_id, product_id) DO UPDATE SET rating = EXCLUDED.rating, comment = EXCLUDED.comment, updated_at = CURRENT_TIMESTAMP").
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build upsert review query: %w", err)
	}

	q := db.Query{Name: "product_repository.UpsertReview", QueryRaw: query}
	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("failed to save review: %w", err)
	}

	return nil
}

func (r *repo) GetFavoriteProductIDsByUser(ctx context.Context, userID int64) ([]int64, error) {
	builder := sq.
		Select("product_id").
		From("favorite_products").
		Where(sq.Eq{"user_id": userID}).
		OrderBy("created_at DESC").
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build favorite ids query: %w", err)
	}

	q := db.Query{Name: "product_repository.GetFavoriteProductIDsByUser", QueryRaw: query}
	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get favorite ids: %w", err)
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan favorite id: %w", err)
		}
		ids = append(ids, id)
	}

	return ids, rows.Err()
}

func (r *repo) IsFavorite(ctx context.Context, userID, productID int64) (bool, error) {
	builder := sq.
		Select("COUNT(*) > 0").
		From("favorite_products").
		Where(sq.Eq{"user_id": userID, "product_id": productID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return false, fmt.Errorf("failed to build favorite status query: %w", err)
	}

	q := db.Query{Name: "product_repository.IsFavorite", QueryRaw: query}
	var isFavorite bool
	if err := r.db.DB().QueryRowContext(ctx, q, args...).Scan(&isFavorite); err != nil {
		return false, fmt.Errorf("failed to get favorite status: %w", err)
	}
	return isFavorite, nil
}

func (r *repo) AddFavorite(ctx context.Context, userID, productID int64) error {
	builder := sq.
		Insert("favorite_products").
		Columns("user_id", "product_id").
		Values(userID, productID).
		Suffix("ON CONFLICT (user_id, product_id) DO NOTHING").
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build add favorite query: %w", err)
	}

	q := db.Query{Name: "product_repository.AddFavorite", QueryRaw: query}
	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("failed to add favorite: %w", err)
	}
	return nil
}

func (r *repo) RemoveFavorite(ctx context.Context, userID, productID int64) error {
	builder := sq.
		Delete("favorite_products").
		Where(sq.Eq{"user_id": userID, "product_id": productID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build remove favorite query: %w", err)
	}

	q := db.Query{Name: "product_repository.RemoveFavorite", QueryRaw: query}
	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("failed to remove favorite: %w", err)
	}
	return nil
}

// GetTotalProducts returns the total number of products
func (r *repo) GetTotalProducts(ctx context.Context, queryParam *model.ProductListQuery) (int, error) {
	builder := sq.
		Select("COUNT(*)").
		From(productsTbl + " AS p").
		LeftJoin(productsSupplierTbl + " AS ps ON ps." + psProductIDCol + " = p." + pIdCol +
			" AND ps." + psSupplierIDCol + " = p." + pLowestSupplierIDCol).
		PlaceholderFormat(sq.Dollar)

	builder = applyProductListFilters(builder, queryParam)

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build get total products query: %w", err)
	}

	q := db.Query{
		Name:     "product_repository.GetTotalProducts",
		QueryRaw: query,
	}

	var total int
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to get total products: %w", err)
	}

	return total, nil
}

func (r *repo) GetProductSuggestions(ctx context.Context, queryParam *model.SuggestionQuery) ([]string, error) {
	builder := sq.
		Select(pNameCol).
		From(productsTbl).
		Where("LOWER(" + pNameCol + ") LIKE LOWER(?)", "%"+strings.TrimSpace(queryParam.Query)+"%").
		GroupBy(pNameCol).
		OrderBy("CHAR_LENGTH(" + pNameCol + ") ASC", pNameCol+" ASC").
		Limit(uint64(queryParam.Limit)).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build product suggestions query: %w", err)
	}

	q := db.Query{
		Name:     "product_repository.GetProductSuggestions",
		QueryRaw: query,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get product suggestions: %w", err)
	}
	defer rows.Close()

	var suggestions []string
	for rows.Next() {
		var suggestion string
		if err := rows.Scan(&suggestion); err != nil {
			return nil, fmt.Errorf("failed to scan product suggestion: %w", err)
		}
		suggestions = append(suggestions, suggestion)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating product suggestions: %w", err)
	}

	return suggestions, nil
}

func (r *repo) CreateProduct(ctx context.Context, input *model.AdminProductInput) (*model.Product, error) {
	builder := sq.Insert(productsTbl).
		PlaceholderFormat(sq.Dollar).
		Columns(pNameCol, pImageUrlCol, pGTINCol).
		Values(strings.TrimSpace(input.Name), strings.TrimSpace(input.ImageUrl), input.GTIN).
		Suffix("RETURNING " + pIdCol + ", " + pNameCol + ", " + pImageUrlCol + ", " + pGTINCol + ", " + pCreatedAtCol + ", " + pUpdatedAtCol)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build create product query: %w", err)
	}

	q := db.Query{Name: "product_repository.CreateProduct", QueryRaw: query}
	var product repoModel.Product
	if err := r.db.DB().QueryRowContext(ctx, q, args...).Scan(
		&product.ID,
		&product.Name,
		&product.ImageUrl,
		&product.GTIN,
		&product.CreatedAt,
		&product.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return converter.ToProductFromRepo(product), nil
}

func (r *repo) UpdateProduct(ctx context.Context, input *model.AdminProductInput) (*model.Product, error) {
	builder := sq.Update(productsTbl).
		PlaceholderFormat(sq.Dollar).
		Set(pNameCol, strings.TrimSpace(input.Name)).
		Set(pImageUrlCol, strings.TrimSpace(input.ImageUrl)).
		Set(pGTINCol, input.GTIN).
		Set(pUpdatedAtCol, sq.Expr("NOW()")).
		Where(sq.Eq{pIdCol: input.ID}).
		Suffix("RETURNING " + pIdCol + ", " + pNameCol + ", " + pImageUrlCol + ", " + pGTINCol + ", " + pCreatedAtCol + ", " + pUpdatedAtCol)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build update product query: %w", err)
	}

	q := db.Query{Name: "product_repository.UpdateProduct", QueryRaw: query}
	var product repoModel.Product
	if err := r.db.DB().QueryRowContext(ctx, q, args...).Scan(
		&product.ID,
		&product.Name,
		&product.ImageUrl,
		&product.GTIN,
		&product.CreatedAt,
		&product.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return converter.ToProductFromRepo(product), nil
}

func (r *repo) DeleteProduct(ctx context.Context, productID int64) error {
	queries := []db.Query{
		{Name: "product_repository.DeleteProductSuppliers", QueryRaw: "DELETE FROM " + productsSupplierTbl + " WHERE " + psProductIDCol + " = $1"},
		{Name: "product_repository.DeleteProductFavorites", QueryRaw: "DELETE FROM favorite_products WHERE product_id = $1"},
		{Name: "product_repository.DeleteProductReviews", QueryRaw: "DELETE FROM product_reviews WHERE product_id = $1"},
		{Name: "product_repository.DeleteProduct", QueryRaw: "DELETE FROM " + productsTbl + " WHERE " + pIdCol + " = $1"},
	}

	for _, q := range queries {
		if _, err := r.db.DB().ExecContext(ctx, q, productID); err != nil {
			return fmt.Errorf("failed to execute %s: %w", q.Name, err)
		}
	}

	return nil
}

// GetProductPriceBySupplier retrieves the price of a product from a specific supplier
func (r *repo) GetProductPriceBySupplier(ctx context.Context, productID, supplierID int64) (int, error) {
	builder := sq.
		Select(psPriceCol).
		From(productsSupplierTbl).
		Where(sq.And{
			sq.Eq{psProductIDCol: productID},
			sq.Eq{psSupplierIDCol: supplierID},
		}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build get product price query: %w", err)
	}

	q := db.Query{
		Name:     "product_repository.GetProductPriceBySupplier",
		QueryRaw: query,
	}

	var price int
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&price)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, model.ErrNoRows
		}
		return 0, fmt.Errorf("failed to get product price: %w", err)
	}
	return price, nil
}

func applyProductListFilters(builder sq.SelectBuilder, queryParam *model.ProductListQuery) sq.SelectBuilder {
	if queryParam == nil {
		return builder
	}

	if queryParam.Search != "" {
		builder = builder.Where("LOWER(p."+pNameCol+") LIKE LOWER(?)", "%"+strings.TrimSpace(queryParam.Search)+"%")
	}

	if queryParam.MinPrice > 0 {
		builder = builder.Where("COALESCE(ps."+psPriceCol+", 0) >= ?", queryParam.MinPrice)
	}

	if queryParam.MaxPrice > 0 {
		builder = builder.Where("COALESCE(ps."+psPriceCol+", 0) <= ?", queryParam.MaxPrice)
	}

	return builder
}

func applyProductListSorting(builder sq.SelectBuilder, queryParam *model.ProductListQuery) sq.SelectBuilder {
	sortBy := "p." + pIdCol
	sortOrder := "DESC"

	if queryParam != nil {
		switch strings.ToLower(queryParam.SortBy) {
		case "name":
			sortBy = "p." + pNameCol
		case "price":
			sortBy = "COALESCE(ps." + psPriceCol + ", 0)"
		case "newest":
			sortBy = "p." + pIdCol
		}

		if strings.ToLower(queryParam.SortOrder) == "asc" {
			sortOrder = "ASC"
		}
	}

	return builder.OrderBy(sortBy+" "+sortOrder, "p."+pIdCol+" DESC")
}
