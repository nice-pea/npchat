package adapter

type ServiceDiscovery interface {
	NpcApiPubUrl() string
}

type ServiceDiscoveryBase struct{}

func (s *ServiceDiscoveryBase) NpcApiPubUrl() string {
	return "https://npchat.dsaime.com/api"
}
