package mongodb

import (
	"context"
	"log"
	"time"

	"github.com/Chetan177/pismo/pkg/db/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoClient db interface implementation for mongodb
type MongoClient struct {
	client                *mongo.Client
	accountCollection     *mongo.Collection
	transactionCollection *mongo.Collection
}

// NewMongoClient get a new instance of MongoClient and test the connection
func NewMongoClient(dbAddress, dbName, accCollection, transCollection string) *MongoClient {
	log.Println("connecting to client ", dbAddress)
	clientOptions := options.Client().ApplyURI(dbAddress)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB!")

	return &MongoClient{
		client:                client,
		accountCollection:     client.Database(dbName).Collection(accCollection),
		transactionCollection: client.Database(dbName).Collection(transCollection),
	}
}

// CreateAccount create a document entry in the accounts collection
func (m *MongoClient) CreateAccount(account *model.Account) (string, error) {
	account.Id = primitive.NewObjectID()
	account.TimeStamp = time.Now().String()
	_, err := m.accountCollection.InsertOne(context.Background(), account)
	if err != nil {
		return "", err
	}
	return account.Id.Hex(), nil

}

// GetAccount get a document from accounts collection for a specific id
func (m *MongoClient) GetAccount(id string) (*model.Account, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objectID}
	account := &model.Account{}
	err = m.accountCollection.FindOne(context.Background(), filter).Decode(&account)
	if err != nil {
		return nil, err
	}
	return account, nil

}

// IsAccountExsists check if an account exists or not
func (m *MongoClient) IsAccountExsists(id string) bool {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false
	}
	filter := bson.M{"_id": objectID}
	account := &model.Account{}
	err = m.accountCollection.FindOne(context.Background(), filter).Decode(&account)
	if err != nil {
		return false
	}
	return true

}

// CreateTransaction  create a document entry in the transaction collection
func (m *MongoClient) CreateTransaction(transaction *model.Transaction) (string, error) {
	transaction.Id = primitive.NewObjectID()
	transaction.TimeStamp = time.Now().String()
	_, err := m.transactionCollection.InsertOne(context.Background(), transaction)
	if err != nil {
		return "", err
	}
	return transaction.Id.Hex(), nil
}

// ShutDown stop disconnect client connection
func (m *MongoClient) ShutDown() {
	m.client.Disconnect(context.Background())
}
