package repository

import (
	"context"

	"github.com/algorithm9/flash-deal/internal/shared/ent/gen"
	"github.com/algorithm9/flash-deal/internal/shared/ent/gen/user"
	"github.com/algorithm9/flash-deal/internal/shared/idgen"
)

type UserRepository interface {
	GetByID(ctx context.Context, id uint64) (*gen.User, error)
	GetByPhone(ctx context.Context, phone string) (*gen.User, error)
	Create(ctx context.Context, phone, password string) (*gen.User, error)
}

type userRepoImp struct {
	client *gen.Client
	idGen  idgen.IDGenerator
}

func NewUserRepo(client *gen.Client, idGen idgen.IDGenerator) UserRepository {
	return &userRepoImp{client: client, idGen: idGen}
}

func (r *userRepoImp) GetByID(ctx context.Context, id uint64) (*gen.User, error) {
	return r.client.User.Get(ctx, id)
}

func (r *userRepoImp) GetByPhone(ctx context.Context, phone string) (*gen.User, error) {
	return r.client.User.Query().Where(user.Phone(phone)).Only(ctx)
}

func (r *userRepoImp) Create(ctx context.Context, phone, password string) (*gen.User, error) {
	id, err := r.idGen.NextID()
	if err != nil {
		return nil, err
	}
	return r.client.User.
		Create().
		SetID(id).
		SetPhone(phone).
		SetPasswordHash(password).
		Save(ctx)
}
