package database

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/k1nky/tookhook/internal/entity"
	log "github.com/k1nky/tookhook/internal/logger"
	"github.com/stretchr/testify/suite"
)

type fileStoreTestSuite struct {
	suite.Suite
	tmpFile *os.File
	fs      *FileStore
}

func (suite *fileStoreTestSuite) SetupTest() {
	suite.fs = NewFileStore("", &log.Blackhole{})
	f, err := os.CreateTemp("/tmp", "tookhook-fs-test")
	if err != nil {
		panic(err)
	}
	suite.tmpFile = f
	suite.fs = NewFileStore(suite.tmpFile.Name(), &log.Blackhole{})
}

func (suite *fileStoreTestSuite) TearDownTest() {
	os.Remove(suite.tmpFile.Name())
	if suite.tmpFile != nil {
		suite.tmpFile.Close()
	}
}

func (suite *fileStoreTestSuite) write(data string) {
	io.WriteString(suite.tmpFile, data)
}

func (suite *fileStoreTestSuite) TestOpenFileNotExist() {
	ctx := context.TODO()
	suite.fs.DSN = "file_not_exists"
	err := suite.fs.Open(ctx)
	suite.Assert().ErrorIs(err, os.ErrNotExist)
}

func (suite *fileStoreTestSuite) TestReadRulesEmpty() {
	ctx := context.TODO()
	suite.write("")
	err := suite.fs.ReadRules(ctx)
	suite.NoError(err)
	suite.Len(suite.fs.rules.Hooks, 0)
	suite.Nil(suite.fs.rules.Templates)
}

func (suite *fileStoreTestSuite) TestReadRulesInvalidYaml() {
	ctx := context.TODO()
	suite.write(`
		templates:
		hooks
	`)
	err := suite.fs.ReadRules(ctx)
	suite.Error(err)
	suite.Len(suite.fs.rules.Hooks, 0)
	suite.Nil(suite.fs.rules.Templates)
}

func (suite *fileStoreTestSuite) TestReadRulesValid() {
	ctx := context.TODO()
	suite.write(`
templates:
  A:
hooks:
 - income: test
   outcome:
     - type: plugin_name
       template:
         - template: T
       target: my_target
       token: my_token
`)
	err := suite.fs.ReadRules(ctx)
	suite.NoError(err)
	// TODO: assert loaded rules
}

func (suite *fileStoreTestSuite) TestReadRulesInvalid() {
	ctx := context.TODO()
	suite.write(`
templates:
  A:
hooks:
 - income: test
   outcome:
     - type:
       template:
         - template: T
       target: my_target
       token: my_token
`)
	err := suite.fs.ReadRules(ctx)
	suite.Error(err)
}

func (suite *fileStoreTestSuite) TestGetIncomeHookByName() {
	suite.fs.rules.Hooks = []entity.Hook{
		{Income: "A", Outcome: []entity.Receiver{{Type: "null"}}},
	}
	got, err := suite.fs.GetIncomeHookByName(context.TODO(), "A")
	suite.NoError(err)
	suite.NotNil(got)
	suite.Equal(suite.fs.rules.Hooks[0], *got)
}

func (suite *fileStoreTestSuite) TestGetIncomeHookByNameNotFound() {
	suite.fs.rules.Hooks = []entity.Hook{
		{Income: "A", Outcome: []entity.Receiver{{Type: "null"}}},
	}
	got, err := suite.fs.GetIncomeHookByName(context.TODO(), "B")
	suite.NoError(err)
	suite.Nil(got)
}

func TestFileStore(t *testing.T) {
	suite.Run(t, new(fileStoreTestSuite))
}
