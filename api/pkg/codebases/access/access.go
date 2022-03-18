package access

import (
	"getsturdy.com/api/pkg/codebases"
	codebaseDB "getsturdy.com/api/pkg/codebases/db"
	"getsturdy.com/api/pkg/users"
)

func UserHasAccessToCodebase(repo codebaseDB.CodebaseUserRepository, userID users.ID, codebaseID codebases.ID) bool {
	_, err := repo.GetByUserAndCodebase(userID, codebaseID)
	return err == nil
}
