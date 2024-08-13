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
		"professionals",
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
	patientId := uuid.New().String()
	patient, err := domain.NewPatient(
		patientId,
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

	professional, err := domain.NewProfessional(
		uuid.NewString(),
		"berilo",
		"12365478",
		"",
	)

	if err != nil {
		panic(err)
	}

	id := uuid.NewString()

	session, err := domain.NewSession(
		id,
		100,
		"notes de doido",
		time.Now(),
		time.Hour,
		patient,
		"unimed",
		professional)

	if err != nil {
		panic(err)
	}

	err = mongoRepo.Create(ctx, session)

	if err != nil {
		panic(err)
	}

	found, err := mongoRepo.FindPatient(ctx, application.FindPatientRepositoryInput{
		PatientId: patientId,
	})

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
	patientId := uuid.New().String()
	patient, err := domain.NewPatient(
		patientId,
		"berilo",
		"12365478",
		"",
		[]domain.Phone{},
	)

	if err != nil {
		panic(err)
	}

	professionalId := uuid.New().String()
	professional, err := domain.NewProfessional(
		professionalId,
		"berilo",
		"12365478",
		"",
	)

	if err != nil {
		panic(err)
	}

	id := uuid.NewString()

	session, err := domain.NewSession(id, 100, "notes de doido", time.Now(), time.Hour, patient, "unimed", professional)

	if err != nil {
		panic(err)
	}

	err = mongoRepo.Create(ctx, session)

	if err != nil {
		panic(err)
	}

	inputList := application.ListByProfessionalRepositoryInput{
		ProfessionalId: professionalId,
		PageSize:       10,
		Page:           1,
	}

	list, err := mongoRepo.ListByProfessional(ctx, inputList)

	if err != nil {
		panic(err)
	}

	assert.Nil(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, session.ID, list[0].ID)
	assert.Equal(t, session.Notes, list[0].Notes)
	assert.Equal(t, session.Professional, list[0].Professional)
}

func TestFindPatient(t *testing.T) {
	before()
	defer after()
	patientId1 := uuid.NewString()
	patientId2 := uuid.NewString()
	patient1, err := domain.NewPatient(
		patientId1,
		"berilo",
		"12365478",
		"",
		[]domain.Phone{},
	)

	patient2, err := domain.NewPatient(
		patientId2,
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

	err = mongoRepo.CreatePatient(ctx, patient2)

	if err != nil {
		panic(err)
	}

	input := application.FindPatientRepositoryInput{
		PatientId: patientId2,
	}

	found, err := mongoRepo.FindPatient(ctx, input)

	if err != nil {
		panic(err)
	}

	assert.Nil(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, patient2.Name, found.Name)
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
	patient3, err := domain.NewPatient(
		uuid.NewString(),
		"berilo grande",
		"12365478",
		"",
		[]domain.Phone{},
	)
	if err != nil {
		panic(err)
	}
	patient4, err := domain.NewPatient(
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
	err = mongoRepo.CreatePatient(ctx, patient4)
	if err != nil {
		panic(err)
	}

	inputSearch := application.SearchPatientsByNameRepositoryInput{
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
	assert.Equal(t, 3, len(founds))
	assert.Equal(t, patient1.Name, founds[0].Name)
	assert.Equal(t, patient2.Name, founds[1].Name)
	assert.Equal(t, patient3.Name, founds[2].Name)

	inputSearch = application.SearchPatientsByNameRepositoryInput{
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
	assert.Equal(t, patient4.Name, founds[0].Name)
}

func TestListByProfessionalSessions(t *testing.T) {
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

	professionalId := uuid.New().String()
	professional, err := domain.NewProfessional(
		professionalId,
		"berilo",
		"12365478",
		"",
	)

	if err != nil {
		panic(err)
	}

	session1, err := domain.NewSession(uuid.NewString(), 100, "notes de doido", time.Now(), time.Hour, patient1, "unimed", professional)
	if err != nil {
		panic(err)
	}
	session2, err := domain.NewSession(uuid.NewString(), 100, "notes de doido 1", time.Now(), time.Hour, patient1, "unimed", professional)
	if err != nil {
		panic(err)
	}
	session3, err := domain.NewSession(uuid.NewString(), 100, "notes de doido 2", time.Now(), time.Hour, patient1, "unimed", professional)
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

	inputList := application.ListByProfessionalRepositoryInput{
		ProfessionalId: professionalId,
		PageSize:       10,
		Page:           1,
	}
	founds, err := mongoRepo.ListByProfessional(ctx, inputList)

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
	)
	if err != nil {
		panic(err)
	}

	professionalId := uuid.New().String()
	professional, err := domain.NewProfessional(
		professionalId,
		"berilo",
		"12365478",
		"",
	)

	if err != nil {
		panic(err)
	}

	session1, err := domain.NewSession(uuid.NewString(), 100, "notes de doido", time.Now(), time.Hour, patient1, "unimed", professional)
	if err != nil {
		panic(err)
	}
	session2, err := domain.NewSession(uuid.NewString(), 100, "notes de doido 1", time.Now(), time.Hour, patient1, "unimed", professional)
	if err != nil {
		panic(err)
	}
	session3, err := domain.NewSession(uuid.NewString(), 100, "notes de doido 2", time.Now(), time.Hour, patient1, "unimed", professional)
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

	inputList := application.ListByProfessionalRepositoryInput{
		ProfessionalId: professionalId,
		PageSize:       10,
		Page:           1,
	}

	founds, err := mongoRepo.ListByProfessional(ctx, inputList)

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
		SessionId: session2.ID,
	}

	err = mongoRepo.Delete(ctx, input)

	if err != nil {
		panic(err)
	}

	founds, err = mongoRepo.ListByProfessional(ctx, inputList)

	if err != nil {
		panic(err)
	}

	assert.Nil(t, err)
	assert.NotNil(t, founds)
	assert.Equal(t, len(founds), 2)
	assert.Equal(t, session1.ID, founds[0].ID)
	assert.Equal(t, session3.ID, founds[1].ID)
}

func TestListSessions(t *testing.T) {
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

	professional, err := domain.NewProfessional(
		uuid.NewString(),
		"berilo",
		"12365478",
		"",
	)

	if err != nil {
		panic(err)
	}

	session1, err := domain.NewSession(uuid.NewString(), 100, "notes de doido", time.Now(), time.Hour, patient1, "unimed", professional)
	if err != nil {
		panic(err)
	}
	session2, err := domain.NewSession(uuid.NewString(), 100, "notes de doido 1", time.Now(), time.Hour, patient1, "unimed", professional)
	if err != nil {
		panic(err)
	}
	session3, err := domain.NewSession(uuid.NewString(), 100, "notes de doido 2", time.Now(), time.Hour, patient1, "unimed", professional)
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
