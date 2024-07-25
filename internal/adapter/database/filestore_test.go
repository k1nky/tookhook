package database

import (
	"context"
	"io"
	"os"
	"testing"

	log "github.com/k1nky/tookhook/pkg/logger"
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

func (suite *fileStoreTestSuite) TestGetRulesFileNotExist() {
	ctx := context.TODO()
	suite.fs.DSN = "file_not_exists"
	_, err := suite.fs.GetRules(ctx)
	suite.Assert().ErrorIs(err, os.ErrNotExist)
}

func (suite *fileStoreTestSuite) TestGetRulesEmpty() {
	ctx := context.TODO()
	suite.write("")
	rules, err := suite.fs.GetRules(ctx)
	suite.NoError(err)
	suite.Len(rules.Hooks, 0)
	suite.Nil(rules.Templates)
}

func (suite *fileStoreTestSuite) TestReadRulesInvalidYaml() {
	ctx := context.TODO()
	suite.write(`
		templates:
		hooks
	`)
	rules, err := suite.fs.GetRules(ctx)
	suite.Error(err)
	suite.Nil(rules)
}

func TestFileStore(t *testing.T) {
	suite.Run(t, new(fileStoreTestSuite))
}
