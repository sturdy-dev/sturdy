package access

import (
	codebaseDB "mash/pkg/codebase/db"
)

func UserHasAccessToCodebase(repo codebaseDB.CodebaseUserRepository, userID, codebaseID string) bool {
	_, err := repo.GetByUserAndCodebase(userID, codebaseID)
	return err == nil
}
