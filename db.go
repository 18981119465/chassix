package chassis

import (
	"errors"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm/logger"
	"log"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"c6x.io/chassis/config"
	"c6x.io/chassis/logx"
)

type MultiDBSource struct {
	lock sync.RWMutex
	dbs  []*gorm.DB
}

var (
	ErrNoDatabaseConfiguration = errors.New("there isn't any database setting in the configuration file")
)

var (
	multiDBSource *MultiDBSource
	initOnce      sync.Once
)

func initMultiDBSource() {
	initOnce.Do(func() {
		multiCfg := config.Databases()
		multiDBSource = new(MultiDBSource)
		multiDBSource.lock.Lock()
		defer multiDBSource.lock.Unlock()
		for _, v := range multiCfg {
			multiDBSource.dbs = append(multiDBSource.dbs, mustConnectDB(v))
		}
	})
}

func mustConnectDB(dbCfg *config.DatabaseConfig) *gorm.DB {
	log := logx.New().Service("chassis").Category("gorm")
	dialectConfig := dbCfg.Dialect
	var dialect gorm.Dialector
	switch dialectConfig {
	case "mysql":
		dialect = mysql.Open(dbCfg.DSN)
	case "":
		dialect = mysql.Open(dbCfg.DSN)	
	case "postgres":
		dialect = postgres.Open(dbCfg.DSN)
	case "sqlite3":
		dialect = sqlite.Open(dbCfg.DSN)
	default:
		log.Fatalln("database driver config error: invalid dialect " + dbCfg.DSN)
	}
	if "" == dialectConfig {
		dialect = mysql.Open(dbCfg.DSN)
	}
	gCfg := &gorm.Config{}
	if false == dbCfg.ShowSQL {
		gCfg.Logger = logger.Default.LogMode(logger.Silent)
	}

	db, err := gorm.Open(dialect, gCfg)
	if err != nil {
		log.Fatalln(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalln(err)
	}

	if dbCfg.MaxIdle > 0 {
		sqlDB.SetMaxIdleConns(dbCfg.MaxIdle)
	}
	if dbCfg.MaxOpen > 0 && dbCfg.MaxOpen > dbCfg.MaxIdle {
		sqlDB.SetMaxOpenConns(100)
	}
	if dbCfg.MaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(dbCfg.MaxLifetime) * time.Second)
	}
	return db
}

//DB get the default(first) *Db connection
func DB() (*gorm.DB, error) {
	if dbs, err := DBs(); nil != err {
		return nil, err
	} else {
		return dbs[0], nil
	}
}

//DBs get all database connections
func DBs() ([]*gorm.DB, error) {
	if initMultiDBSource(); 0 == multiDBSource.Size() {
		return nil, ErrNoDatabaseConfiguration
	}
	return multiDBSource.dbs, nil
}

//Close close all db connection
func CloseAllDB() error {
	if 0 == multiDBSource.Size() {
		return ErrNoDatabaseConfiguration
	}
	for _, v := range multiDBSource.dbs {
		db, err := v.DB()
		if err != nil {
			log.Fatalln(err)
		}
		if err := db.Close(); nil != err {
			return err
		}
	}
	return nil
}

//Size get db connection size
func (s MultiDBSource) Size() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return len(s.dbs)
}
