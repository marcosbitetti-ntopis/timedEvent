package collection_managment

import (
	"context"
	"errors"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/ivanmeca/timedEvent/application/modules/database"
	"github.com/ivanmeca/timedEvent/application/modules/database/data_types"
)

type Collection struct {
	Db               driver.Database
	Collection       string
	CollectionDriver driver.Collection
}

func (coll *Collection) DeleteItem(keyList []string) ([]data_types.ArangoCloudEvent, error) {
	var oldDocs []data_types.ArangoCloudEvent
	ctx := driver.WithReturnOld(context.Background(), oldDocs)
	for _, key := range keyList {
		_, err := coll.CollectionDriver.RemoveDocument(ctx, key)
		if err != nil {
			return nil, err
		}
	}
	return oldDocs, nil
}

func (coll *Collection) Insert(item *data_types.ArangoCloudEvent) (*data_types.ArangoCloudEvent, error) {
	var newDoc data_types.ArangoCloudEvent
	ctx := driver.WithReturnNew(context.Background(), newDoc)
	_, err := coll.CollectionDriver.CreateDocument(ctx, item)
	if err != nil {
		return nil, err
	}
	return &newDoc, nil
}

func (coll *Collection) Update(patch map[string]interface{}, key string) (*data_types.ArangoCloudEvent, error) {
	var newDoc data_types.ArangoCloudEvent
	ctx := driver.WithReturnNew(context.Background(), newDoc)
	_, err := coll.CollectionDriver.UpdateDocument(ctx, key, patch)
	if err != nil {
		return nil, err
	}
	return &newDoc, nil
}

func (coll *Collection) Read(filters []database.AQLComparator) ([]data_types.ArangoCloudEvent, error) {
	var item data_types.ArangoCloudEvent
	var list []data_types.ArangoCloudEvent

	bindVars := map[string]interface{}{}
	query := fmt.Sprintf("FOR item IN %s ", coll.Collection)
	glueStr := "FILTER"
	bindVarsNames := 0
	for _, filter := range filters {
		bindVars[string('A'+bindVarsNames)] = filter.Value
		query += fmt.Sprintf(" %s item.%s %s @%s", glueStr, filter.Field, filter.Comparator, string('A'+bindVarsNames))
		glueStr = "AND"
		bindVarsNames++
	}
	query += fmt.Sprintf(" SORT item.Context.time DESC RETURN item")
	cursor, err := coll.Db.Query(nil, query, bindVars)
	defer cursor.Close()
	if err != nil {
		return nil, errors.New("internal error: " + err.Error())
	}
	for cursor.HasMore() == true {
		_, err = cursor.ReadDocument(nil, &item)
		if err != nil {
			return nil, errors.New("internal error: " + err.Error())
		}
		list = append(list, item)
	}
	return list, nil
}

func (coll *Collection) ReadItem(key string) (*data_types.ArangoCloudEvent, error) {
	var item data_types.ArangoCloudEvent
	_, err := coll.CollectionDriver.ReadDocument(nil, key, &item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}
