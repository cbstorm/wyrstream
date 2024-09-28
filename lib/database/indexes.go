package database

import "go.mongodb.org/mongo-driver/bson"

type Index struct {
	CollectionName string
	Name           string
	Fields         bson.D
	Uniq           bool
	Weight         map[string]interface{}
}

func (d *Database) CreateIndexes() error {
	if err := d.CreateIndex(&Index{
		CollectionName: "admins",
		Name:           "email_1",
		Fields:         bson.D{{Key: "email", Value: 1}},
		Uniq:           true,
	}); err != nil {
		return err
	}
	if err := d.CreateIndex(&Index{
		CollectionName: "users",
		Name:           "email_1",
		Fields:         bson.D{{Key: "email", Value: 1}},
		Uniq:           true,
	}); err != nil {
		return err
	}
	if err := d.CreateIndex(&Index{
		CollectionName: "users",
		Name:           "email_name_text",
		Fields:         bson.D{{Key: "email", Value: "text"}, {Key: "name", Value: "text"}},
	}); err != nil {
		return err
	}
	if err := d.CreateIndex(&Index{
		CollectionName: "streams",
		Name:           "stream_id_1",
		Fields:         bson.D{{Key: "stream_id", Value: 1}},
		Uniq:           true,
	}); err != nil {
		return err
	}
	if err := d.CreateIndex(&Index{
		CollectionName: "streams",
		Name:           "title_desc_text",
		Fields:         bson.D{{Key: "title", Value: "text"}, {Key: "description", Value: "text"}},
	}); err != nil {
		return err
	}
	if err := d.CreateIndex(&Index{
		CollectionName: "stream_logs",
		Name:           "stream_obj_id_1",
		Fields:         bson.D{{Key: "stream_obj_id", Value: 1}},
	}); err != nil {
		return err
	}
	if err := d.CreateIndex(&Index{
		CollectionName: "vods",
		Name:           "owner_id_1",
		Fields:         bson.D{{Key: "owner_id", Value: 1}},
	}); err != nil {
		return err
	}
	if err := d.CreateIndex(&Index{
		CollectionName: "vods",
		Name:           "title_desc_text",
		Fields:         bson.D{{Key: "title", Value: "text"}, {Key: "description", Value: "text"}},
	}); err != nil {
		return err
	}
	return nil
}
