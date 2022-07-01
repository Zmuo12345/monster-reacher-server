package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongodb struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoDBDriver(uri string, db string, collection string) (DBDriver, error) {
	driver := &mongodb{}
	ctx, cancle := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancle()
	var err error = nil
	driver.client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	driver.collection = driver.client.Database(db).Collection(collection)
	return driver, nil
}

func (db *mongodb) SelectOne(ctx context.Context, filter interface{}) interface{} {
	result := db.collection.FindOne(ctx, filter)
	return result
}
func (db *mongodb) PushOne(ctx context.Context, data interface{}) (interface{}, error) {
	result, err := db.collection.InsertOne(ctx, data)
	return result, err
}
func (db *mongodb) DeleteOne(ctx context.Context, filter interface{}) error {
	_, err := db.collection.DeleteOne(ctx, filter)
	return err
}
func (db *mongodb) UpdateOne(ctx context.Context, filter interface{}, data interface{}) error {
	update := bson.D{primitive.E{Key: "$set", Value: data}}
	_, err := db.collection.UpdateOne(ctx, filter, update)
	return err
}
func (db *mongodb) Close() error {
	return db.client.Disconnect(context.Background())
}

func MongoDBSelectOneQueryFilterOne(key string, value interface{}) interface{} {
	return bson.D{primitive.E{Key: key, Value: value}}
}

func MongoDBSelectOneQueryFilterMany(keys []string, values []interface{}) interface{} {
	filter := bson.D{}
	for i := range keys {
		filter = append(filter, primitive.E{Key: keys[i], Value: values[i]})
	}
	return filter
}

func MongoDBSelectOneResultGetError(result interface{}) error {
	return result.(*mongo.SingleResult).Err()
}

func MongoDBDecodeResultToStruct(result interface{}, output interface{}) error {
	if _, ok := result.(*mongo.SingleResult); ok {
		return result.(*mongo.SingleResult).Decode(output)
	}
	return nil
}

func MongoDBDecodeResultToID(result interface{}) string {
	if _, ok := result.(*mongo.InsertOneResult); ok {
		return result.(*mongo.InsertOneResult).InsertedID.(primitive.ObjectID).Hex()
	}
	return ""
}
