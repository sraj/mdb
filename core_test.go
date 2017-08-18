package mdb_test

import (
	"testing"

	"github.com/sraj/mdb"

	"gopkg.in/check.v1"
	"gopkg.in/mgo.v2/bson"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { check.TestingT(t) }

type (
	MdbTestSuite struct{}
)

type TestDocument struct {
	Id      bson.ObjectId `bson:"_id,omitempty"`
	DocName string        `bson:"name,omitempty"`
}

func (self *TestDocument) Validate() (bool, []error) {
	resultErrors := make([]error, 0, 0)
	return len(resultErrors) == 0, resultErrors
}

var _ = check.Suite(&MdbTestSuite{})

func (s *MdbTestSuite) Test_Connection(c *check.C) {
	conn, err := mdb.New(&mdb.Config{
		ConnectionURL: "mongo:27017",
		Database:      "mdbtest",
	})
	defer conn.Close()
	c.Assert(err, check.IsNil)
}

func MongoConnect() (*mdb.Connection, error) {
	return mdb.New(&mdb.Config{
		ConnectionURL: "mongo:27017",
		Database:      "mdbtest",
	})
}

func (s *MdbTestSuite) Test_RunString(c *check.C) {
	conn, err := MongoConnect()
	defer conn.Close()

	result := struct{ Ok int }{}
	err = conn.Session().Run("ping", &result)
	c.Assert(err, check.IsNil)
	c.Assert(result.Ok, check.Equals, 1)
}

func (s *MdbTestSuite) Test_Insert(c *check.C) {
	conn, err := MongoConnect()
	defer conn.Close()

	docs := make([]interface{}, 100)
	for i := 0; i < 100; i++ {
		docs[i] = TestDocument{
			DocName: "Test Collection Insert",
		}
	}

	err = conn.Insert("tests", docs...)
	c.Assert(err, check.IsNil)

	n, err := conn.Count("tests", mdb.Query{})
	c.Assert(err, check.IsNil)
	c.Assert(n, check.Equals, 100)

	conn.DropCollection("tests")
}

func (s *MdbTestSuite) Test_FindOne(c *check.C) {
	conn, err := MongoConnect()
	defer conn.Close()

	err = conn.Insert("tests", mdb.M{"x": 5, "y": 4, "type": "robot"})
	c.Assert(err, check.IsNil)

	var result mdb.M
	err = conn.FindOne("tests", mdb.Query{Filter: mdb.M{"type": "robot"}}, &result)
	c.Assert(err, check.IsNil)
	c.Assert(result["x"], check.Equals, 5)

	conn.DropCollection("tests")
}

func (s *MdbTestSuite) Test_FindAll(c *check.C) {
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
	c.Assert(err, check.IsNil)

	var result []TestDocument
	err = conn.Find("tests", mdb.Query{Limit: 10, Sort: []string{"_id"}}, &result)
	c.Assert(err, check.IsNil)
	c.Assert(result[0].Id, check.Equals, docs[0].(TestDocument).Id)
	c.Assert(len(result), check.Equals, 10)

	conn.DropCollection("tests")
}

func (s *MdbTestSuite) Test_Update(c *check.C) {
	conn, err := MongoConnect()
	defer conn.Close()

	err = conn.Insert("tests", mdb.M{"x": 42, "y": 42})
	c.Assert(err, check.IsNil)

	var result interface{}
	err = conn.Update("tests", mdb.Query{Filter: mdb.M{"x": 42}}, mdb.Change{
		Update:    mdb.M{"$inc": mdb.M{"y": 1}},
		ReturnNew: true,
	}, &result)
	c.Assert(err, check.IsNil)

	conn.DropCollection("tests")
}

func (s *MdbTestSuite) Test_RemoveOne(c *check.C) {
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
	c.Assert(err, check.IsNil)

	err = conn.RemoveOne("tests", mdb.Query{})
	c.Assert(err, check.IsNil)

	n, err := conn.Count("tests", mdb.Query{})
	c.Assert(err, check.IsNil)
	c.Assert(n, check.Equals, 99)

	conn.DropCollection("tests")
}

func (s *MdbTestSuite) Test_RemoveAll(c *check.C) {
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
	c.Assert(err, check.IsNil)

	err = conn.RemoveAll("tests", mdb.Query{})
	c.Assert(err, check.IsNil)

	n, err := conn.Count("tests", mdb.Query{})
	c.Assert(err, check.IsNil)
	c.Assert(n, check.Equals, 0)

	conn.DropCollection("tests")
}
