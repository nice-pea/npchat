package sqlite

//func TestUsersRepository_Mapping(t *testing.T) {
//	t.Run("один в domain", func(t *testing.T) {
//		repoUser := user{
//			ID: uuid.NewString(),
//		}
//		domainUser := userToDomain(repoUser)
//		assert.Equal(t, repoUser.ID, domainUser.ID)
//	})
//	t.Run("один из domain", func(t *testing.T) {
//		domainUser := domain.User{
//			ID: uuid.NewString(),
//		}
//		repoUser := userFromDomain(domainUser)
//		assert.Equal(t, domainUser.ID, repoUser.ID)
//	})
//	t.Run("несколько в domain", func(t *testing.T) {
//		repoUsers := []user{
//			{ID: uuid.NewString()},
//			{ID: uuid.NewString()},
//			{ID: uuid.NewString()},
//		}
//		domainUsers := usersToDomain(repoUsers)
//		for i, repoUser := range repoUsers {
//			assert.Equal(t, repoUser.ID, domainUsers[i].ID)
//		}
//	})
//	t.Run("несколько из domain", func(t *testing.T) {
//		domainUsers := []domain.User{
//			{ID: uuid.NewString()},
//			{ID: uuid.NewString()},
//			{ID: uuid.NewString()},
//		}
//		repoUsers := usersFromDomain(domainUsers)
//		for i, domainUser := range domainUsers {
//			assert.Equal(t, domainUser.ID, repoUsers[i].ID)
//		}
//	})
//}
