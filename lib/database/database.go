package database

import (
	"context"
	"sync"
	"time"

	"github.com/cbstorm/wyrstream/lib/configs"
	"github.com/cbstorm/wyrstream/lib/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

var instance *Database
var instance_sync sync.Once

type Database struct {
	client *mongo.Client
	db     *mongo.Database
	logger *logger.Logger
	config *configs.Config
}

func GetDatabase() *Database {
	if instance == nil {
		instance_sync.Do(func() {
			instance = &Database{
				logger: logger.NewLogger("DATABASE"),
				config: configs.GetConfig(),
			}
		})
	}
	return instance
}

func (d *Database) Client() *mongo.Client {
	return d.client
}

func (d *Database) DB() *mongo.Database {
	return d.db
}

func (d *Database) Connect() error {
	clientOptions := options.Client().ApplyURI(d.config.MONGODB_URL)
	clientOptions.SetMaxPoolSize(30)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		return err
	}
	d.client = client
	d.db = client.Database(d.config.MONGODB_DB_NAME, options.Database().SetWriteConcern(writeconcern.Majority()), options.Database().SetReadPreference(readpref.Secondary()))
	if err := d.db.Client().Ping(context.Background(), options.Client().ReadPreference); err != nil {
		return err
	}
	d.logger.Info("Connected to MongoDB successfully.")
	d.CreateIndexes()
	return nil
}

func (d *Database) Close() error {
	return d.client.Disconnect(context.Background())
}

func (d *Database) CreateIndex(i *Index) error {
	coll := d.db.Collection(i.CollectionName)
	opts := options.Index().SetUnique(i.Uniq).SetName(i.Name)
	if i.Weight != nil {
		opts.SetWeights(bson.M(i.Weight))
	}
	indexModel := mongo.IndexModel{
		Keys:    i.Fields,
		Options: opts,
	}
	indexName, err := coll.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		d.logger.Error("Could not create [%s] on collection [%s] with error: %v", indexName, i.CollectionName, err)
		return err
	}
	d.logger.Info("Created index [%s] [%s] successfully", i.CollectionName, indexName)
	return nil
}
