package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Base struct {
	collection *mongo.Collection
}

func (ins *Base) Add(ctx context.Context, obj interface{}) (id string, err error) {
	rslt, err := ins.collection.InsertOne(ctx, obj)
	if err != nil {
		return "", err
	}
	id = rslt.InsertedID.(bson.ObjectID).Hex()
	return id, nil
}
func (ins *Base) ById(ctx context.Context, id string) (interface{}, error) {
	var obj interface{}
	err := ins.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}
func (ins *Base) Get(ctx context.Context, filter interface{}, result interface{}) error {
	err := ins.collection.FindOne(ctx, filter).Decode(result)
	return err
}
func (ins *Base) GetAll(ctx context.Context, filter bson.M, results interface{}) error {
	cursor, err := ins.collection.Find(ctx, filter)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, results)
	if err != nil {
		return err
	}
	return nil
}
func (ins *Base) Update(ctx context.Context, filter interface{}, update interface{}) error {
	_, err := ins.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
func (ins *Base) Delete(ctx context.Context, filter interface{}) error {
	_, err := ins.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
func (ins *Base) Count(ctx context.Context, filter interface{}) (int64, error) {
	count, err := ins.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}
