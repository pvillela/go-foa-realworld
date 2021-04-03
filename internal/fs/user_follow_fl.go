package fs

// CommentAddSfl is the stereotype instance for the service flow that
// causes the current user start following a given other user.
type UserFollowFl struct {
	UserGetByNameDaf UserGetByNameDafT
	UserUpdateDaf    UserUpdateDafT
}

// UserFollowFlT is the function type instantiated by UserFollowFl.
type UserFollowFlT = func(username string, followedUsername string, follow bool) (PwUser, error)

func (s UserFollowFl) Make() UserFollowFlT {
	return func(username string, followedUsername string, follow bool) (PwUser, error) {
		pwUser, err := s.UserGetByNameDaf(username)
		if err != nil {
			return nil, err
		}
		user := pwUser.Entity()

		*user = user.UpdateFollowees(followedUsername, follow)

		if pwUser, err = s.UserUpdateDaf(pwUser); err != nil {
			return nil, err
		}

		return pwUser, nil
	}
}
