/*
 * SPDX-License-Identifier: Apache-2.0
 *
 * The OpenSearch Contributors require contributions made to
 * this file be licensed under the Apache-2.0 license or a
 * compatible open source license.
 *
 * Modifications Copyright OpenSearch Contributors. See
 * GitHub history for details.
 */

package platform

import (
	"context"
	"encoding/json"
	"opensearch-cli/entity/platform"
	osg "opensearch-cli/gateway/platform"
	mapper "opensearch-cli/mapper/platform"

	"fmt"
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen  -destination=mocks/mock_platform.go -package=mocks . Controller

//Controller is an interface for OpenSearch
type Controller interface {
	GetDistinctValues(ctx context.Context, index string, field string) ([]interface{}, error)
	Curl(ctx context.Context, param platform.CurlCommandRequest) ([]byte, error)
}

type controller struct {
	gateway osg.Gateway
}

//New returns new instance of Controller
func New(gateway osg.Gateway) Controller {
	return &controller{
		gateway,
	}
}

//GetDistinctValues get only unique values for given index, given field name
func (c controller) GetDistinctValues(ctx context.Context, index string, field string) ([]interface{}, error) {
	if len(index) == 0 || len(field) == 0 {
		return nil, fmt.Errorf("index and field cannot be empty")
	}
	response, err := c.gateway.SearchDistinctValues(ctx, index, field)
	if err != nil {
		return nil, err
	}
	var data platform.Response
	err = json.Unmarshal(response, &data)
	if err != nil {
		return nil, err
	}

	var values []interface{}
	for _, bucket := range data.Aggregations.Items.Buckets {
		values = append(values, bucket.Key)
	}
	return values, nil
}

//Curl accept user request and convert to format which OpenSearch can understand
func (c controller) Curl(ctx context.Context, param platform.CurlCommandRequest) ([]byte, error) {
	curlRequest, err := mapper.CommandToCurlRequestParameter(param)
	if err != nil {
		return nil, err
	}
	return c.gateway.Curl(ctx, curlRequest)
}
