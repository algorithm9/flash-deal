package repository

import (
	"context"
	"encoding/json"
	"net/http"

	goodsdto "github.com/algorithm9/flash-deal/internal/module/product/dto"
	"github.com/algorithm9/flash-deal/internal/shared/ent/gen"
	"github.com/algorithm9/flash-deal/pkg/errorx"
)

type GoodsRepository interface {
	GetSKUWithProduct(ctx context.Context, skuID uint64) (*goodsdto.SKUProductDetail, error)
}

type goodsRepo struct {
	db *gen.Client
}

func NewGoodsRepo(client *gen.Client) GoodsRepository {
	return &goodsRepo{db: client}
}

func (r *goodsRepo) GetSKUWithProduct(ctx context.Context, skuID uint64) (*goodsdto.SKUProductDetail, error) {
	skuDetail, err := r.db.QueryContext(ctx, skuWithProductSql(), skuID)
	if err != nil {
		return nil, errorx.Wrap(goodsdto.SKUInternalError.Int(), http.StatusInternalServerError, "failed to query sku", err)
	}
	defer skuDetail.Close()

	var s goodsdto.SKUProductDetail
	var specsRaw []byte
	for skuDetail.Next() {
		if err = skuDetail.Scan(&s.SKUID, &s.ProductID,
			&s.ProductTitle, &s.Description, &specsRaw,
			&s.Price, &s.SeckillPrice, &s.Stock, &s.IsSeckill, &s.IsActive); err != nil {
			return nil, errorx.Wrap(goodsdto.SKUInternalError.Int(), http.StatusInternalServerError, "failed to query sku", err)
		}
	}
	if err = skuDetail.Err(); err != nil {
		return nil, errorx.Wrap(goodsdto.SKUInternalError.Int(), http.StatusInternalServerError, "failed to query sku", err)
	}

	if s.SKUID == 0 {
		return nil, nil
	}
	var specs map[string]string
	if specsRaw != nil {
		if err := json.Unmarshal(specsRaw, &specs); err != nil {
			return nil, err
		}
	}

	s.Specs = specs

	return &s, nil
}

func skuWithProductSql() string {
	return `SELECT 
              s.id AS sku_id,
  			  s.product_id,
  			  p.title,
  			  p.description,
  			  s.specs,  			  
  			  s.price,
  			  s.seckill_price,
  			  s.stock,
  			  s.is_seckill,
  			  s.is_active
			FROM skus s
			JOIN products p ON s.product_id = p.id
			WHERE s.id = ?
`
}
