package config

import (
	"context"
	"database/sql"
	"dc-transaction-service/log"
	"fmt"
	"net/url"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v4/stdlib"
	// "github.com/cinemaxxi-team/xxi-library/encryptor"
)

func InitDBConnection(cfg *ApplicationConfig) (map[string]*sql.DB, error) {
	// decode secret key
	// remove encryption for db password
	// decodedKey, _ := base64.StdEncoding.DecodeString(cfg.SecretKey)

	// connect to database
	dbConfig := make(map[string]*sql.DB)
	for dbName, db := range cfg.Database {
		// db.Password, _ = encryptor.Decrypt(string(decodedKey), db.Password)
		var dbCon *sql.DB
		var err error
		if db.Driver == "postgresql" {
			dbCon, err = DbConnection(db)
		} else if db.Driver == "mysql" {
			dbCon, err = DbConnectionMysql(db)
		}
		if err != nil {
			return nil, err
		}

		dbConfig[dbName] = dbCon
	}

	return dbConfig, nil
}

func DbConnection(cfg DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn(cfg))
	if err != nil {
		log.Logging.Errorf("Error %s when opening DB\n", err)
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(time.Minute * cfg.ConnMaxLifetime)

	ctx, cancelfunc := context.WithTimeout(context.Background(), cfg.Timeout*time.Second)
	defer cancelfunc()
	err = db.PingContext(ctx)
	if err != nil {
		log.Logging.Errorf("Errors %s pinging DB\n", err)
		return nil, err
	}
	log.Logging.Infof("Connected to DB %s successfully\n", cfg.DBname)
	return db, nil
}

func dsn(cfg DatabaseConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s",
		url.QueryEscape(cfg.Username), url.QueryEscape(cfg.Password),
		cfg.Hostname, cfg.DBname)
}

func DbConnectionMysql(cfg DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsnMySQL(cfg))
	if err != nil {
		log.Logging.Errorf("Error %s when opening DB\n", err)
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(time.Minute * cfg.ConnMaxLifetime)

	ctx, cancelfunc := context.WithTimeout(context.Background(), cfg.Timeout*time.Second)
	defer cancelfunc()
	err = db.PingContext(ctx)
	if err != nil {
		log.Logging.Errorf("Errors %s pinging DB\n", err)
		return nil, err
	}
	log.Logging.Infof("Connected to DB %s successfully\n", cfg.DBname)
	return db, nil
}

func dsnMySQL(cfg DatabaseConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s",
		cfg.Username, cfg.Password,
		cfg.Hostname, cfg.DBname)
}
