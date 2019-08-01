package storage_service

import (
	"context"
	"github.com/nettyrnp/url-shortener/config"
	"github.com/pkg/errors"
)

const (
	getQuery = `SELECT xxx FROM yyy WHERE id = $1 AND version = $2`
)

type StorageService struct {
	config.StorageConfig
}

func NewDBStorageService(ctx context.Context, storageConfig config.StorageConfig) (*StorageService, error) {
	dbs := &StorageService{
		StorageConfig: storageConfig,
	}
	if err := dbs.initializeDatabase(ctx); err != nil {
		return nil, errors.Wrap(err, "initialize DB http.")
	}
	return dbs, nil
}

func (dbs *StorageService) initializeDatabase(ctx context.Context) error {
	return nil
}

//func (dbs *storageService) GetRow(ctx context.Context, datasetId string, version string) (*api.Row, error) {
//	var row *api.Row
//	data, err := dbs.retrieveData(ctx, getRowQuery, datasetId, version, row)
//	if err != nil {
//		return nil, errors.Wrap(err, "retrieve materialization row from database")
//	}
//	if data == nil {
//		return nil, nil
//	}
//	return data.(*api.Row), nil
//}
