package route

type ExternalApi struct {
	ExternalApiOption
	Path       string
	HTTPMethod string
}

type ExternalApiGroup struct {
	ExternalApiOption
	BasePath string
	Apis     map[string]*ExternalApi
}


type ExternalApiOption struct {
	Plugins map[string]*Plugin
}
