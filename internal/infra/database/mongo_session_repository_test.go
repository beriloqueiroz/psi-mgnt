package infra

import (
	"context"
	"fmt"
	"github.com/beriloqueiroz/psi-mgnt/pkg/helpers"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"os"
	"testing"
	"time"

	"github.com/beriloqueiroz/psi-mgnt/internal/application"
	domain "github.com/beriloqueiroz/psi-mgnt/internal/domain/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var ctx context.Context
var mongoRepo *MongoSessionRepository
var mongoContainer testcontainers.Container

const mongoUser = "root"
const mongoPassword = "example"
const mongoDatabase = "psimgnt"

func TestMain(m *testing.M) {
	setup()
	defer teardown()

	// Run all the tests
	code := m.Run()

	// Exit with the received code
	os.Exit(code)
}

func setup() {
	log.Println("setup suite")
	ctx = context.Background()
	req := testcontainers.ContainerRequest{
		Image: "docker.io/mongo",
		Env:   map[string]string{"MONGO_INITDB_ROOT_USERNAME": mongoUser, "MONGO_INITDB_ROOT_PASSWORD": mongoPassword, "MONGO_INITDB_DATABASE": mongoDatabase},
	}
	var err error

	mongoContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		log.Fatal(err)
	}
}

func teardown() {
	log.Println("teardown suite")
	func() {
		if err := mongoContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()
}

func before() {
	log.Println("setup test")
	host, err := mongoContainer.PortEndpoint(ctx, nat.Port("27017"), "")
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	uri := fmt.Sprintf("mongodb://%s:%s@%s", mongoUser, mongoPassword, host)
	mongoRepo, err = NewMongoSessionRepository(
		ctx,
		uri,
		"patients",
		"professionals",
		"sessions",
		mongoDatabase,
	)
	if err != nil {
		panic(err)
	}
	clearDb()
}

func after() {
	log.Println("teardown test")
	clearDb()
	mongoRepo.Client.Disconnect(ctx)
}

func clearDb() {
	_, err := mongoRepo.PatientCollection.DeleteMany(ctx, bson.D{{}})
	if err != nil {
		panic(err)
	}
	_, err = mongoRepo.SessionCollection.DeleteMany(ctx, bson.D{{}})
	if err != nil {
		panic(err)
	}
	_, err = mongoRepo.ProfessionalCollection.DeleteMany(ctx, bson.D{{}})
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

	assert.Nil(t, err)

	err = mongoRepo.CreatePatient(ctx, patient)

	assert.Nil(t, err)

	professional, err := domain.NewProfessional(
		uuid.NewString(),
		"berilo",
		"12365478",
		"",
	)

	assert.Nil(t, err)

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

	assert.Nil(t, err)

	err = mongoRepo.Create(ctx, session)

	assert.Nil(t, err)

	found, err := mongoRepo.FindPatient(ctx, application.FindPatientRepositoryInput{
		PatientId: patientId,
	})

	assert.Nil(t, err)

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

	assert.Nil(t, err)

	professionalId := uuid.New().String()
	professional, err := domain.NewProfessional(
		professionalId,
		"berilo",
		"12365478",
		"",
	)

	assert.Nil(t, err)

	id := uuid.NewString()

	session, err := domain.NewSession(id, 100, "notes de doido", time.Now(), time.Hour, patient, "unimed", professional)

	assert.Nil(t, err)

	err = mongoRepo.Create(ctx, session)

	assert.Nil(t, err)

	listConfig := helpers.ListConfig{
		PageSize: 10,
		Page:     1,
	}

	inputList := application.ListByProfessionalRepositoryInput{
		ProfessionalId: professionalId,
		ListConfig:     listConfig,
	}

	list, err := mongoRepo.ListByProfessional(ctx, inputList)

	assert.Nil(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, session.ID, list.Content[0].ID)
	assert.Equal(t, session.Notes, list.Content[0].Notes)
	assert.Equal(t, session.Professional, list.Content[0].Professional)
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
	assert.Nil(t, err)

	professionalId := uuid.New().String()
	professional, err := domain.NewProfessional(
		professionalId,
		"berilo",
		"12365478",
		"",
	)

	assert.Nil(t, err)

	session1, err := domain.NewSession(uuid.NewString(), 100, "notes de doido", time.Now(), time.Hour, patient1, "unimed", professional)
	assert.Nil(t, err)
	session2, err := domain.NewSession(uuid.NewString(), 100, "notes de doido 1", time.Now(), time.Hour, patient1, "unimed", professional)
	assert.Nil(t, err)
	session3, err := domain.NewSession(uuid.NewString(), 100, "notes de doido 2", time.Now(), time.Hour, patient1, "unimed", professional)
	assert.Nil(t, err)

	err = mongoRepo.Create(ctx, session1)
	assert.Nil(t, err)
	err = mongoRepo.Create(ctx, session2)
	assert.Nil(t, err)
	err = mongoRepo.Create(ctx, session3)
	assert.Nil(t, err)

	listConfig := helpers.ListConfig{
		PageSize: 10,
		Page:     1,
	}

	inputList := application.ListByProfessionalRepositoryInput{
		ProfessionalId: professionalId,
		ListConfig:     listConfig,
	}

	founds, err := mongoRepo.ListByProfessional(ctx, inputList)

	assert.Nil(t, err)
	assert.NotNil(t, founds)
	assert.Equal(t, len(founds.Content), 3)
	assert.Equal(t, session1.ID, founds.Content[0].ID)
	assert.Equal(t, session2.ID, founds.Content[1].ID)
	assert.Equal(t, session3.ID, founds.Content[2].ID)

	input := application.DeleteRepositoryInput{
		SessionId: session2.ID,
	}

	err = mongoRepo.Delete(ctx, input)

	assert.Nil(t, err)

	founds, err = mongoRepo.ListByProfessional(ctx, inputList)

	assert.Nil(t, err)
	assert.NotNil(t, founds)
	assert.Equal(t, len(founds.Content), 2)
	assert.Equal(t, session1.ID, founds.Content[0].ID)
	assert.Equal(t, session3.ID, founds.Content[1].ID)
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
	assert.Nil(t, err)

	professional, err := domain.NewProfessional(
		uuid.NewString(),
		"berilo",
		"12365478",
		"",
	)

	assert.Nil(t, err)

	session1, err := domain.NewSession(uuid.NewString(), 100, "notes de doido", time.Now(), time.Hour, patient1, "unimed", professional)
	assert.Nil(t, err)

	session2, err := domain.NewSession(uuid.NewString(), 100, "notes de doido 1", time.Now(), time.Hour, patient1, "unimed", professional)
	assert.Nil(t, err)

	session3, err := domain.NewSession(uuid.NewString(), 100, "notes de doido 2", time.Now(), time.Hour, patient1, "unimed", professional)
	assert.Nil(t, err)

	err = mongoRepo.Create(ctx, session1)
	assert.Nil(t, err)

	err = mongoRepo.Create(ctx, session2)
	assert.Nil(t, err)

	err = mongoRepo.Create(ctx, session3)
	assert.Nil(t, err)

	listConfig := helpers.ListConfig{
		PageSize: 10,
		Page:     1,
	}

	inputList := application.ListRepositoryInput{
		ListConfig: listConfig,
	}
	founds, err := mongoRepo.List(ctx, inputList)

	assert.Nil(t, err)
	assert.NotNil(t, founds)
	assert.Equal(t, len(founds.Content), 3)
	assert.Equal(t, session1.ID, founds.Content[0].ID)
	assert.Equal(t, session2.ID, founds.Content[1].ID)
	assert.Equal(t, session3.ID, founds.Content[2].ID)
}

func TestUpdateSession(t *testing.T) {
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

	assert.Nil(t, err)

	professionalId := uuid.New().String()
	professional, err := domain.NewProfessional(
		professionalId,
		"berilo",
		"12365478",
		"",
	)

	assert.Nil(t, err)

	id := uuid.NewString()

	session, err := domain.NewSession(id, 100, "notes de doido", time.Now(), time.Hour, patient, "unimed", professional)

	assert.Nil(t, err)
	err = mongoRepo.Create(ctx, session)

	assert.Nil(t, err)

	session.Notes = "12365478 notes"

	err = mongoRepo.Update(ctx, session)

	assert.Nil(t, err)

	sessionUpdated, err := mongoRepo.Find(ctx, application.FindSessionRepositoryInput{ID: id})
	assert.Nil(t, err)

	assert.Equal(t, session.ID, sessionUpdated.ID)
	assert.Equal(t, session.Notes, sessionUpdated.Notes)

}

func TestFindSession(t *testing.T) {
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

	assert.Nil(t, err)

	professionalId := uuid.New().String()
	professional, err := domain.NewProfessional(
		professionalId,
		"berilo",
		"12365478",
		"",
	)

	assert.Nil(t, err)

	id := uuid.NewString()

	session, err := domain.NewSession(id, 100, "notes de doido", time.Now(), time.Hour, patient, "unimed", professional)

	assert.Nil(t, err)
	err = mongoRepo.Create(ctx, session)

	assert.Nil(t, err)

	sessionFound, err := mongoRepo.Find(ctx, application.FindSessionRepositoryInput{ID: id})
	assert.Nil(t, err)

	assert.Equal(t, session.ID, sessionFound.ID)
	assert.Equal(t, session.Notes, sessionFound.Notes)
}

// others

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

	listConfig := helpers.ListConfig{
		PageSize: 10,
		Page:     1,
		ExpressionFilters: []helpers.ExpressionFilter{
			{PropertyName: "name", Value: "berilo"},
		},
	}

	inputSearch := application.SearchPatientsByNameRepositoryInput{
		ListConfig: listConfig,
	}

	founds, err := mongoRepo.SearchPatientsByName(ctx, inputSearch)

	if err != nil {
		panic(err)
	}

	assert.Nil(t, err)
	assert.NotNil(t, founds)
	assert.Equal(t, 3, len(founds.Content))
	assert.Equal(t, patient1.Name, founds.Content[0].Name)
	assert.Equal(t, patient2.Name, founds.Content[1].Name)
	assert.Equal(t, patient3.Name, founds.Content[2].Name)

	listConfig = helpers.ListConfig{
		PageSize: 10,
		Page:     1,
		ExpressionFilters: []helpers.ExpressionFilter{
			{PropertyName: "name", Value: "é"},
		},
	}

	inputSearch = application.SearchPatientsByNameRepositoryInput{
		ListConfig: listConfig,
	}

	founds, err = mongoRepo.SearchPatientsByName(ctx, inputSearch)

	if err != nil {
		panic(err)
	}

	assert.Nil(t, err)
	assert.NotNil(t, founds)
	assert.Equal(t, len(founds.Content), 1)
	assert.Equal(t, patient4.Name, founds.Content[0].Name)
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

	listConfig := helpers.ListConfig{
		PageSize: 10,
		Page:     1,
	}

	inputList := application.ListByProfessionalRepositoryInput{
		ProfessionalId: professionalId,
		ListConfig:     listConfig,
	}
	founds, err := mongoRepo.ListByProfessional(ctx, inputList)

	if err != nil {
		panic(err)
	}

	assert.Nil(t, err)
	assert.NotNil(t, founds)
	assert.Equal(t, len(founds.Content), 3)
	assert.Equal(t, session1.ID, founds.Content[0].ID)
	assert.Equal(t, session2.ID, founds.Content[1].ID)
	assert.Equal(t, session3.ID, founds.Content[2].ID)
}

func TestCreateProfessional(t *testing.T) {
	before()
	defer after()
	professionalId := uuid.New().String()
	professional, err := domain.NewProfessional(
		professionalId,
		"berilo",
		"12365478",
		"berio@gmail.com",
	)

	if err != nil {
		panic(err)
	}

	err = mongoRepo.CreateProfessional(ctx, professional)

	if err != nil {
		panic(err)
	}

	assert.Nil(t, err)
}
