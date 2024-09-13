package dtos

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cbstorm/wyrstream/lib/entities"
	"github.com/cbstorm/wyrstream/lib/exceptions"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FetchArgs struct {
	Page     int64                  `json:"page,omitempty"`
	Limit    int64                  `json:"limit,omitempty"`
	Order    map[string]interface{} `json:"order,omitempty"`
	Filter   map[string]interface{} `json:"filter"`
	Search   string                 `json:"search,omitempty"`
	Includes map[string]int         `json:"includes"`
	Fields   map[string]int         `json:"fields"`
	Location []float64              `json:"location"`
}

func NewFetchArgs() *FetchArgs {
	return &FetchArgs{
		Page:  1,
		Limit: 10,
		Order: map[string]interface{}{
			"_id": -1,
		},
		Filter:   map[string]interface{}{},
		Includes: map[string]int{},
		Fields:   map[string]int{},
		Location: []float64{0, 0},
	}
}

func (fetchArgs *FetchArgs) ParseQueries(input map[string]string) *FetchArgs {
	page, err := strconv.Atoi(input["page"])
	if err != nil || page < 0 {
		page = 1
	} else {
		fetchArgs.Page = int64(page)
	}
	limit, err := strconv.Atoi(input["limit"])
	if err != nil || limit < 0 {
		limit = 10
	} else {
		fetchArgs.Limit = int64(limit)
	}
	includes_str := input["includes"]
	if includes_str != "" {
		includes := strings.Split(includes_str, ",")
		for _, v := range includes {
			fetchArgs.Includes[v] = 1
		}
	}
	filter_str := input["filter"]
	if filter_str != "" {
		filter := strings.Split(filter_str, ",")
		for _, v := range filter {
			e := strings.Split(v, ":")
			obj_id, err := primitive.ObjectIDFromHex(e[1])
			if err == nil {
				fetchArgs.Filter[e[0]] = obj_id
				continue
			}
			if e[1] == "true" || e[1] == "false" {
				if e[1] == "true" {
					fetchArgs.Filter[e[0]] = true
				}
				if e[1] == "false" {
					fetchArgs.Filter[e[0]] = false
				}
				continue
			}
			fetchArgs.Filter[e[0]] = e[1]
		}
	}

	ignore_str := input["ignore"]
	if ignore_str != "" {
		ignores := strings.Split(ignore_str, ",")
		for _, v := range ignores {
			e := strings.Split(v, ":")
			obj_id, err := primitive.ObjectIDFromHex(e[1])
			if err == nil {
				fetchArgs.Filter[e[0]] = map[string]interface{}{
					"$ne": obj_id,
				}
				continue
			}
			fetchArgs.Filter[e[0]] = map[string]interface{}{
				"$ne": e[1],
			}
		}
	}

	sort_str := input["sort"]
	if sort_str != "" {
		fetchArgs.Order = map[string]interface{}{}
		sort := strings.Split(sort_str, ",")
		for _, v := range sort {
			e := strings.Split(v, ":")
			e_val, err := strconv.Atoi(e[1])
			if err == nil {
				fetchArgs.Order[e[0]] = e_val
			}
		}
	}
	location_str := input["location"]
	if location_str != "" {
		lat_lng := strings.Split(location_str, ",")
		lat, lat_err := strconv.ParseFloat(lat_lng[0], 64)
		lng, lng_err := strconv.ParseFloat(lat_lng[1], 64)
		if lat_err == nil && lng_err == nil {
			fetchArgs.Location = []float64{lat, lng}
		}
	}
	fetchArgs.Search = input["search"]
	return fetchArgs
}

func (fetchArgs *FetchArgs) SetFilter(key string, value interface{}) *FetchArgs {
	fetchArgs.Filter[key] = value
	return fetchArgs
}

func (fetchArgs *FetchArgs) SetOrder(key string, value interface{}) *FetchArgs {
	fetchArgs.Order[key] = value
	return fetchArgs
}
func (fetchArgs *FetchArgs) SetLimit(limit int64) *FetchArgs {
	fetchArgs.Limit = limit
	return fetchArgs
}

func (fetchArgs *FetchArgs) IsIncludes(key string) bool {
	return fetchArgs.Includes[key] == 1
}
func (fetchArgs *FetchArgs) IsHavingLocation() bool {
	return fetchArgs.Location[0] != 0 && fetchArgs.Location[1] != 0
}

func (fetchArgs *FetchArgs) IsOrderByLocation() bool {
	return fetchArgs.Order["location"] == 1
}
func (fetchArgs *FetchArgs) GetLat() float64 {
	return fetchArgs.Location[0]
}
func (fetchArgs *FetchArgs) GetLng() float64 {
	return fetchArgs.Location[1]
}

type FetchOutput[T entities.IEntity] struct {
	Total  int64 `json:"total"`
	Result []T   `json:"result"`
}

func NewFetchOutput[T entities.IEntity](total int64, result []T) *FetchOutput[T] {
	return &FetchOutput[T]{
		Total:  total,
		Result: result,
	}
}

type LoadInput struct {
	Ids []string `json:"ids,omitempty"`
}
type LoadOutput[T any] struct {
	Data []*T `json:"data,omitempty"`
}

type GetOneInput struct {
	Id primitive.ObjectID `json:"_id,omitempty"`
}

func NewGetOneInput() *GetOneInput {
	return &GetOneInput{}
}

func (d *GetOneInput) SetId(id string) (*GetOneInput, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, exceptions.Err_BAD_REQUEST().SetMessage(fmt.Sprintf("id invalid: %s", id))
	}
	return &GetOneInput{Id: objId}, nil
}

type DeleteOneInput struct {
	Id primitive.ObjectID `json:"_id,omitempty"`
}

func (i *DeleteOneInput) SetId(id string) (*DeleteOneInput, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, exceptions.Err_BAD_REQUEST().SetMessage(fmt.Sprintf("id invalid: %s", id))
	}
	i.Id = objId
	return i, nil
}

func NewDeleteOneInput() *DeleteOneInput {
	return &DeleteOneInput{}
}

func ToObjectIds(ids []string) (*[]primitive.ObjectID, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	obj_ids := make([]primitive.ObjectID, len(ids))
	for i, v := range ids {
		obj_id, err := primitive.ObjectIDFromHex(v)
		if err != nil {
			return nil, exceptions.Err_BAD_REQUEST().SetMessage(fmt.Sprintf("id invalid with %s", v))
		}
		obj_ids[i] = obj_id
	}
	return &obj_ids, nil
}
