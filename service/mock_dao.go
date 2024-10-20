package service

import (
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) GetDB() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB) // Return a mock GORM DB instance
}

func (m *MockDB) Init() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockDB) Create(value interface{}) error {
	args := m.Called(value)
	return args.Error(0)
}

func (m *MockDB) Where(query interface{}, args ...interface{}) Database {
	//argsCalled := m.Called(query, args)
	//return argsCalled.Get(0).(Database)
	m.Called(query, args)
	return m
}

// Mock First method
func (m *MockDB) First(out interface{}, where ...interface{}) error {
	args := m.Called(out)
	return args.Error(0)
}

func (m *MockDB) Save(value interface{}) error {
	args := m.Called(value)
	return args.Error(0)
}

// Mock the Model method
func (m *MockDB) Model(value interface{}) Database {
	m.Called(value)
	return m
}

// Mock the Take method
func (m *MockDB) Take(out interface{}, where ...interface{}) error {
	args := m.Called(out)
	return args.Error(0)

	/*args := m.Called(out)
	user := out.(*models.User)
	user.ID = 1
	user.Username = "testuser"
	user.Password = "testuserpassword"
	user.Role = "user"
	return args.Error(0)*/
}

func (m *MockDB) Delete(value interface{}, where ...interface{}) error {
	args := m.Called(value, where)
	return args.Error(0)
}

func (m *MockDB) Find(out interface{}, where ...interface{}) error {
	args := m.Called(out)
	return args.Error(0)
}

func (m *MockDB) Updates(values interface{}) error {
	args := m.Called(values)
	return args.Error(0)
}

/*type mockDAO struct {
	users []int
}

func (d *mockDAO) CreateUser() error {
	d.users = append(d.users, 1)
	return nil
}

func (d *mockDAO) CreatePlayer() error {
	return nil
}

func (d *mockDAO) CreateTask() error {
	return nil
}

func (d *mockDAO) GetTask(id string) Task {
	return nil
}
*/
