package mdb

import (
	"errors"
	"time"

	"gopkg.in/mgo.v2"
)

type (
	Config struct {
		ConnectionURL string
		Database      string
	}

	Connection struct {
		Config  *Config
		session *mgo.Session
	}
)

func NewConnection(config *Config) (*Connection, error) {
	conn := &Connection{
		Config: config,
	}
	err := conn.Connect()
	return conn, err
}

func (t *Connection) Connect() (err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else if e, ok := r.(string); ok {
				err = errors.New(e)
			} else {
				err = errors.New("mdb: error connecting mongo!")
			}
		}
	}()

	session, err := mgo.DialWithTimeout(t.Config.ConnectionURL, 5*time.Second)
	if err != nil {
		return err
	}

	t.session = session
	t.session.SetMode(mgo.Monotonic, true)
	t.session.SetSafe(&mgo.Safe{})

	return nil
}

func (t *Connection) Session() *mgo.Session {
	return t.session
}

func (t *Connection) Close() {
	if t.session != nil {
		t.session.Close()
	}
}
