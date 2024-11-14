package database

import (
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // Import PostgreSQL dialect

	"test-task/shared/config"
	"test-task/shared/utils"
)

// ITransaction is
type IConnection interface {
	GetDB() *gorm.DB // this is for get the database with two mode with transaction or without transaction

	CommitTransaction()   // commit the transaction
	RollbackTransaction() // rollback transaction
	RollbackOnException() // for emergency rollback

	CreateNew(obj interface{}) error   // insert new data
	SaveChanges(obj interface{}) error // update existing data
}

// GormDB is
type connection struct {
	db *gorm.DB

	readonly      bool
	isTranscation bool
}

var db *gorm.DB
var slaveDB *gorm.DB

var connectionOnce sync.Once

func Init(config config.IConfig) {
	connectionOnce.Do(func() {
		var err error

		connectionString := config.Database().ConnectionString
		dialect := config.Database().Dialect
		logMode := config.Database().LogMode

		db, err = gorm.Open(dialect, connectionString)
		if err != nil {
			panic(err)
		}
		db.LogMode(logMode)

		if len(config.Database().SlaveConnectionString) > 0 {
			connectionString = config.Database().SlaveConnectionString
		}

		slaveDB, err = gorm.Open(dialect, connectionString)
		if err != nil {
			panic(err)
		}
		slaveDB.LogMode(logMode)
	})
}

func NewConnection() IConnection {
	return &connection{
		db,
		false,
		true,
	}
}

func NewTransaction() IConnection {
	return &connection{
		db.Begin(),
		false,
		true,
	}
}

func NewSlaveConnection() IConnection {
	return &connection{
		slaveDB,
		true,
		false,
	}
}

func Close() error {
	if db != nil {
		return db.Close()
	}
	if slaveDB != nil {
		return slaveDB.Close()
	}
	return nil
}

// GetDB is
func (self *connection) GetDB() *gorm.DB {
	return self.db
}

// Commitconnection is
func (self *connection) CommitTransaction() {
	if !self.isTranscation {
		panic("Commit on a non-transcation connection")
	}
	self.db.Commit()
}

// RollbackTransaction is
func (self *connection) RollbackTransaction() {
	if !self.isTranscation {
		panic("Rollback on a non-transcation connection")
	}
	self.db.Rollback()
}

// RollbackOnException is common handler for rollback the transaction
// to avoid database lock when something goes wrong in transaction state
// use with defer right after we call GetDB(true)
func (self *connection) RollbackOnException() {
	if !self.isTranscation {
		panic("Rollback on a non-transcation connection")
	}

	// catch the error
	if err := recover(); err != nil {
		// rollback it!
		self.db.Rollback()
		// repanic so we can get where it happen in log!
		panic(err)
	}
}

// CreateNew is
func (self *connection) CreateNew(obj interface{}) error {
	if self.readonly {
		return utils.NewError("Disallow write on slave connection!")
	}
	return self.db.Create(obj).Error
}

// SaveChanges is
func (self *connection) SaveChanges(obj interface{}) error {
	if self.readonly {
		return utils.NewError("Disallow write on slave connection!")
	}
	return self.db.Save(obj).Error
}
