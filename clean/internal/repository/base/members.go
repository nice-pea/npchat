package base

// Эта строка будет выдавать ошибку во время компиляции,
// если тип справа не имплементирует интерфейс слева
//var _ domain.MembersRepository = (*MembersRepository)(nil) // TODO: раскомментить.

type MembersRepository struct {
}
