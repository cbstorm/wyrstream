package repositories

import (
	"context"
	"time"

	"github.com/cbstorm/wyrstream/lib/database"
	"github.com/cbstorm/wyrstream/lib/dtos"
	"github.com/cbstorm/wyrstream/lib/entities"
	"github.com/cbstorm/wyrstream/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type CRUDOption struct {
	ctx             context.Context
	bulkWriteOrder  bool
	bulkWriteUpsert bool
	sort            map[string]interface{}
}

func _NewCRUDOption() *CRUDOption {
	return &CRUDOption{}
}

type CURDOptionFunc func(*CRUDOption)

func WithContext(ctx context.Context) CURDOptionFunc {
	return func(c *CRUDOption) {
		c.ctx = ctx
	}
}

func WithBulkWriteOrder(bulkWriteOrder bool) CURDOptionFunc {
	return func(c *CRUDOption) {
		c.bulkWriteOrder = bulkWriteOrder
	}
}

func WithBulkWriteUpsert(bulkWriteUpsert bool) CURDOptionFunc {
	return func(c *CRUDOption) {
		c.bulkWriteUpsert = bulkWriteUpsert
	}
}

func WithSort(sort map[string]interface{}) CURDOptionFunc {
	return func(c *CRUDOption) {
		c.sort = sort
	}
}

type CRUDRepository[T entities.IEntity] struct {
	collection *mongo.Collection
}

func (r *CRUDRepository[T]) FindOneById(id primitive.ObjectID, out T, opts ...CURDOptionFunc) (error, bool) {
	return r.FindOne(map[string]interface{}{
		"_id": id,
	}, out, opts...)

}

func (r *CRUDRepository[T]) FindOne(filter map[string]interface{}, out T, opts ...CURDOptionFunc) (error, bool) {
	ctx := context.Background()
	find_one_opts := options.FindOne()
	o := _NewCRUDOption()
	for _, v := range opts {
		v(o)
	}
	if o.ctx != nil {
		ctx = o.ctx
	}
	if len(o.sort) != 0 {
		find_one_opts.SetSort(o.sort)
	}
	err := r.collection.FindOne(ctx, bson.M(filter), find_one_opts).Decode(out)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, true
		}
	}
	return err, false
}

func (r *CRUDRepository[T]) Find(filter map[string]interface{}, results *[]T, opts ...CURDOptionFunc) error {
	ctx := context.Background()
	find_opts := options.Find()
	o := _NewCRUDOption()
	for _, v := range opts {
		v(o)
	}
	if o.ctx != nil {
		ctx = o.ctx
	}
	if len(o.sort) != 0 {
		find_opts.SetSort(o.sort)
	}
	cursor, err := r.collection.Find(ctx, bson.M(filter), find_opts)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	return cursor.All(ctx, results)
}
func (r *CRUDRepository[T]) FindManyByIds(ids []primitive.ObjectID, results *[]T, opts ...CURDOptionFunc) error {
	filter := map[string]interface{}{
		"_id": map[string]interface{}{
			"$in": ids,
		},
	}
	err := r.Find(filter, results, opts...)
	if err != nil {
		return err
	}
	return nil
}

func (r *CRUDRepository[T]) InsertOne(doc T, opts ...CURDOptionFunc) error {
	ctx := context.Background()
	o := _NewCRUDOption()
	for _, v := range opts {
		v(o)
	}
	if o.ctx != nil {
		ctx = o.ctx
	}
	_, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return err
	}
	return nil
}

func (r *CRUDRepository[T]) InsertMany(docs []T, opts ...CURDOptionFunc) ([]interface{}, error) {
	ctx := context.Background()
	o := _NewCRUDOption()
	for _, v := range opts {
		v(o)
	}
	if o.ctx != nil {
		ctx = o.ctx
	}
	interface_slice := make([]interface{}, len(docs))
	for i, v := range docs {
		interface_slice[i] = v
	}
	result, err := r.collection.InsertMany(ctx, interface_slice)
	if err != nil {
		return nil, err
	}
	return result.InsertedIDs, nil
}

func (r *CRUDRepository[T]) UpdateOne(filter map[string]interface{}, update interface{}, out interface{}, opts ...CURDOptionFunc) error {
	ctx := context.Background()
	o := _NewCRUDOption()
	for _, v := range opts {
		v(o)
	}
	if o.ctx != nil {
		ctx = o.ctx
	}
	update_opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := r.collection.FindOneAndUpdate(ctx, bson.M(filter), bson.M{"$set": update}, update_opts)
	if result.Err() != nil {
		return result.Err()
	}
	err := result.Decode(out)
	return err
}

func (r *CRUDRepository[T]) UpdateOneById(id primitive.ObjectID, update interface{}, out interface{}, opts ...CURDOptionFunc) error {
	err := r.UpdateOne(map[string]interface{}{
		"_id": id,
	}, update, out, opts...)
	if err != nil {
		return err
	}
	return nil
}

func (r *CRUDRepository[T]) Update(filter map[string]interface{}, update interface{}, opts ...CURDOptionFunc) error {
	ctx := context.Background()
	o := _NewCRUDOption()
	for _, v := range opts {
		v(o)
	}
	if o.ctx != nil {
		ctx = o.ctx
	}
	_, err := r.collection.UpdateMany(ctx, bson.M(filter), bson.M{"$set": update})

	return err
}

func (r *CRUDRepository[T]) Upsert(filter map[string]interface{}, update interface{}, out interface{}, opts ...CURDOptionFunc) error {
	ctx := context.Background()
	o := _NewCRUDOption()
	for _, v := range opts {
		v(o)
	}
	if o.ctx != nil {
		ctx = o.ctx
	}
	upsert_opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	result := r.collection.FindOneAndUpdate(ctx, bson.M(filter), bson.M{"$set": update}, upsert_opts)
	if err := result.Err(); err != nil && err != mongo.ErrNoDocuments {
		return err
	}
	if err := result.Decode(out); err != nil && err != mongo.ErrNoDocuments {
		return err
	}
	return nil
}

func (r *CRUDRepository[T]) Increase(filter map[string]interface{}, update interface{}, out T, opts ...CURDOptionFunc) error {
	ctx := context.Background()
	o := _NewCRUDOption()
	for _, v := range opts {
		v(o)
	}
	if o.ctx != nil {
		ctx = o.ctx
	}
	update_opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := r.collection.FindOneAndUpdate(ctx, bson.M(filter), bson.M{"$inc": update}, update_opts)
	if result.Err() != nil {
		return result.Err()
	}
	err := result.Decode(out)
	return err
}

type BulkUpdateRecord struct {
	Filter     map[string]interface{}
	UpdateData map[string]interface{}
}

func (r *CRUDRepository[T]) BulkUpdate(records []*BulkUpdateRecord, opts ...CURDOptionFunc) error {
	ctx := context.Background()
	ordered := false
	upsert := false
	o := _NewCRUDOption()
	for _, v := range opts {
		v(o)
	}
	if o.ctx != nil {
		ctx = o.ctx
	}
	if o.bulkWriteOrder {
		ordered = o.bulkWriteOrder
	}
	if o.bulkWriteUpsert {
		upsert = o.bulkWriteUpsert
	}
	bulk_update_models := make([]mongo.WriteModel, 0)
	for _, v := range records {
		bulk_update_models = append(bulk_update_models, mongo.NewUpdateOneModel().SetFilter(bson.M(v.Filter)).SetUpdate(bson.M{"$set": v.UpdateData}).SetUpsert(upsert))
	}

	bulk_update_opts := options.BulkWrite().SetOrdered(ordered)
	_, err := r.collection.BulkWrite(ctx, bulk_update_models, bulk_update_opts)
	return err
}

func (r *CRUDRepository[T]) DeleteOne(filter map[string]interface{}, out T, opts ...CURDOptionFunc) (error, bool) {
	ctx := context.Background()
	o := _NewCRUDOption()
	for _, v := range opts {
		v(o)
	}
	if o.ctx != nil {
		ctx = o.ctx
	}
	result := r.collection.FindOneAndDelete(ctx, bson.M(filter))
	err := result.Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, true
		}
		return err, false
	}
	err = result.Decode(out)
	if err != nil {
		return err, false
	}
	return nil, false
}

func (r *CRUDRepository[T]) DeleteOneById(id primitive.ObjectID, out T, opts ...CURDOptionFunc) (error, bool) {
	err, isNotFound := r.DeleteOne(map[string]interface{}{
		"_id": id,
	}, out, opts...)
	if isNotFound {
		return nil, true
	}
	if err != nil {
		return err, false
	}
	return nil, false
}

func (r *CRUDRepository[T]) Delete(filter map[string]interface{}, opts ...CURDOptionFunc) error {
	ctx := context.Background()
	o := _NewCRUDOption()
	for _, v := range opts {
		v(o)
	}
	if o.ctx != nil {
		ctx = o.ctx
	}
	_, err := r.collection.DeleteMany(ctx, bson.M(filter))
	return err
}

func (r *CRUDRepository[T]) Count(filter map[string]interface{}, opts ...CURDOptionFunc) (int64, error) {
	ctx := context.Background()
	o := _NewCRUDOption()
	for _, v := range opts {
		v(o)
	}
	if o.ctx != nil {
		ctx = o.ctx
	}
	return r.collection.CountDocuments(ctx, bson.M(filter))
}

func (r *CRUDRepository[T]) Fetch(fetchArgs *dtos.FetchArgs, out *[]T, opts ...CURDOptionFunc) (*dtos.FetchOutput[T], error) {
	ctx := context.Background()
	o := _NewCRUDOption()
	for _, v := range opts {
		v(o)
	}
	if o.ctx != nil {
		ctx = o.ctx
	}
	findOptions := options.Find()
	findOptions.SetSkip((fetchArgs.Page - 1) * (fetchArgs.Limit))
	findOptions.SetLimit(fetchArgs.Limit)
	var sort bson.D
	if len(fetchArgs.Order) > 0 {
		for k, v := range fetchArgs.Order {
			sort = append(sort, bson.E{Key: k, Value: v})
		}
	}
	if !fetchArgs.IsOrderByLocation() {
		findOptions.SetSort(sort)
	}
	filter := fetchArgs.Filter
	if fetchArgs.IsHavingLocation() && fetchArgs.IsOrderByLocation() {
		filter["location"] = map[string]interface{}{
			"$nearSphere": map[string]interface{}{
				"$geometry": map[string]interface{}{
					"type":        "Point",
					"coordinates": []float64{fetchArgs.GetLng(), fetchArgs.GetLat()},
				},
				"$maxDistance": 500000,
			},
		}
	}
	if fetchArgs.Search != "" {
		filter["$text"] = bson.D{{Key: "$search", Value: fetchArgs.Search}}
		findOptions.SetProjection(bson.D{{Key: "score", Value: bson.D{{Key: "$meta", Value: "textScore"}}}})
		findOptions.SetSort(bson.D{{Key: "score", Value: bson.D{{Key: "$meta", Value: "textScore"}}}})
	}
	count_filter := map[string]interface{}{}
	utils.CopyMap(&filter, &count_filter, []string{"location"})
	count, err := r.collection.CountDocuments(ctx, bson.M(count_filter))
	if err != nil {
		return nil, err
	}
	cursor, err := r.collection.Find(ctx, bson.M(filter), findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	err = cursor.All(ctx, out)
	if err != nil {
		return nil, err
	}
	fetchOut := dtos.NewFetchOutput(count, *out)
	return fetchOut, nil
}

func (r *CRUDRepository[T]) WithTransaction(process func(ctx mongo.SessionContext) error) error {
	client := database.GetDatabase().Client()
	wc := writeconcern.Majority()
	txnOptions := options.Transaction().SetWriteConcern(wc)
	session, err := client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(context.Background())

	if err := mongo.WithSession(context.Background(), session, func(ctx mongo.SessionContext) error {
		if err := session.StartTransaction(txnOptions); err != nil {
			return err
		}
		if err := process(ctx); err != nil {
			return err
		}

		if err := session.CommitTransaction(ctx); err != nil {
			return err
		}
		return nil
	}); err != nil {
		if err := session.AbortTransaction(context.Background()); err != nil {
			return err
		}
		return err
	}
	return nil
}

type EAggregationStageTypes string
type EAggregationStageOpTypes string

const (
	MATCH_STAGE  EAggregationStageTypes = "$match"
	GROUP_STAGE  EAggregationStageTypes = "$group"
	LOOKUP_STAGE EAggregationStageTypes = "$lookup"
	SORT_STAGE   EAggregationStageTypes = "$sort"
	LIMIT_STAGE  EAggregationStageTypes = "$limit"
	UNWIND_STAGE EAggregationStageTypes = "$unwind"
	UNSET_STAGE  EAggregationStageTypes = "$unset"
)
const (
	OP_TYPE_M EAggregationStageOpTypes = "M"
	OP_TYPE_A EAggregationStageOpTypes = "A"
	OP_TYPE_S EAggregationStageOpTypes = "S"
	OP_TYPE_N EAggregationStageOpTypes = "N"
)

type AggregationStage struct {
	stage  EAggregationStageTypes
	opM    map[string]interface{}
	opA    []interface{}
	opS    string
	opN    int
	opType EAggregationStageOpTypes
}

func MatchStage(op map[string]interface{}) *AggregationStage {
	return &AggregationStage{
		stage:  MATCH_STAGE,
		opType: OP_TYPE_M,
		opM:    op,
	}
}
func GroupStage(op map[string]interface{}) *AggregationStage {
	return &AggregationStage{
		stage:  GROUP_STAGE,
		opType: OP_TYPE_M,
		opM:    op,
	}
}

type LookupStageOp struct {
	From         string
	LocalField   string
	ForeignField string
	As           string
}

func LookupStage(op *LookupStageOp) *AggregationStage {
	return &AggregationStage{
		stage:  LOOKUP_STAGE,
		opType: OP_TYPE_M,
		opM: map[string]interface{}{
			"from":         op.From,
			"localField":   op.LocalField,
			"foreignField": op.ForeignField,
			"as":           op.As,
		},
	}
}

func SortStage(op map[string]interface{}) *AggregationStage {
	return &AggregationStage{
		stage:  SORT_STAGE,
		opType: OP_TYPE_M,
		opM:    op,
	}
}

func LimitStage(op int) *AggregationStage {
	return &AggregationStage{
		stage:  LIMIT_STAGE,
		opType: OP_TYPE_N,
		opN:    op,
	}
}
func UnwindStage(op string) *AggregationStage {
	return &AggregationStage{
		stage:  UNWIND_STAGE,
		opType: OP_TYPE_S,
		opS:    op,
	}
}

func UnsetStage(op []interface{}) *AggregationStage {
	return &AggregationStage{
		stage: UNSET_STAGE,
		opA:   op,
	}
}

func (r *CRUDRepository[T]) Aggregate(stages []*AggregationStage, out interface{}, opts ...CURDOptionFunc) error {
	pipeline := mongo.Pipeline{}
	for _, v := range stages {
		if v.opType == OP_TYPE_M {
			pipeline = append(pipeline, bson.D{{Key: string(v.stage), Value: bson.M(v.opM)}})
			continue
		}
		if v.opType == OP_TYPE_A {
			pipeline = append(pipeline, bson.D{{Key: string(v.stage), Value: bson.A(v.opA)}})
			continue
		}
		if v.opType == OP_TYPE_S {
			pipeline = append(pipeline, bson.D{{Key: string(v.stage), Value: v.opS}})
			continue
		}
		if v.opType == OP_TYPE_N {
			pipeline = append(pipeline, bson.D{{Key: string(v.stage), Value: v.opN}})
			continue
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	o := _NewCRUDOption()
	for _, v := range opts {
		v(o)
	}
	if o.ctx != nil {
		cancel()
		ctx = o.ctx
	}
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}
	if err = cursor.All(ctx, out); err != nil {
		return err
	}
	return nil
}
