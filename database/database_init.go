package database

import (
	"fmt"

	config "crossplatform_chatbot/configs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//var DB *gorm.DB

// Database interface defines the methods for database operations
type Database interface {
	Init() error
	GetDB() *gorm.DB
	//GetRedis() *redis.Redis
}

// database2 struct holds the connection details and the gorm DB instance
type database struct {
	conf *config.Config
	//user string
	//pwd  string
	db *gorm.DB
	//redis *redis.Redis
	//mongo *mongo.Mongo
}

// NewDatabase2 creates a new instance of database2 with the provided config
func NewDatabase(config *config.Config) Database {
	return &database{
		conf: config, // TODO: Need whole config?
		//user: config.GetDBUser(),
		//pwd:  config.GetDBPwd(),
	}
}

// Init initializes the database connection and performs migrations
func (db2 *database) Init() error {
	if err := db2.initPostgres(); err != nil {
		return err
	}
	/*if err := db2.initRedis(); err != nil {
		return err
	}*/
	fmt.Println("Database connected!")

	return nil

	// dbstr := db2.conf.DBString

	// db, err := gorm.Open(postgres.Open(dbstr), &gorm.Config{})
	// if err != nil {
	// 	return fmt.Errorf("failed to connect to database: %w", err)
	// }

	// // Auto migrate the User schema
	// if err := db.AutoMigrate(&models.User{}); err != nil {
	// 	return fmt.Errorf("migration failed: %w", err)
	// }

	// // Assign the initialized database connection to the db field
	// db2.db = db
	// fmt.Println("Database connected!")

	// return nil
}

// GetDB returns the gorm DB instance
func (db2 *database) GetDB() *gorm.DB {
	return db2.db
}

// Initialize the database connection (Postgres)
func (db2 *database) initPostgres() error {
	dbstr := db2.conf.DBString

	db, err := gorm.Open(postgres.Open(dbstr), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto migrate the User schema
	/*if err := db.AutoMigrate(&models.User{}); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}*/

	// Assign the initialized database connection to the db field
	db2.db = db

	return nil
}

// func (db2 *database) initRedis() error {
// 	if err := redisClient.Connect("url", "pwd"); err != nil {
// 		return fmt.Errorf("failed to connect to redis: %w", err)
// 	}

// 	// Assign the initialized database connection to the db field
// 	db2.redis = redis

// 	return nil
// }

//var Client *mongo.Client
//var ItemCollection *mongo.Collection

// Initialize the database connection (MongoDB)
/*
func InitMongoDB() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	Client = client
	ItemCollection = client.Database("testDB").Collection("items")
	log.Println("Connected to MongoDB!")

	return client, nil
}*/
