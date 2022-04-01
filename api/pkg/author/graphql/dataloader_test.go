package graphql

//go:generate mockgen -destination internal/mock_service/user_service_mock.go getsturdy.com/api/pkg/users/service Service

import (
	"context"
	"testing"

	"getsturdy.com/api/pkg/author/graphql/internal/mock_service"
	gqldataloader "getsturdy.com/api/pkg/graphql/dataloader"
	"getsturdy.com/api/pkg/users"

	"github.com/golang/mock/gomock"
	"github.com/graph-gophers/graphql-go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestDataloader(t *testing.T) {
	ctrl := gomock.NewController(t)
	db := mock_service.NewMockService(ctrl)

	db.EXPECT().GetByID(gomock.Any(), gomock.Eq(users.ID("user-id"))).
		Return(&users.User{ID: users.ID("user-id"), Name: "foo"}, nil).Times(1)
	db.EXPECT().GetByID(gomock.Any(), gomock.Eq(users.ID("user-id2"))).
		Return(&users.User{ID: users.ID("user-id2"), Name: "foo"}, nil).Times(0)

	root := NewResolver(db, zap.NewNop())
	ctx := gqldataloader.NewContext(context.Background())

	// Get many times, expect to be cached, and the db to only be called once
	for i := 0; i < 5; i++ {
		author, err := root.Author(ctx, "user-id")
		assert.NoError(t, err)
		assert.Equal(t, graphql.ID("user-id"), author.ID())
	}
}

func TestDataloaderRequestedOncePerCtx(t *testing.T) {
	ctrl := gomock.NewController(t)
	db := mock_service.NewMockService(ctrl)

	db.EXPECT().GetByID(gomock.Any(), gomock.Eq(users.ID("user-id"))).Do(func(_ context.Context, userID users.ID) {
		t.Log("CALLED", userID)
	}).Return(&users.User{ID: users.ID("user-id"), Name: "foo"}, nil).Times(5)

	root := NewResolver(db, zap.NewNop())

	// Get many times, expect to be cached, and the db to only be called once
	for i := 0; i < 5; i++ {
		ctx := gqldataloader.NewContext(context.Background())

		author, err := root.Author(ctx, "user-id")
		assert.NoError(t, err)
		assert.Equal(t, graphql.ID("user-id"), author.ID())

		author2, err := root.Author(ctx, "user-id")
		assert.NoError(t, err)
		assert.Equal(t, graphql.ID("user-id"), author2.ID())
	}
}
