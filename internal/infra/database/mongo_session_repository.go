package infra

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/beriloqueiroz/psi-mgnt/internal/application"
	domain "github.com/beriloqueiroz/psi-mgnt/internal/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoSessionRepository struct {
	Client            *mongo.Client
	PatientCollection *mongo.Collection
	SessionCollection *mongo.Collection
}

func NewMongoSessionRepository(
	ctx context.Context, uri string,
	patientCollection string, sessionCollection string,
	databaseName string) (*MongoSessionRepository, error) {
	client, err := connectToMongoDb(uri, ctx)
	if err != nil {
		return nil, err
	}
	database := client.Database(databaseName)
	initCollections(ctx, database, patientCollection, sessionCollection)
	PatientCollection := database.Collection(patientCollection)
	SessionCollection := database.Collection(sessionCollection)
	return &MongoSessionRepository{
		client,
		PatientCollection,
		SessionCollection,
	}, nil
}

func initCollections(ctx context.Context, database *mongo.Database, patientCollection string, sessionCollection string) error {
	result, err := database.ListCollectionNames(ctx, bson.D{{}})
	if err != nil {
		return err
	}
	if !slices.Contains(result, patientCollection) {
		err = database.CreateCollection(ctx, patientCollection)
		if err != nil {
			return err
		}
	}
	if !slices.Contains(result, sessionCollection) {
		err = database.CreateCollection(ctx, sessionCollection)
	}
	return err
}

func connectToMongoDb(uri string, ctx context.Context) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err

	}
	fmt.Println(" --------------- App has been Connected to MongoDB! --------------- ")
	return client, nil
}

func (mr *MongoSessionRepository) Create(ctx context.Context, session *domain.Session) error {
	fmt.Println(session)
	_, err := mr.SessionCollection.InsertOne(ctx, session)
	if err != nil {
		return err
	}
	return nil
}
func (mr *MongoSessionRepository) Delete(ctx context.Context, input application.DeleteRepositoryInput) error {
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"id", input.SessionId}},
				bson.D{{"owner_id", input.ProfessionalId}},
			}},
	}
	_, err := mr.SessionCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
func (mr *MongoSessionRepository) List(ctx context.Context, input application.ListByProfessionalRepositoryInput) ([]*domain.Session, error) {
	l := int64(input.PageSize)
	skip := int64(input.Page*input.PageSize - input.PageSize)
	findOptions := options.FindOptions{Limit: &l, Skip: &skip}

	var results []*domain.Session
	cur, err := mr.SessionCollection.Find(ctx, bson.D{{Key: "owner_id", Value: input.ProfessionalId}}, &findOptions)
	if err != nil {
		return nil, err
	}

	for cur.Next(ctx) {
		var elem domain.Session
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}

		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	cur.Close(ctx)
	return results, nil
}
func (mr *MongoSessionRepository) FindPatientByName(ctx context.Context, input application.FindPatientRepositoryInput) (*domain.Patient, error) {
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"name", input.Name}},
				bson.D{{"owner_id", input.ProfessionalId}},
			}},
	}
	var result domain.Patient
	err := mr.PatientCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}
func (mr *MongoSessionRepository) CreatePatient(ctx context.Context, patient *domain.Patient) error {
	_, err := mr.PatientCollection.InsertOne(ctx, patient)
	if err != nil {
		return err
	}
	return nil
}
func (mr *MongoSessionRepository) SearchPatientsByName(ctx context.Context, input application.SearchPatientsByNameRepositoryInput) ([]*domain.Patient, error) {
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{Key: "name", Value: primitive.Regex{
					Pattern: "/*" + input.Term + ".*",
					Options: "i",
				}}},
				bson.D{{"owner_id", input.ProfessionalId}},
			}},
	}
	l := int64(input.PageSize)
	skip := int64(input.Page*input.PageSize - input.PageSize)
	findOptions := options.FindOptions{Limit: &l, Skip: &skip}

	var results []*domain.Patient
	cur, err := mr.PatientCollection.Find(ctx, filter, &findOptions)
	if err != nil {
		return nil, err
	}

	for cur.Next(ctx) {
		var elem domain.Patient
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}

		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	cur.Close(ctx)
	return results, nil
}

func Connect(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return &mongo.Client{}, err
	}
	return client, nil
}
