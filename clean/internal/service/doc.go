// Package service содержит уровня Юзкейсов (сценариев/случаев использования бизнеса),
// пока что Юзкейсы объединены по subject из домена, и составляют
// такую структуру: service/<модель.go>/service(struct) + usecases(methods), если с ней
// станет сложно ориентироваться, можно рассмотреть такую: service/<Юзкейс.go>/usecase(struct)+Execute(method).
package service
