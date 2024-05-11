package infra

import (
	"context"
	"testing"
	"time"

	"github.com/beriloqueiroz/psi-mgnt/internal/application"
	domain "github.com/beriloqueiroz/psi-mgnt/internal/domain/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

var ctx context.Context
var mongoRepo *MongoSessionRepository

func after() {
	clear()
	defer mongoRepo.Client.Disconnect(ctx)
}

func before() {
	ctx = context.Background()
	var err error

	mongoRepo, err = NewMongoSessionRepository(
		ctx,
		"mongodb://root:example@localhost:27017",
		"patients",
		"sessions",
		"psimgnt_test",
	)
	if err != nil {
		panic(err)
	}
	clear()
}

func clear() {
	_, err := mongoRepo.PatientCollection.DeleteMany(ctx, bson.D{{}})
	if err != nil {
		panic(err)
	}
	_, err = mongoRepo.SessionCollection.DeleteMany(ctx, bson.D{{}})
	if err != nil {
		panic(err)
	}
}

func TestCreateSession_WhenPatientAlreadyExist(t *testing.T) {
	before()
	defer after()
	patient, err := domain.NewPatient(
		uuid.NewString(),
		"berilo",
		"12365478",
		"",
		[]domain.Phone{},
		"123",
	)

	if err != nil {
		panic(err)
	}

	err = mongoRepo.CreatePatient(ctx, patient)

	if err != nil {
		panic(err)
	}

	id := uuid.NewString()

	session, err := domain.NewSession(id, 100, "notes de doido", time.Now(), time.Now(), time.Hour, patient, "123")

	if err != nil {
		panic(err)
	}

	err = mongoRepo.Create(ctx, session)

	if err != nil {
		panic(err)
	}

	found, err := mongoRepo.FindPatientByName(ctx, application.FindPatientByNameRepositoryInput{
		OwnerId: "123",
		Name:    "berilo",
	})

	if err != nil {
		panic(err)
	}

	assert.Nil(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, session.Patient.ID, found.ID)
	assert.Equal(t, patient.OwnerId, found.OwnerId)
}

func TestCreateSession_WhenPatientNotExist(t *testing.T) {
	before()
	defer after()
	patient, err := domain.NewPatient(
		uuid.NewString(),
		"berilo",
		"12365478",
		"",
		[]domain.Phone{},
		"123",
	)

	if err != nil {
		panic(err)
	}

	id := uuid.NewString()

	session, err := domain.NewSession(id, 100, "notes de doido", time.Now(), time.Now(), time.Hour, patient, "123")

	if err != nil {
		panic(err)
	}

	err = mongoRepo.Create(ctx, session)

	if err != nil {
		panic(err)
	}

	inputList := application.ListRepositoryInput{
		OwnerId:  "123",
		PageSize: 10,
		Page:     1,
	}

	list, err := mongoRepo.List(ctx, inputList)

	if err != nil {
		panic(err)
	}

	assert.Nil(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, session.ID, list[0].ID)
	assert.Equal(t, session.Notes, list[0].Notes)
	assert.Equal(t, session.OwnerId, list[0].OwnerId)
}

func TestFindPatientByName(t *testing.T) {
	before()
	defer after()
	patient1, err := domain.NewPatient(
		uuid.NewString(),
		"berilo",
		"12365478",
		"",
		[]domain.Phone{},
		"123",
	)

	patient2, err := domain.NewPatient(
		uuid.NewString(),
		"berilo",
		"12365478",
		"",
		[]domain.Phone{},
		"1234",
	)

	if err != nil {
		panic(err)
	}

	err = mongoRepo.CreatePatient(ctx, patient1)

	if err != nil {
		panic(err)
	}

	err = mongoRepo.CreatePatient(ctx, patient2)

	if err != nil {
		panic(err)
	}

	input := application.FindPatientByNameRepositoryInput{
		OwnerId: "1234",
		Name:    "berilo",
	}

	found, err := mongoRepo.FindPatientByName(ctx, input)

	if err != nil {
		panic(err)
	}

	assert.Nil(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, patient2.Name, found.Name)
	assert.Equal(t, patient2.OwnerId, found.OwnerId)
}

func TestSearchPatientsByTermName(t *testing.T) {
	before()
	defer after()
	patient1, err := domain.NewPatient(
		uuid.NewString(),
		"berilo jose",
		"12365478",
		"",
		[]domain.Phone{},
		"123",
	)
	if err != nil {
		panic(err)
	}
	patient2, err := domain.NewPatient(
		uuid.NewString(),
		"berilo grande",
		"12365478",
		"",
		[]domain.Phone{},
		"123",
	)
	patientNotOwner, err := domain.NewPatient(
		uuid.NewString(),
		"berilo grande",
		"12365478",
		"",
		[]domain.Phone{},
		"1234",
	)
	if err != nil {
		panic(err)
	}
	patient3, err := domain.NewPatient(
		uuid.NewString(),
		"não é ",
		"12365478",
		"",
		[]domain.Phone{},
		"123",
	)

	if err != nil {
		panic(err)
	}

	err = mongoRepo.CreatePatient(ctx, patient1)
	if err != nil {
		panic(err)
	}
	err = mongoRepo.CreatePatient(ctx, patient2)
	if err != nil {
		panic(err)
	}
	err = mongoRepo.CreatePatient(ctx, patient3)
	if err != nil {
		panic(err)
	}
	err = mongoRepo.CreatePatient(ctx, patientNotOwner)
	if err != nil {
		panic(err)
	}

	inputSearch := application.SearchPatientsByNameRepositoryInput{
		OwnerId:  "123",
		PageSize: 10,
		Page:     1,
		Term:     "berilo",
	}

	founds, err := mongoRepo.SearchPatientsByName(ctx, inputSearch)

	if err != nil {
		panic(err)
	}

	assert.Nil(t, err)
	assert.NotNil(t, founds)
	assert.Equal(t, len(founds), 2)
	assert.Equal(t, patient1.Name, founds[0].Name)
	assert.Equal(t, patient1.OwnerId, founds[0].OwnerId)
	assert.Equal(t, patient2.Name, founds[1].Name)
	assert.Equal(t, patient2.OwnerId, founds[1].OwnerId)

	inputSearch = application.SearchPatientsByNameRepositoryInput{
		OwnerId:  "123",
		PageSize: 10,
		Page:     1,
		Term:     "é",
	}

	founds, err = mongoRepo.SearchPatientsByName(ctx, inputSearch)

	if err != nil {
		panic(err)
	}

	assert.Nil(t, err)
	assert.NotNil(t, founds)
	assert.Equal(t, len(founds), 1)
	assert.Equal(t, patient3.Name, founds[0].Name)
}

func TestListPatientsByTermName(t *testing.T) {
	before()
	defer after()
	patient1, err := domain.NewPatient(
		uuid.NewString(),
		"berilo jose",
		"12365478",
		"",
		[]domain.Phone{},
		"123",
	)
	if err != nil {
		panic(err)
	}

	session1, err := domain.NewSession(uuid.NewString(), 100, "notes de doido", time.Now(), time.Now(), time.Hour, patient1, "123")
	if err != nil {
		panic(err)
	}
	session2, err := domain.NewSession(uuid.NewString(), 100, "notes de doido 1", time.Now(), time.Now(), time.Hour, patient1, "123")
	if err != nil {
		panic(err)
	}
	session3, err := domain.NewSession(uuid.NewString(), 100, "notes de doido 2", time.Now(), time.Now(), time.Hour, patient1, "123")
	if err != nil {
		panic(err)
	}

	err = mongoRepo.Create(ctx, session1)
	if err != nil {
		panic(err)
	}
	err = mongoRepo.Create(ctx, session2)
	if err != nil {
		panic(err)
	}
	err = mongoRepo.Create(ctx, session3)

	if err != nil {
		panic(err)
	}

	inputList := application.ListRepositoryInput{
		OwnerId:  "123",
		PageSize: 10,
		Page:     1,
	}
	founds, err := mongoRepo.List(ctx, inputList)

	if err != nil {
		panic(err)
	}

	assert.Nil(t, err)
	assert.NotNil(t, founds)
	assert.Equal(t, len(founds), 3)
	assert.Equal(t, session1.ID, founds[0].ID)
	assert.Equal(t, session2.ID, founds[1].ID)
	assert.Equal(t, session3.ID, founds[2].ID)
}

func TestDeleteSession(t *testing.T) {
	before()
	defer after()
	patient1, err := domain.NewPatient(
		uuid.NewString(),
		"berilo jose",
		"12365478",
		"",
		[]domain.Phone{},
		"123",
	)
	if err != nil {
		panic(err)
	}

	session1, err := domain.NewSession(uuid.NewString(), 100, "notes de doido", time.Now(), time.Now(), time.Hour, patient1, "123")
	if err != nil {
		panic(err)
	}
	session2, err := domain.NewSession(uuid.NewString(), 100, "notes de doido 1", time.Now(), time.Now(), time.Hour, patient1, "123")
	if err != nil {
		panic(err)
	}
	session3, err := domain.NewSession(uuid.NewString(), 100, "notes de doido 2", time.Now(), time.Now(), time.Hour, patient1, "123")
	if err != nil {
		panic(err)
	}

	err = mongoRepo.Create(ctx, session1)
	if err != nil {
		panic(err)
	}
	err = mongoRepo.Create(ctx, session2)
	if err != nil {
		panic(err)
	}
	err = mongoRepo.Create(ctx, session3)

	if err != nil {
		panic(err)
	}

	inputList := application.ListRepositoryInput{
		OwnerId:  "123",
		PageSize: 10,
		Page:     1,
	}

	founds, err := mongoRepo.List(ctx, inputList)

	if err != nil {
		panic(err)
	}

	assert.Nil(t, err)
	assert.NotNil(t, founds)
	assert.Equal(t, len(founds), 3)
	assert.Equal(t, session1.ID, founds[0].ID)
	assert.Equal(t, session2.ID, founds[1].ID)
	assert.Equal(t, session3.ID, founds[2].ID)

	input := application.DeleteRepositoryInput{
		OwnerId: session2.OwnerId,
		Id:      session2.ID,
	}

	err = mongoRepo.Delete(ctx, input)

	if err != nil {
		panic(err)
	}

	founds, err = mongoRepo.List(ctx, inputList)

	if err != nil {
		panic(err)
	}

	assert.Nil(t, err)
	assert.NotNil(t, founds)
	assert.Equal(t, len(founds), 2)
	assert.Equal(t, session1.ID, founds[0].ID)
	assert.Equal(t, session3.ID, founds[1].ID)

}
