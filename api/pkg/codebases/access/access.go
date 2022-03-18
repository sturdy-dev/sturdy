package access

import (
	codebaseDB "getsturdy.com/api/pkg/codebases/db"
	"getsturdy.com/api/pkg/users"
)

func UserHasAccessToCodebase(repo codebaseDB.CodebaseUserRepository, userID users.ID, codebaseID string) bool {
	_, err := repo.GetByUserAndCodebase(userID, codebaseID)
	return err == nil
}
