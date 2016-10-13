package mdb_test

import (
	"testing"

	"mdb"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type (
	MdbTestSuite struct{}
)

var _ = Suite(&MdbTestSuite{})

func (s *MdbTestSuite) Test_Connection(c *C) {
	conn, err := mdb.NewConnection(&mdb.Config{
		ConnectionURL: "localhost:37017",
		Database:      "mdbtest",
	})
	defer conn.Close()
	c.Assert(err, IsNil)
}

func (s *MdbTestSuite) Test_Connect(c *C) {
  conn := &mdb.Connection{
		Config: &mdb.Config{
  		ConnectionURL: "localhost:37017",
  		Database:      "mdbtest",
  	},
	}
	err := conn.Connect()
	defer conn.Close()
	c.Assert(err, IsNil)
}

func MongoConnect() (*mdb.Connection, error) {
	return mdb.NewConnection(&mdb.Config{
		ConnectionURL: "localhost:37017",
		Database:      "mongotest",
	})
}

func (s *MdbTestSuite) Test_RunString(c *C) {
  conn, err := MongoConnect()
	defer conn.Close()

	result := struct{ Ok int }{}
	err = conn.Session().Run("ping", &result)
	c.Assert(err, IsNil)
	c.Assert(result.Ok, Equals, 1)
}
