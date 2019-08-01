package storage_service

import (
	"bytes"
	"context"
	"encoding/csv"
	"github.com/nettyrnp/url-shortener/config"
	"github.com/nettyrnp/url-shortener/util"
	"github.com/pkg/errors"
	"io"
)

const (
	dbfile = "storage/db.csv"
)

var (
	// List of allowed values of 'action' (in route '/v1/{action}/...') together with their short versions
	FullToShortMap map[string]string
	// Reversed map
	ShortToFullMap map[string]string
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
	FullToShortMap = dbFromCsv("FULL", "SHORT", dbfile)
	ShortToFullMap = dbFromCsv("SHORT", "FULL", dbfile)
	return nil
}

func dbFromCsv(id1, id2, filename string) map[string]string {
	m, err := toMap(id1, filename)
	util.Die(err)
	m2 := map[string]string{}
	for k, v := range m {
		m2[k] = v[id2]
	}
	return m2
}

//func (dbs *storageService) GetRow(ctx context.Context, id string, version string) (*api.Row, error) {
//	var row *api.Row
//	data, err := dbs.retrieveData(ctx, getRowQuery, id, version, row)
//	if err != nil {
//		return nil, errors.Wrap(err, "retrieve row from database")
//	}
//	if data == nil {
//		return nil, nil
//	}
//	return data.(*api.Row), nil
//}

func toMap(idField, fname string) (map[string]map[string]string, error) {
	text, err := util.ReadFile(fname)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBufferString(text)
	arr := CSVToMap(buf)

	m, err := arrToMap(idField, arr)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func arrToMap(idField string, arr []map[string]string) (map[string]map[string]string, error) {
	m := map[string]map[string]string{}
	for _, row := range arr {
		id, ok := row[idField]
		if ok {
			m[id] = row
		} else {
			return nil, errors.Errorf("Now column '%v'", id)
		}
	}
	return m, nil
}

func CSVToMap(reader io.Reader) []map[string]string {
	r := csv.NewReader(reader)
	r.Comma = ';'
	rows := []map[string]string{}
	var header []string
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		util.Die(err)
		if header == nil {
			header = record
		} else {
			dict := map[string]string{}
			for i := range header {
				dict[header[i]] = record[i]
			}
			rows = append(rows, dict)
		}
	}
	return rows
}
