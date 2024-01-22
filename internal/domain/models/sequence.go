package models

type Sequence struct {
	CollectionName string `bson:"collection_name"`
	Counter        int64  `bson:"counter"`
}
