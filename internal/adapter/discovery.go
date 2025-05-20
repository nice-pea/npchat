package adapter

type ServiceDiscovery interface {
	NpcApiPubUrl() string
}

type ServiceDiscoveryBase struct {
	Debug bool
}

func (s *ServiceDiscoveryBase) NpcApiPubUrl() string {
	if s.Debug {
		return "http://localhost:8080"
	} else {
		return "https://npchat.dsaime.com/api"
	}
}
