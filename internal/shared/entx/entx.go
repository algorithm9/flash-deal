package entx

import (
	"context"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"github.com/algorithm9/flash-deal/internal/model"
	"github.com/algorithm9/flash-deal/internal/shared/ent/gen"
	"github.com/algorithm9/flash-deal/pkg/errorx"
	"github.com/algorithm9/flash-deal/pkg/logger"
)

func NewEntClient(opts *model.DatabaseConfig) (*gen.Client, func(), error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=%s",
		opts.UserName, opts.Password, opts.Host, opts.Port, opts.DBName, "UTC",
	)

	var entOpts []gen.Option
	if opts.ShowSQL {
		entOpts = append(entOpts, gen.Debug(), gen.Log(log.Println))
	}

	entClient, err := gen.Open(opts.Driver, dsn, entOpts...)
	if err != nil {
		panic(err)
	}

	// 清理函数
	cleanup := func() {
		logger.L().Info().Msg("closing ent client")
		if err := entClient.Close(); err != nil {
			logger.L().Err(err).Msgf("ent close error")
		}
	}

	return entClient, cleanup, nil
}

func WithTx(ctx context.Context, client *gen.Client, fn func(tx *gen.Tx) error) error {
	tx, err := client.Tx(ctx)
	if err != nil {
		return errorx.Wrap(errorx.CodeInternal.Int(), http.StatusInternalServerError, "failed to get tx", err)
	}
	defer func() {
		if v := recover(); v != nil {
			tx.Rollback()
			panic(v)
		}
	}()
	if err := fn(tx); err != nil {
		return Rollback(tx, err)
	}
	if err := tx.Commit(); err != nil {
		return errorx.Wrap(errorx.CodeInternal.Int(), http.StatusInternalServerError, "committing transaction", err)
	}
	return nil
}

// Rollback calls to tx.Rollback and wraps the given error with the rollback error if occurred.
func Rollback(tx *gen.Tx, err error) error {
	if rerr := tx.Rollback(); rerr != nil {
		err = fmt.Errorf("%w: rolling back transaction: %v", err, rerr)
	}
	return errorx.Wrap(errorx.CodeInternal.Int(), http.StatusInternalServerError, "rolling back transaction", err)
}
