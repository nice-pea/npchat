package adapter

// ServiceDiscovery интерфейс для обнаружения сервисов.
type ServiceDiscovery interface {
	// NpcApiPubUrl возвращает URL для публичного API NPC.
	NpcApiPubUrl() string
}

// ServiceDiscoveryBase представляет собой базовую реализацию интерфейса ServiceDiscovery.
type ServiceDiscoveryBase struct {
	Debug bool // Флаг, указывающий, включен ли режим отладки
}

// NpcApiPubUrl возвращает URL для публичного API NPC в зависимости от режима отладки.
func (s *ServiceDiscoveryBase) NpcApiPubUrl() string {
	if s.Debug {
		// Если режим отладки включен, возвращаем локальный URL
		return "http://localhost:8080"
	} else {
		// В противном случае возвращаем публичный URL
		return "https://npchat.dsaime.com/api"
	}
}
