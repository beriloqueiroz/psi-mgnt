package infra

import (
	"context"
	"testing"
	"time"

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
	)

	if err != nil {
		panic(err)
	}

	err = mongoRepo.CreatePatient(ctx, patient)

	if err != nil {
		panic(err)
	}

	id := uuid.NewString()

	session, err := domain.NewSession(id, 100, "notes de doido", time.Now(), time.Now(), time.Hour, patient)

	if err != nil {
		panic(err)
	}

	err = mongoRepo.Create(ctx, session)

	if err != nil {
		panic(err)
	}

	found, err := mongoRepo.FindPatientByName(ctx, "berilo")

	if err != nil {
		panic(err)
	}

	assert.Nil(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, session.Patient.ID, found.ID)
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
	)

	if err != nil {
		panic(err)
	}

	id := uuid.NewString()

	session, err := domain.NewSession(id, 100, "notes de doido", time.Now(), time.Now(), time.Hour, patient)

	if err != nil {
		panic(err)
	}

	err = mongoRepo.Create(ctx, session)

	if err != nil {
		panic(err)
	}

	list, err := mongoRepo.List(ctx, 10, 0)

	if err != nil {
		panic(err)
	}

	assert.Nil(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, session.ID, list[0].ID)
	assert.Equal(t, session.Notes, list[0].Notes)
}

func TestFinsPatientByName(t *testing.T) {
	before()
	defer after()
	patient1, err := domain.NewPatient(
		uuid.NewString(),
		"berilo",
		"12365478",
		"",
		[]domain.Phone{},
	)

	if err != nil {
		panic(err)
	}

	err = mongoRepo.CreatePatient(ctx, patient1)

	if err != nil {
		panic(err)
	}
	found, err := mongoRepo.FindPatientByName(ctx, "berilo")

	if err != nil {
		panic(err)
	}

	assert.Nil(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, patient1.Name, found.Name)
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
	founds, err := mongoRepo.SearchPatientsByName(ctx, "berilo", 10, 0)

	if err != nil {
		panic(err)
	}

	assert.Nil(t, err)
	assert.NotNil(t, founds)
	assert.Equal(t, len(founds), 2)
	assert.Equal(t, patient1.Name, founds[0].Name)
	assert.Equal(t, patient2.Name, founds[1].Name)

	founds, err = mongoRepo.SearchPatientsByName(ctx, "é", 10, 0)

	if err != nil {
		panic(err)
	}

	assert.Nil(t, err)
	assert.NotNil(t, founds)
	assert.Equal(t, len(founds), 1)
	assert.Equal(t, patient3.Name, founds[0].Name)
}
