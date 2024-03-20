package elasticex

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/xm-chentl/goresource"
)

type repository struct {
	ctx    context.Context
	client *elasticsearch.Client
}

func (r *repository) Create(entry goresource.IDbModel, args ...interface{}) (err error) {
	entryByte, err := json.Marshal(entry)
	if err != nil {
		return
	}

	req := esapi.IndicesCreateRequest{
		Index: entry.Table(),
		Body:  bytes.NewReader(entryByte),
	}
	resp, err := req.Do(r.ctx, r.client)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.IsError() {
		err = fmt.Errorf("creating the index(%s) err: %v", entry.Table(), err)
		return
	}

	var res map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		err = fmt.Errorf("creating parsing the response body: %v", err)
	}

	return
}

// delete args 0 -> 支持many
func (r *repository) Delete(entry goresource.IDbModel, args ...interface{}) (err error) {

	return
}

// Update args 0 upset 1 filter
func (r *repository) Update(entry goresource.IDbModel, args ...interface{}) (err error) {

	return
}

func (r *repository) Query() goresource.IQuery {
	return nil
}
