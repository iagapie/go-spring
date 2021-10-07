package postgresdb

import (
	"fmt"
	"github.com/iagapie/go-spring/modules/sys/config"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type Database struct {
	*gorm.DB
	logger *logrus.Entry
}

func New(cfg config.DB, log *logrus.Entry) (*Database, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port, cfg.SSLMode, cfg.TimeZone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.New(log, logger.Config{
			SlowThreshold: 200 * time.Millisecond,
			LogLevel:      toGormLogLevel(log.Level),
			Colorful:      true,
		}),
	})

	if err != nil {
		return nil, err
	}

	return &Database{DB: db, logger: log}, nil
}

func (d *Database) Ping() error {
	db, err := d.DB.DB()
	if err != nil {
		d.logger.Error(err)
		return err
	}
	return db.Ping()
}

func (d *Database) Close() {
	db, err := d.DB.DB()
	if err != nil {
		d.logger.Error(err)
		return
	}

	if err = db.Close(); err != nil {
		d.logger.Error(err)
	}
}

func toGormLogLevel(lvl logrus.Level) logger.LogLevel {
	switch lvl {
	case logrus.FatalLevel, logrus.PanicLevel, logrus.ErrorLevel:
		return logger.Error
	case logrus.WarnLevel:
		return logger.Warn
	}
	return logger.Info
}
