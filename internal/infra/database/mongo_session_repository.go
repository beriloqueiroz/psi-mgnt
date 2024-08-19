package infra

import (
	"context"
	"fmt"
	"github.com/beriloqueiroz/psi-mgnt/internal/application"
	domain "github.com/beriloqueiroz/psi-mgnt/internal/domain/entity"
	"github.com/beriloqueiroz/psi-mgnt/pkg/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"slices"
)

type MongoSessionRepository struct {
	Client                 *mongo.Client
	PatientCollection      *mongo.Collection
	ProfessionalCollection *mongo.Collection
	SessionCollection      *mongo.Collection
}

func NewMongoSessionRepository(ctx context.Context, uri string, patientCollection string, professionalCollection string,
	sessionCollection string, databaseName string) (*MongoSessionRepository, error) {
	client, err := connectToMongoDb(uri, ctx)
	if err != nil {
		return nil, err
	}
	database := client.Database(databaseName)
	initCollections(ctx, database, patientCollection, professionalCollection, sessionCollection)
	PatientCollection := database.Collection(patientCollection)
	ProfessionalCollection := database.Collection(professionalCollection)
	SessionCollection := database.Collection(sessionCollection)
	return &MongoSessionRepository{
		client,
		PatientCollection,
		ProfessionalCollection,
		SessionCollection,
	}, nil
}

func initCollections(ctx context.Context, database *mongo.Database, patientCollection string, professionalCollection string, sessionCollection string) error {
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
	if !slices.Contains(result, professionalCollection) {
		err = database.CreateCollection(ctx, professionalCollection)
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
			}},
	}
	_, err := mr.SessionCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
func (mr *MongoSessionRepository) List(ctx context.Context, input application.ListRepositoryInput) (*helpers.Pages[domain.Session], error) {
	l := int64(input.ListConfig.PageSize)
	skip := int64(input.ListConfig.Page*input.ListConfig.PageSize - input.ListConfig.PageSize)
	findOptions := options.FindOptions{Limit: &l, Skip: &skip}

	var results []domain.Session
	cur, err := mr.SessionCollection.Find(ctx, bson.D{}, &findOptions)
	if err != nil {
		return nil, err
	}

	for cur.Next(ctx) {
		var elem domain.Session
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}

		results = append(results, elem)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	cur.Close(ctx)
	return &helpers.Pages[domain.Session]{
		Content:  results,
		Page:     input.ListConfig.Page,
		PageSize: input.ListConfig.PageSize,
		Size:     100,
	}, nil
}
func (m *MongoSessionRepository) Update(ctx context.Context, session *domain.Session) error {
	return nil
}
func (m *MongoSessionRepository) Find(ctx context.Context, input application.FindSessionRepositoryInput) (*domain.Session, error) {
	return nil, nil
}

func (mr *MongoSessionRepository) ListByProfessional(ctx context.Context, input application.ListByProfessionalRepositoryInput) (*helpers.Pages[domain.Session], error) {
	l := int64(input.ListConfig.PageSize)
	skip := int64(input.ListConfig.Page*input.ListConfig.PageSize - input.ListConfig.PageSize)
	findOptions := options.FindOptions{Limit: &l, Skip: &skip}

	var results []domain.Session
	cur, err := mr.SessionCollection.Find(ctx, bson.D{{Key: "professional.id", Value: input.ProfessionalId}}, &findOptions)
	if err != nil {
		return nil, err
	}

	for cur.Next(ctx) {
		var elem domain.Session
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}

		results = append(results, elem)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	cur.Close(ctx)
	return &helpers.Pages[domain.Session]{
		Content:  results,
		Page:     input.ListConfig.Page,
		PageSize: input.ListConfig.PageSize,
		Size:     100,
	}, nil
}
func (mr *MongoSessionRepository) FindPatient(ctx context.Context, input application.FindPatientRepositoryInput) (*domain.Patient, error) {
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"id", input.PatientId}},
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
func (mr *MongoSessionRepository) FindProfessional(ctx context.Context, input application.FindProfessionalRepositoryInput) (*domain.Professional, error) {
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"id", input.ProfessionalId}},
			}},
	}
	var result domain.Professional
	err := mr.ProfessionalCollection.FindOne(ctx, filter).Decode(&result)
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
func (mr *MongoSessionRepository) SearchPatientsByName(ctx context.Context, input application.SearchPatientsByNameRepositoryInput) (*helpers.Pages[domain.Patient], error) {
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{Key: input.ListConfig.ExpressionFilters[0].PropertyName, Value: primitive.Regex{
					Pattern: "/*" + input.ListConfig.ExpressionFilters[0].Value.(string) + ".*",
					Options: "i",
				}}},
			}},
	}
	l := int64(input.ListConfig.PageSize)
	skip := int64(input.ListConfig.Page*input.ListConfig.PageSize - input.ListConfig.PageSize)
	findOptions := options.FindOptions{Limit: &l, Skip: &skip}

	var results []domain.Patient
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

		results = append(results, elem)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	cur.Close(ctx)
	return &helpers.Pages[domain.Patient]{
		Content:  results,
		Page:     input.ListConfig.Page,
		PageSize: input.ListConfig.PageSize,
		Size:     100,
	}, nil
}
func (mr *MongoSessionRepository) SearchProfessionalsByName(ctx context.Context, input application.SearchProfessionalByNameRepositoryInput) (*helpers.Pages[domain.Professional], error) {
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{Key: input.ListConfig.ExpressionFilters[0].PropertyName, Value: primitive.Regex{
					Pattern: "/*" + input.ListConfig.ExpressionFilters[0].Value.(string) + ".*",
					Options: "i",
				}}},
			}},
	}
	l := int64(input.ListConfig.PageSize)
	skip := int64(input.ListConfig.Page*input.ListConfig.PageSize - input.ListConfig.PageSize)
	findOptions := options.FindOptions{Limit: &l, Skip: &skip}

	var results []domain.Professional
	cur, err := mr.ProfessionalCollection.Find(ctx, filter, &findOptions)
	if err != nil {
		return nil, err
	}

	for cur.Next(ctx) {
		var elem domain.Professional
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}

		results = append(results, elem)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	cur.Close(ctx)
	return &helpers.Pages[domain.Professional]{
		Content:  results,
		Page:     input.ListConfig.Page,
		PageSize: input.ListConfig.PageSize,
		Size:     100,
	}, nil
}
func (mr *MongoSessionRepository) CreateProfessional(ctx context.Context, professional *domain.Professional) error {
	_, err := mr.ProfessionalCollection.InsertOne(ctx, professional)
	if err != nil {
		return err
	}
	return nil
}
