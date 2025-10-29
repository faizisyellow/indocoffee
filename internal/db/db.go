package db

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
)

type TransFnc func(db *sql.DB, ctx context.Context, operation func(*sql.Tx) error) error

func New(dsn string, maxOpenConn, maxIdleConn int, maxIdleTime, maxLifeTime string, usr, passwd, dbname string) (*sql.DB, error) {
	rootCertPool := x509.NewCertPool()
	pem, err := os.ReadFile("ca.pem")
	if err != nil {
		log.Fatal("Failed to read CA file:", err)
	}

	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		log.Fatal("Failed to append CA cert")
	}

	err = mysql.RegisterTLSConfig("aiven", &tls.Config{
		RootCAs: rootCertPool,
	})
	if err != nil {
		log.Fatal("Failed to register TLS config:", err)
	}

	cfg := mysql.Config{
		User:                 usr,
		Passwd:               passwd,
		Net:                  "tcp",
		Addr:                 dsn,
		DBName:               dbname,
		TLSConfig:            "aiven",
		ParseTime:            true,
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("error open db=%v", err)
	}

	db.SetMaxIdleConns(maxIdleConn)
	db.SetMaxOpenConns(maxOpenConn)

	idleTime, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, fmt.Errorf("error parse time=%v", err)
	}

	lifeTime, err := time.ParseDuration(maxLifeTime)
	if err != nil {
		return nil, fmt.Errorf("error parse time=%v", err)
	}

	db.SetConnMaxLifetime(lifeTime)
	db.SetConnMaxIdleTime(idleTime)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("error connect to db=%v", err)
	}

	return db, nil
}

type TransactionDB struct {
	Db *sql.DB
}

func (t *TransactionDB) WithTx(ctx context.Context, fnc func(tx *sql.Tx) error) error {

	tx, err := t.Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fnc(tx); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}
