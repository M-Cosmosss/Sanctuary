package servicecenter

type ServiceCenter interface {
	Run() error
	Load() error
	Store() error
	AddServiceNode(n *ServiceNode) error
	RegisterService(s *Service) error
	RegisterServiceGroup(group *ServiceGroup) error
	GetServiceNodes(s *Service) ([]string, error)
	PrintAllRMapGroups()
}
