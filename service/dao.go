package service

import (
	"gorm.io/gorm"
)

/*type DAO interface {
	CreateUser() error
	//CreatePlayer() error
	//CreateTask() error
	//GetTask(id string) Task
}

type dao struct {
	database database.Database2
	api      api.Api
}

func NewDAO(database database.Database2, api api.Api) DAO {
	return &dao{
		database: database,
		api:      api,
	}
}*/

// Define database interface (for DI)
type Database interface {
	Create(value interface{}) error
	Where(query interface{}, args ...interface{}) Database
	First(out interface{}, where ...interface{}) error
	Save(value interface{}) error
	Model(value interface{}) Database
	Take(out interface{}, where ...interface{}) error
	Delete(value interface{}, where ...interface{}) error
	Find(out interface{}, where ...interface{}) error
	Updates(values interface{}) error
}

// database access object
type GormDB struct {
	DB *gorm.DB
}

func (g *GormDB) Create(value interface{}) error {
	return g.DB.Create(value).Error
}

func (g *GormDB) Where(query interface{}, args ...interface{}) Database {
	return &GormDB{DB: g.DB.Where(query, args...)}
}

func (g *GormDB) First(out interface{}, where ...interface{}) error {
	return g.DB.First(out, where...).Error
}

func (g *GormDB) Save(value interface{}) error {
	return g.DB.Save(value).Error
}

func (g *GormDB) Model(value interface{}) Database {
	g.DB = g.DB.Model(value)
	return g
}

func (g *GormDB) Take(out interface{}, where ...interface{}) error {
	return g.DB.Take(out, where...).Error
}

func (g *GormDB) Delete(value interface{}, where ...interface{}) error {
	return g.DB.Delete(value, where...).Error
}

func (g *GormDB) Find(out interface{}, where ...interface{}) error {
	return g.DB.Find(out, where...).Error
}

func (g *GormDB) Updates(values interface{}) error {
	return g.DB.Updates(values).Error
}

// func (d *dao) CreateUser() error {
// 	// save into postgres
// 	d.database.GetDB().Create(model)
// 	// d.database.GetPostgresDB().Create(model)
// }

// func (d *dao) CreatePlayer() error {
// 	// save into mongodb
// 	d.database.GetMongoDB().Create(model)
// }

// func (d *dao) CreateTask() error {
// 	// save mysql
// 	d.database.GetMongoMySQLDB().Create(model)
// }

// func (d *dao) GetTask(id string) Task {
// 	// remote api server
// 	d.api.GetTak(id)
// }
