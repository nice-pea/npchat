package sqlite

//func TestInvitationsRepository_Mapping(t *testing.T) {
//	t.Run("один в domain", func(t *testing.T) {
//		repoInvitation := invitation{
//			ID:     uuid.NewString(),
//			ChatID: uuid.NewString(),
//			UserID: uuid.NewString(),
//		}
//		domainInvitations := invitationToDomain(repoInvitation)
//		assert.Equal(t, repoInvitation.ID, domainInvitations.ID)
//		assert.Equal(t, repoInvitation.ChatID, domainInvitations.ChatID)
//	})
//	t.Run("один из domain", func(t *testing.T) {
//		domainInvitations := domain.Invitation{
//			ID:     uuid.NewString(),
//			ChatID: uuid.NewString(),
//			UserID: uuid.NewString(),
//		}
//		repoInvitation := invitationFromDomain(domainInvitations)
//		assert.Equal(t, domainInvitations.ID, repoInvitation.ID)
//		assert.Equal(t, domainInvitations.ChatID, repoInvitation.ChatID)
//	})
//	t.Run("несколько в domain", func(t *testing.T) {
//		repoInvitations := []invitation{
//			{ID: uuid.NewString(), ChatID: uuid.NewString(), UserID: uuid.NewString()},
//			{ID: uuid.NewString(), ChatID: uuid.NewString(), UserID: uuid.NewString()},
//			{ID: uuid.NewString(), ChatID: uuid.NewString(), UserID: uuid.NewString()},
//		}
//		domainInvitations := invitationsToDomain(repoInvitations)
//		for i, repoInvitation := range repoInvitations {
//			assert.Equal(t, repoInvitation.ID, domainInvitations[i].ID)
//			assert.Equal(t, repoInvitation.ChatID, domainInvitations[i].ChatID)
//		}
//	})
//	t.Run("несколько из domain", func(t *testing.T) {
//		domainInvitations := []domain.Invitation{
//			{ID: uuid.NewString(), ChatID: uuid.NewString(), UserID: uuid.NewString()},
//			{ID: uuid.NewString(), ChatID: uuid.NewString(), UserID: uuid.NewString()},
//			{ID: uuid.NewString(), ChatID: uuid.NewString(), UserID: uuid.NewString()},
//		}
//		repoInvitations := invitationsFromDomain(domainInvitations)
//		for i, domainInvitation := range domainInvitations {
//			assert.Equal(t, domainInvitation.ID, repoInvitations[i].ID)
//			assert.Equal(t, domainInvitation.ChatID, repoInvitations[i].ChatID)
//		}
//	})
//}
