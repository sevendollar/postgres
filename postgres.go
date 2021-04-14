// PostgreSql package with gorm implemented
package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	DefaultRedisMode = false
	DefaultAddr      = "localhost"
	DefaultPort      = "5432"
	DefaultUsername  = "postgres"
	DefaultPassword  = "postgres"
	DefaultDBName    = "postgres"
	DefaultSSLMode   = "disable"
	DefaultTimeZone  = "Asia/Shanghai"
)

type PostgresStore struct {
	db *gorm.DB

	modelValues *interface{}
}

type Options struct {
	Addr     string
	Port     int
	Username string
	Password string
	DBName   string
	SSLMode  bool
	TimeZone string
}

func NewPostgresStore(opts Options) (*PostgresStore, error) {
	var addr string
	var port string
	var username string
	var password string
	var dbname string
	var sslMode string
	var timeZone string

	// check options
	// addr
	if opts.Addr == "" {
		addr = DefaultAddr
	} else {
		// TODO: check the ip address
		addr = opts.Addr
	}
	// port
	if p := opts.Port; p == 0 {
		port = DefaultPort
	} else {
		if p < 0 || p > 65535 {
			return nil, fmt.Errorf("port should be in range from 1 to 65535")
		}
		port = strconv.Itoa(p)
	}
	// username
	if opts.Username == "" {
		username = DefaultUsername
	} else {
		username = opts.Username
	}
	// password
	if opts.Password == "" {
		password = DefaultPassword
	} else {
		password = opts.Password
	}
	// dbname
	if opts.DBName == "" {
		dbname = DefaultDBName
	} else {
		dbname = opts.DBName
	}
	// sslmode
	if !opts.SSLMode {
		sslMode = DefaultSSLMode
	} else {
		sslMode = "enable"
	}
	// timezone
	if opts.TimeZone == "" {
		timeZone = DefaultTimeZone
	} else {
		// TODO: compare from the timezone list
		timeZone = opts.TimeZone
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		addr,
		port,
		username,
		password,
		dbname,
		sslMode,
		timeZone,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Err() error {
	return s.db.Error
}

func (s *PostgresStore) RowsAffected() int {
	return int(s.db.RowsAffected)
}

// close the postgres store connections
func (s *PostgresStore) Close() error {
	sqldb, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqldb.Close()
}

// check the availability of the underlining database service
func (s *PostgresStore) Ping() error {
	sqldb, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqldb.Ping()
}

// auto migrate
func (s *PostgresStore) AutoMigrate(models ...interface{}) error {
	return s.db.AutoMigrate(models...)
}

// Create
func (s *PostgresStore) Create(ctx context.Context, value interface{}) *PostgresStore {
	s.db = s.db.WithContext(ctx).Create(value)
	return s
}

// Read all
func (s *PostgresStore) ReadAll(ctx context.Context, model interface{}) *PostgresStore {
	s.db = s.db.WithContext(ctx).Find(model)
	if s.Err() == nil && s.RowsAffected() != 0 {
		s.modelValues = new(interface{})
		*s.modelValues = model
	}
	return s
}

// Read by ID
func (s *PostgresStore) ReadByID(ctx context.Context, id int, model interface{}) *PostgresStore {
	if id < 1 {
		s.db.Error = fmt.Errorf("id can not be an negative value")
		return s
	}
	s.db = s.db.WithContext(ctx).Where("id = ?", strconv.Itoa(id)).Find(model)
	return s
}

// returns JSON
func (s *PostgresStore) Json() ([]byte, error) {
	return json.Marshal(s.modelValues)
}

// returns JSON with indentation
func (s *PostgresStore) JsonIndent() ([]byte, error) {
	return json.MarshalIndent(s.modelValues, "", "    ")
}

// func (s *PostgresStore) Update(ctx context.Context)
// func (s *PostgresStore) Delete(ctx context.Context)
