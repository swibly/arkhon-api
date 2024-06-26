package usecase

import (
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/service/repository"
)

type PermissionUseCase struct {
	ur repository.PermissionRepository
}

func NewPermissionUseCase() PermissionUseCase {
	return PermissionUseCase{ur: repository.NewPermissionRepository()}
}

func (p PermissionUseCase) GetPermissions(userID uint) ([]*model.Permission, error) {
	following, err := p.ur.GetPermissions(userID)
	return following, err
}
