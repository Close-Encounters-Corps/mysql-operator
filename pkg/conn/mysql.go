package conn

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/Close-Encounters-Corps/mysql-operator/pkg/config"
)

type Connection struct {
	Db *sql.DB
}

func MysqlConnection(cfg *config.MysqlCfg) Connection {
	db, err := sql.Open("mysql", fmt.Sprintf("%s;%s@%s/%s?%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.DefaultDb,
		cfg.UriArgs,
	))
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("failed to connect to MySQL server: %s", err.Error())
	}
	log.Println("Connected to MySQL")
	return Connection{db}
}

func (c *Connection) NewDatabase(database string, user string, password string) error {
	_, err := c.Db.Exec(fmt.Sprintf("CREATE DATABASE `%s`", database))
	if err != nil {
		return err
	}
	_, err = c.Db.Exec(fmt.Sprintf("GRANT ALL ON `%s`.* TO `%s`@'%%' IDENTIFIED BY '%s'", database, user, password))
	return err
}

func (c *Connection) DropDatabase(database string, user string) error {
	_, err := c.Db.Exec(fmt.Sprintf("DROP USER '%s'", user))
	if err != nil {
		return err
	}
	_, err = c.Db.Exec(fmt.Sprintf("DROP DATABASE `%s`", database))
	return err
}
