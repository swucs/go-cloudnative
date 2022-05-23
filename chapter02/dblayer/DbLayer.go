package dblayer

import "go-cloudnative/repository"

type DBTYPE string

const (
	MONGODB  DBTYPE = "mongodb"
	DYNAMODB DBTYPE = "dynamodb"
)

func NewPersistenceLayer(options DBTYPE, connection string) (repository.DatabaseHandler, error) {
	switch options {
	case MONGODB:
		return repository.NewMongoDBLayer(connection)
	}
	return nil, nil
}
