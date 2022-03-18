package routes

import (
	"context"

	"getsturdy.com/api/pkg/author"
	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/codebases/db"
	service_user "getsturdy.com/api/pkg/users/service"
)

func membersAsAuthors(ctx context.Context, codebaseUserRepo db.CodebaseUserRepository, userService service_user.Service, codebaseID codebases.ID) ([]author.Author, error) {
	// Get members
	members, err := codebaseUserRepo.GetByCodebase(codebaseID)
	if err != nil {
		return nil, err
	}

	var memberAuthors []author.Author
	for _, m := range members {
		userAuthor, err := userService.GetAsAuthor(ctx, m.UserID)
		if err != nil {
			return nil, err
		}
		memberAuthors = append(memberAuthors, *userAuthor)
	}

	return memberAuthors, nil
}
