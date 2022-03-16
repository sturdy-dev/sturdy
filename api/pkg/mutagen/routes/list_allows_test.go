//nolint:bodyclose
package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"getsturdy.com/api/pkg/auth"
	service_auth "getsturdy.com/api/pkg/auth/service"
	"getsturdy.com/api/pkg/codebase/acl"
	provider_acl "getsturdy.com/api/pkg/codebase/acl/provider"
	"getsturdy.com/api/pkg/internal/inmemory"
	"getsturdy.com/api/pkg/users"
	db_users "getsturdy.com/api/pkg/users/db"
	service_user "getsturdy.com/api/pkg/users/service"
	"getsturdy.com/api/pkg/view"
)

func TestListAllows(t *testing.T) {
	viewRepo := inmemory.NewInMemoryViewRepo()
	aclRepo := inmemory.NewInMemoryAclRepo()
	userRepo := db_users.NewMemory()

	userService := service_user.New(
		zap.NewNop(),
		userRepo,
		nil,
	)

	aclProvider := provider_acl.New(
		aclRepo,
		nil,
		nil,
	)

	authService := service_auth.New(
		nil,
		nil,
		userService,
		nil,
		aclProvider,
		nil,
	)

	type listAllowsResponse struct {
		Allows []string `json:"allows"`
	}

	route := ListAllows(
		zap.NewNop(),
		viewRepo,
		authService,
	)

	cases := []struct {
		name      string
		resources string
		expected  []string
	}{
		{
			name:      "no files",
			resources: "",
			expected:  []string{"!.git", "!.git/**/*"},
		},
		{
			name:      "subset",
			resources: `"files::pkg", "files::pkg/**/*"`,
			expected:  []string{"pkg", "pkg/**/*", "!.git", "!.git/**/*"},
		},
		{
			name:      "all",
			resources: `"files::*"`,
			expected:  []string{"*", "!.git", "!.git/**/*"},
		},
		{
			name:      "all but one",
			resources: `"files::*", "files::!README.md"`,
			expected:  []string{"*", "!README.md", "!.git", "!.git/**/*"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			userID := users.ID(uuid.NewString())
			viewID := uuid.NewString()
			codebaseID := uuid.NewString()

			usr := &users.User{
				ID: userID,
			}
			err := userRepo.Create(usr)
			assert.NoError(t, err)

			vw := view.View{
				ID:         viewID,
				UserID:     userID,
				CodebaseID: codebaseID,
			}
			err = viewRepo.Create(vw)
			assert.NoError(t, err)

			aclID := uuid.NewString()

			rawPolicy := strings.NewReplacer(
				"__USER_ID__", userID.String(),
				"__CODEBASE_ID__", codebaseID,
				"__RESOURCES__", tc.resources,
			).Replace(`{
  "groups": [
    {
      "id": "our-user",
      "members": ["users::__USER_ID__"],
    },
  ],
  "tests": [
    {
      "id": "can manage access control",
      "principal": "users::__USER_ID__",
      "allow": "write",
      "resource": "acls::__CODEBASE_ID__",
    },
  ],
  "rules": [
    {
      "id": "everyone can manage access control",
      "principals": ["groups::our-user"],
      "action": "write",
      "resources": ["acls::__CODEBASE_ID__"],
    },
    {
      "id": "user can access some files",
      "principals": ["users::__USER_ID__"],
      "action": "write",
      "resources": [ __RESOURCES__ ],
    }
  ],
}`)

			rule := acl.ACL{
				ID:         acl.ID(aclID),
				CodebaseID: codebaseID,
				CreatedAt:  time.Now(),
				RawPolicy:  rawPolicy,
			}

			err = aclRepo.Create(context.Background(), rule)
			assert.NoError(t, err)

			var res listAllowsResponse

			requestWithParams(t, userID, route, nil, &res, gin.Params{gin.Param{
				Key:   "id",
				Value: viewID,
			}})

			assert.Equal(t, tc.expected, res.Allows)
		})
	}
}

func request(t *testing.T, userID users.ID, route func(*gin.Context), request, response any) {
	requestWithParams(t, userID, route, request, response, nil)
}

func requestWithParams(t *testing.T, userID users.ID, route func(*gin.Context), request, response any, params []gin.Param) {
	res := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(res)
	c.Params = params

	data, err := json.Marshal(request)
	assert.NoError(t, err)

	c.Request, err = http.NewRequest("POST", "/", bytes.NewReader(data))
	c.Request = c.Request.WithContext(auth.NewContext(context.Background(), &auth.Subject{ID: userID.String(), Type: auth.SubjectUser}))
	assert.NoError(t, err)
	route(c)
	assert.Equal(t, http.StatusOK, res.Result().StatusCode)
	content, err := ioutil.ReadAll(res.Result().Body)
	assert.NoError(t, err)

	if len(content) > 0 {
		err = json.Unmarshal(content, response)
		assert.NoError(t, err)
	}
}
