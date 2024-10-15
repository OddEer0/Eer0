package eapp

type (
	info struct {
		name, version, description string
	}

	app struct {
		info           info
		configs        *Configs
		err            error
		beforeHandlers map[string]BeforeHandler
		afterHandlers  map[string]AfterHandler
		jobs           map[string]Job
	}
)
