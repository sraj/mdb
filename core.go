package mdb

import (
	"errors"
	_ "fmt"
	"time"

	"gopkg.in/mgo.v2"
)

type M map[string]interface{}

type Config struct {
	ConnectionURL string
	Database      string
}

type Connection struct {
	Config  *Config
	session *mgo.Session
}

type Query struct {
	Filter M
	Fields M
	Sort   string
	Limit  int
	Skip   int
}

var defaultLimit = 20

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

func (t *Connection) DropDatabase() error {
	return t.session.DB(t.Config.Database).DropDatabase()
}

func (t *Connection) Collection(name string) *mgo.Collection {
	return t.session.DB(t.Config.Database).C(name)
}

func (t *Connection) Insert(name string, docs ...interface{}) error {
	return t.Collection(name).Insert(docs...)
}

func (t *Connection) processQuery(iq Query, mq *mgo.Query) *mgo.Query {

	if iq.Fields != nil {
		mq = mq.Select(iq.Fields)
	}
	if iq.Sort != "" {
		mq = mq.Sort(iq.Sort)
	}
	if iq.Skip > 0 {
		mq = mq.Skip(iq.Skip)
	}
	if iq.Limit == 0 {
		iq.Limit = defaultLimit
	}
	mq = mq.Limit(iq.Limit)

	return mq
}

func (t *Connection) Find(name string, query Query, result interface{}) error {
	mq := t.Collection(name).Find(query.Filter)
	mq = t.processQuery(query, mq)
	return mq.All(result)
}

func (t *Connection) FindOne(name string, query Query, result interface{}) error {
	mq := t.Collection(name).Find(query.Filter)
	mq = t.processQuery(query, mq)
	return mq.One(result)
}
