package api

//import (
//	"bytes"
//	"encoding/json"
//	"fmt"
//	"time"
//
//	"github.com/pkg/errors"
//)
//
//type Row struct {
//	// DataSetID identifies a data set.
//	DataSetID string `json:"dataSetId"`
//	// Version identifies a unique version of a data set.
//	Version string `json:"version"`
//	// Comment is a comment about a version of a data set.
//	Comment string `json:"comment"`
//	Caches []CacheRow `json:"caches"`
//	// CreatedAt is a creation timestamp.
//	CreatedAt time.Time
//	// UpdatedAt is a last-update timestamp.
//	UpdatedAt time.Time
//	// DeletedAt is soft-deletion timestamp. A nil value means it isn't deleted.
//	DeletedAt *time.Time
//}
//
//// GetCahceRow returns the cache row and true for an existing ID; otherwise false.
//func (row Row) GetCache(id string) (CacheRow, bool) {
//	for _, cacheRow := range row.Caches {
//		if cacheRow.ID == id {
//			return cacheRow, true
//		}
//	}
//	return CacheRow{}, false
//}
