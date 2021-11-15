package servicecenter

type ServicesGroupMap map[string]*ServiceGroup
type ServicesMap map[string]*Service

type Service struct {
	Name      string
	GroupName string
	Nodes     []string
}

type ServiceNode struct {
	Url              string `json:"url" binding:"required"`
	ServiceName      string `json:"service_name" binding:"required"`
	ServiceGroupName string `json:"service_group_name" binding:"required"`
}

type ServiceGroup struct {
	Name     string
	Services ServicesMap
}
