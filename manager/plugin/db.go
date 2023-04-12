package plugin

import (
	"database/sql"
	"fmt"
	tp "github.com/dsrvlabs/vatz/manager/types"
	"sync"
	"time"
	// Load sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
)

var (
	once = sync.Once{}

	db *pluginDB
)

type pluginEntry struct {
	Name           string
	Repository     string
	BinaryLocation string
	Version        string
	InstalledAt    time.Time
}

type dbWriter interface {
	AddPlugin(e pluginEntry) error
	DeletePlugin(name string) error
	UpdatePlugin() error
}

type dbReader interface {
	List() ([]pluginEntry, error)
	Get(name string) (*pluginEntry, error)
}

type pluginDB struct {
	conn *sql.Conn
	ctx  context.Context
}

func (p *pluginDB) AddPlugin(e pluginEntry) error {
	log.Info().Str("module", "db").Msg("AddPlugin")

	err := p.createPluginTable()
	if err != nil {
		log.Info().Str("module", "db").Err(err)
		return err
	}

	opts := &sql.TxOptions{Isolation: sql.LevelDefault}
	tx, err := p.conn.BeginTx(p.ctx, opts)
	if err != nil {
		log.Info().Str("module", "db").Err(err)
		return err
	}

	q := `
INSERT INTO plugin(name, repository, binary_location, version, installed_at) VALUES(?, ?, ?, ?, ?)
`

	_, err = tx.Exec(q, e.Name, e.Repository, e.BinaryLocation, e.Version, e.InstalledAt)
	if err != nil {
		log.Info().Str("module", "db").Err(err)
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Info().Str("module", "db").Err(err)
		return err
	}

	return nil
}

func (p *pluginDB) updatePluginTable(tableName string, columnName string) (bool, error) {

}

func (p *pluginDB) isFieldExist(tableName string, columnName string) (bool, error) {

	q := fmt.Sprintf("SELECT %s FROM %s LIMIT 1", columnName, tableName)

	rows, err := p.conn.QueryContext(p.ctx, q)

	if err != nil {
		log.Info().Str("module", "db").Err(err)
		return false, err
	}

	defer rows.Close()

	// Check the columns of the result set
	columns, err := rows.Columns()
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	found := false
	for _, column := range columns {
		if column == columnName {
			found = true
			break
		}
	}
	return found, nil

}

func (p *pluginDB) createPluginTable() error {
	opts := &sql.TxOptions{Isolation: sql.LevelDefault}
	tx, err := p.conn.BeginTx(p.ctx, opts)
	if err != nil {
		return err
	}

	q := `
CREATE TABLE IF NOT EXISTS plugin (
    plugin_id varchar(256) PRIMARY KEY
	name varchar(256),
    is_enabled varchar(256),
	repository varchar(256),
	binary_location varchar(256),
	version varchar(256),
	installed_at DATE)
`

	_, err = tx.Exec(q)
	if err != nil {
		log.Info().Str("module", "db").Err(err)
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Info().Str("module", "db").Err(err)
		return err
	}

	isExist, fieldError := p.isFieldExist("plugin", "is_enabled")
	if fieldError != nil {
		log.Error().Str("module", "db").Msgf("createPluginTable Error: %s", fieldError)
		return err
	}

	if !isExist {

	}

	return nil
}

func (p *pluginDB) createPluginTableByCondition(tableType tp.PluginTableCond) error {
	opts := &sql.TxOptions{Isolation: sql.LevelDefault}
	tx, err := p.conn.BeginTx(p.ctx, opts)
	if err != nil {
		return err
	}

	tableName := "plugin"
	if tableType == tp.UPDATE {
		tableName = "plugin_re"
	}

	q := `CREATE TABLE IF NOT EXISTS ` + tableName + `(
    plugin_id varchar(256) PRIMARY KEY
	name varchar(256),
    is_enabled varchar(256),
	repository varchar(256),
	binary_location varchar(256),
	version varchar(256),
	installed_at DATE)`

	_, err = tx.Exec(q)
	if err != nil {
		log.Info().Str("module", "db").Err(err)
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Info().Str("module", "db").Err(err)
		return err
	}

	return nil
}

func (p *pluginDB) DeletePlugin(name string) error {
	log.Info().Str("module", "db").Msg("DeletePlugin")

	opts := &sql.TxOptions{Isolation: sql.LevelDefault}
	tx, err := p.conn.BeginTx(p.ctx, opts)
	if err != nil {
		log.Info().Str("module", "db").Err(err)
		return err
	}

	_, err = tx.Exec("DELETE FROM plugin WHERE name = ?", name)
	if err != nil {
		log.Info().Str("module", "db").Err(err)
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Info().Str("module", "db").Err(err)
		return err
	}

	return nil
}

func (p *pluginDB) UpdatePlugin() error {
	// TODO:
	return nil
}

func (p *pluginDB) List() ([]pluginEntry, error) {
	log.Info().Str("module", "db").Msg("List Plugin")

	q := `SELECT name, repository, binary_location, version, installed_at FROM plugin`
	rows, err := p.conn.QueryContext(p.ctx, q)
	if err != nil {
		log.Info().Str("module", "db").Err(err)
		return nil, err
	}

	defer rows.Close()

	retPlugins := make([]pluginEntry, 0)

	for rows.Next() {
		e := pluginEntry{}
		err := rows.Scan(&e.Name, &e.Repository, &e.BinaryLocation, &e.Version, &e.InstalledAt)
		if err != nil {
			log.Info().Str("module", "db").Err(err)
			return nil, err
		}

		retPlugins = append(retPlugins, e)
	}

	return retPlugins, nil
}

func (p *pluginDB) Get(name string) (*pluginEntry, error) {
	log.Info().Str("module", "db").Msgf("Get %s", name)

	q := `SELECT name, repository, binary_location, version, installed_at FROM plugin WHERE name=?`
	e := pluginEntry{}
	err := p.conn.
		QueryRowContext(p.ctx, q, name).
		Scan(&e.Name, &e.Repository, &e.BinaryLocation, &e.Version, &e.InstalledAt)

	if err != nil {
		log.Info().Str("module", "db").Err(err)
		return nil, err
	}

	return &e, err
}

func newWriter(dbfile string) (dbWriter, error) {
	log.Info().Str("module", "db").Msgf("newWriter %s", dbfile)

	chanErr := make(chan error, 1)

	once.Do(func() {
		log.Info().Str("module", "db").Msg("Create DB Instance")

		ctx := context.Background()
		conn, err := getDBConnection(ctx, dbfile)
		if err != nil {
			chanErr <- err
		}

		db = &pluginDB{ctx: ctx, conn: conn}
		chanErr <- nil
	})

	var err error
	if db == nil {
		err = <-chanErr
	}

	return db, err
}

func newReader(dbfile string) (dbReader, error) {
	log.Info().Str("module", "db").Msgf("newReader %s", dbfile)

	chanErr := make(chan error, 1)

	once.Do(func() {
		log.Info().Str("module", "db").Msg("Create DB Instance")

		ctx := context.Background()
		conn, err := getDBConnection(ctx, dbfile)
		if err != nil {
			log.Error().Str("module", "db").Err(err).Msg("")
			chanErr <- err
			return
		}

		db = &pluginDB{ctx: ctx, conn: conn}
		chanErr <- nil
	})

	var err error
	if db == nil {
		err = <-chanErr
	}

	return db, err
}

func getDBConnection(ctx context.Context, dbfile string) (*sql.Conn, error) {
	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		return nil, err
	}

	conn, err := db.Conn(ctx)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
