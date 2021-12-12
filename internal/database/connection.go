package database

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	_mysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// NewFromEnv sets up the database connections using the configuration in the
// process's environment variables. This should be called just once per server
// instance.
func NewFromEnv(ctx context.Context, cfg *Config) (*DB, error) {
	db, err := gorm.Open(_mysql.Open(dbToMysqlDSN(cfg)), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %v", err)
	}
	sqlDb, err := db.DB()
	sqlDb.SetMaxOpenConns(cfg.PoolMaxConnections)
	sqlDb.SetMaxIdleConns(cfg.PoolMaxIdleConnections)
	_db := &DB{db: db}
	return _db, nil
}

// dbToMysqlDSN builds a connection string suitable for the mysql driver, using
// the values of vars.
func dbToMysqlDSN(cfg *Config) string {
	mySqlConfig := mysql.NewConfig()
	mySqlConfig.Addr = cfg.Address
	mySqlConfig.Passwd = cfg.Password
	mySqlConfig.Net = cfg.Protocol
	mySqlConfig.User = cfg.User
	mySqlConfig.Timeout = time.Duration(cfg.ConnectionTimeout) * time.Second
	mySqlConfig.DBName = cfg.Name
	mySqlConfig.ParseTime = true

	if cfg.SSLMode != "disable" {
		tlsConfigName := "custom"
		rootCertPool := x509.NewCertPool()
		pem, err := ioutil.ReadFile(cfg.SSLRootCertPath)
		if err != nil {
			panic(err)
		}
		if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
			panic("Failed to append PEM.")
		}
		clientCert := make([]tls.Certificate, 0, 1)
		certs, err := tls.LoadX509KeyPair(cfg.SSLCertPath, cfg.SSLKeyPath)
		if err != nil {
			panic(err)
		}
		clientCert = append(clientCert, certs)
		err = mysql.RegisterTLSConfig(tlsConfigName, &tls.Config{
			RootCAs:      rootCertPool,
			Certificates: clientCert,
		})
		if err != nil {
			panic(err)
		}
		mySqlConfig.TLSConfig = tlsConfigName
	}
	return mySqlConfig.FormatDSN()
}

// dbToDSN builds a connection string suitable for the pgx Postgres driver, using
// the values of vars.
func dbToDSN(cfg *Config) string {
	vals := dbValues(cfg)
	p := make([]string, 0, len(vals))
	for k, v := range vals {
		p = append(p, fmt.Sprintf("%s=%s", k, v))
	}
	fmt.Println(strings.Join(p, " "))
	return strings.Join(p, " ")
}

func setIfNotEmpty(m map[string]string, key, val string) {
	if val != "" {
		m[key] = val
	}
}

func setIfPositive(m map[string]string, key string, val int) {
	if val > 0 {
		m[key] = fmt.Sprintf("%d", val)
	}
}

func setIfPositiveDuration(m map[string]string, key string, d time.Duration) {
	if d > 0 {
		m[key] = d.String()
	}
}

func dbValues(cfg *Config) map[string]string {
	p := map[string]string{}
	hostAndPort := strings.Split(cfg.Address, ":")
	setIfNotEmpty(p, "dbname", cfg.Name)
	setIfNotEmpty(p, "user", cfg.User)
	setIfNotEmpty(p, "host", hostAndPort[0])
	setIfNotEmpty(p, "port", hostAndPort[1])
	setIfNotEmpty(p, "sslmode", cfg.SSLMode)
	setIfPositive(p, "connect_timeout", cfg.ConnectionTimeout)
	setIfNotEmpty(p, "password", cfg.Password)
	setIfNotEmpty(p, "sslcert", cfg.SSLCertPath)
	setIfNotEmpty(p, "sslkey", cfg.SSLKeyPath)
	setIfNotEmpty(p, "sslrootcert", cfg.SSLRootCertPath)
	//setIfNotEmpty(p, "pool_max_conns", fmt.Sprint(cfg.PoolMaxConnections))
	return p
}
