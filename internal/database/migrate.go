package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var gormMigrationScripts = []*gormigrate.Migration{
	{
		ID: "00001-CreateRepository",
		Migrate: func(tx *gorm.DB) error {
			return tx.Exec(`
CREATE TABLE IF NOT EXISTS repositories
(
    id                    bigint unsigned  not null auto_increment,
    created_at            datetime         not null default CURRENT_TIMESTAMP,
    updated_at            datetime on update CURRENT_TIMESTAMP,
    name     varchar(255)     not null,
    url                  varchar(2048)     not null,
    PRIMARY KEY (id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci;
				`).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("repositories")
		},
	},
	{
		ID: "00002-CreateScan",
		Migrate: func(tx *gorm.DB) error {
			return tx.Exec(`
CREATE TABLE IF NOT EXISTS scans
(
    id            bigint unsigned  not null auto_increment,
    created_at    datetime         not null default CURRENT_TIMESTAMP,
    updated_at    datetime on update CURRENT_TIMESTAMP,
    repository_id bigint unsigned  not null,
    status        tinyint unsigned not null,
    queued_at     datetime,
    scanning_at   datetime,
    finished_at   datetime,
    findings      json,
	note		  varchar(255),
    FOREIGN KEY (repository_id) REFERENCES repositories (id),
    PRIMARY KEY (id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci;
				`).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("scans")
		},
	},
}

// Migrate migrate schema
func (_db *DB) Migrate(ctx context.Context) error {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Disable color
		},
	)
	db := _db.GetDB().Session(&gorm.Session{Logger: newLogger})

	m := gormigrate.New(db, gormigrate.DefaultOptions, gormMigrationScripts)
	return m.Migrate()
}

// MigrateDown rollback schema
func (_db *DB) MigrateDown(ctx context.Context) error {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Disable color
		},
	)
	db := _db.GetDB().Session(&gorm.Session{Logger: newLogger})

	m := gormigrate.New(db, gormigrate.DefaultOptions, gormMigrationScripts)

	for i := range gormMigrationScripts {
		migration := gormMigrationScripts[len(gormMigrationScripts)-i-1]
		err := m.RollbackMigration(migration)
		if err != nil {
			newLogger.Error(ctx, "cannot rollback script", migration.ID)
			return err
		}
	}

	return nil
}
