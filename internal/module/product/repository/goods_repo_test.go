package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/algorithm9/flash-deal/internal/config"
	"github.com/algorithm9/flash-deal/internal/shared/entx"
)

func TestGoodsRepo_GetSKUWithProduct(t *testing.T) {
	cfgPath := "../../../../conf.toml"
	configConfig := config.LoadConfig(cfgPath)
	databaseConfig := config.ProvideDB(configConfig)
	client, cleanup, err := entx.NewEntClient(databaseConfig)
	require.NoError(t, err)
	defer cleanup()
	repo := NewGoodsRepo(client)
	product, err := repo.GetSKUWithProduct(context.Background(), 2)
	require.NoError(t, err)
	fmt.Println(product)
}
