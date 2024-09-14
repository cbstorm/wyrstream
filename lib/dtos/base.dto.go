package dtos

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cbstorm/wyrstream/lib/exceptions"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FetchArgs struct {
	page     int64
	limit    int64
	order    map[string]interface{}
	filter   map[string]interface{}
	search   string
	includes map[string]int
	fields   map[string]int
	location []float64
}

func NewFetchArgs() *FetchArgs {
	return &FetchArgs{
		page:  1,
		limit: 10,
		order: map[string]interface{}{
			"_id": -1,
		},
		filter:   map[string]interface{}{},
		includes: map[string]int{},
		fields:   map[string]int{},
		location: []float64{0, 0},
	}
}

func (fetchArgs *FetchArgs) ParseQueries(input map[string]string) *FetchArgs {
	page, err := strconv.Atoi(input["page"])
	if err != nil || page < 0 {
		page = 1
	} else {
		fetchArgs.page = int64(page)
	}
	limit, err := strconv.Atoi(input["limit"])
	if err != nil || limit < 0 {
		limit = 10
	} else {
		fetchArgs.limit = int64(limit)
	}
	includes_str := input["includes"]
	if includes_str != "" {
		includes := strings.Split(includes_str, ",")
		for _, v := range includes {
			fetchArgs.includes[v] = 1
		}
	}
	filter_str := input["filter"]
	if filter_str != "" {
		filter := strings.Split(filter_str, ",")
		for _, v := range filter {
			e := strings.Split(v, ":")
			obj_id, err := primitive.ObjectIDFromHex(e[1])
			if err == nil {
				fetchArgs.filter[e[0]] = obj_id
				continue
			}
			if e[1] == "true" || e[1] == "false" {
				if e[1] == "true" {
					fetchArgs.filter[e[0]] = true
				}
				if e[1] == "false" {
					fetchArgs.filter[e[0]] = false
				}
				continue
			}
			fetchArgs.filter[e[0]] = e[1]
		}
	}

	ignore_str := input["ignore"]
	if ignore_str != "" {
		ignores := strings.Split(ignore_str, ",")
		for _, v := range ignores {
			e := strings.Split(v, ":")
			obj_id, err := primitive.ObjectIDFromHex(e[1])
			if err == nil {
				fetchArgs.filter[e[0]] = map[string]interface{}{
					"$ne": obj_id,
				}
				continue
			}
			fetchArgs.filter[e[0]] = map[string]interface{}{
				"$ne": e[1],
			}
		}
	}

	sort_str := input["sort"]
	if sort_str != "" {
		fetchArgs.order = map[string]interface{}{}
		sort := strings.Split(sort_str, ",")
		for _, v := range sort {
			e := strings.Split(v, ":")
			e_val, err := strconv.Atoi(e[1])
			if err == nil {
				fetchArgs.order[e[0]] = e_val
			}
		}
	}
	location_str := input["location"]
	if location_str != "" {
		lat_lng := strings.Split(location_str, ",")
		lat, lat_err := strconv.ParseFloat(lat_lng[0], 64)
		lng, lng_err := strconv.ParseFloat(lat_lng[1], 64)
		if lat_err == nil && lng_err == nil {
			fetchArgs.location = []float64{lat, lng}
		}
	}
	fetchArgs.search = input["search"]
	return fetchArgs
}

func (fetchArgs *FetchArgs) SetFilter(key string, value interface{}) *FetchArgs {
	fetchArgs.filter[key] = value
	return fetchArgs
}

func (fetchArgs *FetchArgs) SetOrder(key string, value interface{}) *FetchArgs {
	fetchArgs.order[key] = value
	return fetchArgs
}
func (fetchArgs *FetchArgs) SetLimit(limit int64) *FetchArgs {
	fetchArgs.limit = limit
	return fetchArgs
}

func (fetchArgs *FetchArgs) IsIncludes(key string) bool {
	return fetchArgs.includes[key] == 1
}
func (fetchArgs *FetchArgs) IsHavingLocation() bool {
	return fetchArgs.location[0] != 0 && fetchArgs.location[1] != 0
}

func (fetchArgs *FetchArgs) IsOrderByLocation() bool {
	return fetchArgs.order["location"] == 1
}
func (fetchArgs *FetchArgs) GetLat() float64 {
	return fetchArgs.location[0]
}
func (fetchArgs *FetchArgs) GetLng() float64 {
	return fetchArgs.location[1]
}
func (fetchArgs *FetchArgs) Page() int64 {
	return fetchArgs.page
}
func (fetchArgs *FetchArgs) Limit() int64 {
	return fetchArgs.limit
}
func (fetchArgs *FetchArgs) Order() map[string]interface{} {
	return fetchArgs.order
}
func (fetchArgs *FetchArgs) Filter() map[string]interface{} {
	return fetchArgs.filter
}
func (fetchArgs *FetchArgs) Search() string {
	return fetchArgs.search
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
