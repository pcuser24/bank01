package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	db "github.com/user2410/simplebank/db/sqlc"
	"github.com/user2410/simplebank/storage"
	"github.com/user2410/simplebank/util"
)

func newTestServer(t *testing.T, store db.Store, fStorage storage.Storage) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomAlphanumericStr(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store, fStorage)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
