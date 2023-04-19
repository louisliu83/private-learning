package service

import (
	"context"

	"fedlearn/psi/api/types"
)

// DatasetService wrap functions related dataset
type DatasetService interface {
	GetDatasetCount(ctx context.Context, partyName, dsName string, index int32) (int64, error)
	GetDataset(ctx context.Context, partyName, dsName string) (*types.Dataset, error)
	GetDatasetList(ctx context.Context, partyName string) ([]types.Dataset, error)
}
