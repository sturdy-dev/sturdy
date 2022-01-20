package access

import (
	codebaseDB "getsturdy.com/api/pkg/codebase/db"
)

func UserHasAccessToCodebase(repo codebaseDB.CodebaseUserRepository, userID, codebaseID string) bool {
	_, err := repo.GetByUserAndCodebase(userID, codebaseID)
	return err == nil
}
