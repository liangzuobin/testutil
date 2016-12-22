package testutil

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

	// go sql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

const driver = "mysql"

// MySQL test db (mysql) config
type MySQL struct {
	db                       *sql.DB
	DataSource               string
	Database                 string
	User                     string
	Pwd                      string
	ScriptFile               string
	RestartMySQLFirst        bool
	DropExistsDatabaseFirst  bool
	DropDatabaseAfterTesting bool
}

func (m *MySQL) openDB() error {
	var err error
	m.db, err = sql.Open(driver, m.DataSource)
	return err
}

func (m *MySQL) dbScripts() []string {
	steps := make([]string, 0, 5)
	if m.DropExistsDatabaseFirst {
		steps = append(steps, fmt.Sprintf("DROP DATABASE IF EXISTS %s", m.Database))
	}
	return append(steps,
		fmt.Sprintf("CREATE DATABASE %s", m.Database),
		fmt.Sprintf("USE %s", m.Database),
		fmt.Sprintf("CREATE USER IF NOT EXISTS '%s'@'%%' IDENTIFIED BY '%s'", m.User, m.Pwd),
		fmt.Sprintf("GRANT ALL PRIVILEGES ON %s TO '%s'@'%%'", m.Database, m.User),
	)
}

func (m *MySQL) tableScripts() ([]string, error) {
	if len(m.ScriptFile) == 0 {
		return []string{}, nil
	}
	var name string
	if path.IsAbs(m.ScriptFile) {
		name = m.ScriptFile
	} else {
		wd, err := os.Getwd()
		if err != nil {
			return []string{}, err
		}
		name = wd + "/" + m.ScriptFile
	}
	bs, err := ioutil.ReadFile(name)
	if err != nil {
		return []string{}, err
	}
	return strings.Split(string(bs), ";"), nil
}

func (m *MySQL) execScripts(scripts []string) (string, error) {
	if len(scripts) == 0 {
		return "", errors.New("no script to exec")
	}
	for _, sql := range scripts {
		if len(sql) < 5 {
			continue
		}
		_, err := m.db.Exec(sql)
		if err != nil {
			return sql, err
		}
	}
	return "", nil
}

// Prepare create test database / create user if not exists / create tables if needed.
func (m *MySQL) Prepare() error {
	if err := m.openDB(); err != nil {
		return errors.Wrap(err, "open db failed")
	}

	if m.RestartMySQLFirst {
		err := exec.Command("mysql.server", "restart").Run()
		if err != nil {
			return errors.Wrap(err, "restart mysql failed")
		}
	}

	scripts := m.dbScripts()

	if m.ScriptFile != "" {
		tss, err := m.tableScripts()
		if err != nil {
			return errors.Wrapf(err, "read script file %s failed", m.ScriptFile)
		}
		scripts = append(scripts, tss...)
	}

	script, err := m.execScripts(scripts)
	if err != nil {
		return errors.Wrapf(err, "failed when exec: %s", script)
	}

	return nil
}

// Close close db, drop the testing database according to DropDatabaseAfterTesting.
func (m *MySQL) Close() error {
	if m.db == nil {
		return errors.New("db is nil, may never been opened")
	}
	if m.DropDatabaseAfterTesting {
		_, err := m.db.Exec("DROP DATABASE " + m.Database)
		if err != nil {
			return errors.Wrap(err, "drop database failed")
		}
	}
	return errors.Wrap(m.db.Close(), "close db failed")
}

// NewMySQL create a new mysql testing db
func NewMySQL() *MySQL {
	return new(MySQL)
}
