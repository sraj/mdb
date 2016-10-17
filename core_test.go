package mdb_test

import (
	"testing"

	"mdb"

	. "gopkg.in/check.v1"
	"gopkg.in/mgo.v2/bson"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type (
	MdbTestSuite struct{}
)

type TestDocument struct {
	Id      bson.ObjectId `bson:"_id,omitempty"`
	DocName string        `bson:"name,omitempty"`
}

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
		Database:      "mdbtest",
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

func (s *MdbTestSuite) Test_Insert(c *C) {
	conn, err := MongoConnect()
	defer conn.Close()

	docs := make([]interface{}, 100)
	for i := 0; i < 100; i++ {
		docs[i] = TestDocument{
			DocName: "Test Collection Insert",
		}
	}

	err = conn.Insert("tests", docs...)
	c.Assert(err, IsNil)

	n, err := conn.Count("tests", mdb.Query{})
	c.Assert(err, IsNil)
	c.Assert(n, Equals, 100)

	conn.DropCollection("tests")
}

func (s *MdbTestSuite) Test_FindOne(c *C) {
	conn, err := MongoConnect()
	defer conn.Close()

	err = conn.Insert("tests", mdb.M{"x": 5, "y": 4, "type": "robot"})
	c.Assert(err, IsNil)

	var result mdb.M
	err = conn.FindOne("tests", mdb.Query{Filter: mdb.M{"type": "robot"}}, &result)
	c.Assert(err, IsNil)
	c.Assert(result["x"], Equals, 5)

	conn.DropCollection("tests")
}

func (s *MdbTestSuite) Test_FindAll(c *C) {
	conn, err := MongoConnect()
	defer conn.Close()

	docs := make([]interface{}, 100)
	for i := 0; i < 100; i++ {
		docs[i] = TestDocument{
			Id:      bson.NewObjectId(),
			DocName: "Test Collection Insert",
		}
	}
	err = conn.Insert("tests", docs...)
	c.Assert(err, IsNil)

	var result []TestDocument
	err = conn.Find("tests", mdb.Query{Limit: 10, Sort: "_id"}, &result)
	c.Assert(err, IsNil)
	c.Assert(result[0].Id, Equals, docs[0].(TestDocument).Id)
	c.Assert(len(result), Equals, 10)

	conn.DropCollection("tests")
}

func (s *MdbTestSuite) Test_Update(c *C) {
	conn, err := MongoConnect()
	defer conn.Close()

	err = conn.Insert("tests", mdb.M{"x": 42, "y": 42})
	c.Assert(err, IsNil)

	var result interface{}
	err = conn.Update("tests", mdb.Query{Filter: mdb.M{"x": 42}}, mdb.Change{
		Update:    mdb.M{"$inc": mdb.M{"y": 1}},
		ReturnNew: true,
	}, &result)
	c.Assert(err, IsNil)

	conn.DropCollection("tests")
}

func (s *MdbTestSuite) Test_RemoveOne(c *C) {
	conn, err := MongoConnect()
	defer conn.Close()

	docs := make([]interface{}, 100)
	for i := 0; i < 100; i++ {
		docs[i] = TestDocument{
			Id:      bson.NewObjectId(),
			DocName: "Test Collection Insert",
		}
	}
	err = conn.Insert("tests", docs...)
	c.Assert(err, IsNil)

	err = conn.RemoveOne("tests", mdb.Query{})
	c.Assert(err, IsNil)

	n, err := conn.Count("tests", mdb.Query{})
	c.Assert(err, IsNil)
	c.Assert(n, Equals, 99)

	conn.DropCollection("tests")
}

func (s *MdbTestSuite) Test_RemoveAll(c *C) {
	conn, err := MongoConnect()
	defer conn.Close()

	docs := make([]interface{}, 100)
	for i := 0; i < 100; i++ {
		docs[i] = TestDocument{
			Id:      bson.NewObjectId(),
			DocName: "Test Collection Insert",
		}
	}
	err = conn.Insert("tests", docs...)
	c.Assert(err, IsNil)

	err = conn.RemoveAll("tests", mdb.Query{})
	c.Assert(err, IsNil)

	n, err := conn.Count("tests", mdb.Query{})
	c.Assert(err, IsNil)
	c.Assert(n, Equals, 0)

	conn.DropCollection("tests")
}
