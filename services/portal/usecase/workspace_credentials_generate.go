package usecase

import (
	"context"
	"fmt"

	"github.com/kanthorlabs/kanthor/internal/constants"
	"github.com/kanthorlabs/kanthor/internal/entities"
	"github.com/kanthorlabs/kanthor/pkg/identifier"
	"github.com/kanthorlabs/kanthor/pkg/utils"
	"github.com/kanthorlabs/kanthor/pkg/validator"
	"github.com/kanthorlabs/kanthor/project"
)

type WorkspaceCredentialsGenerateIn struct {
	WsId      string
	Name      string
	ExpiredAt int64
}

func (in *WorkspaceCredentialsGenerateIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ws_id", in.WsId, entities.IdNsWs),
		validator.StringRequired("name", in.Name),
		validator.NumberGreaterThanOrEqual("expired_at", in.ExpiredAt, 0),
	)
}

type WorkspaceCredentialsGenerateOut struct {
	Credentials *entities.WorkspaceCredentials
	Password    string
}

func (uc *workspaceCredentials) Generate(ctx context.Context, in *WorkspaceCredentialsGenerateIn) (*WorkspaceCredentialsGenerateOut, error) {
	now := uc.infra.Timer.Now()
	doc := &entities.WorkspaceCredentials{
		WsId:      in.WsId,
		Name:      in.Name,
		ExpiredAt: in.ExpiredAt,
	}
	doc.Id = identifier.New(entities.IdNsWsc)
	doc.SetAT(now)

	password := fmt.Sprintf("%s.%s", project.RegionCode(), utils.RandomString(constants.PasswordLength))
	// once we got error, reject entirely request instead of do a partial success request
	hash, err := utils.PasswordHash(password)
	if err != nil {
		return nil, err
	}
	doc.Hash = hash

	credentials, err := uc.repositories.Database().WorkspaceCredentials().Create(ctx, doc)
	if err != nil {
		return nil, err
	}

	res := &WorkspaceCredentialsGenerateOut{Credentials: credentials, Password: password}
	return res, nil
}
