package http

var (
	rootRoute      = WithVirtualHost("/")
	dashboardRoute = WithVirtualHost("/dashboard")
	settingsRoute  = WithVirtualHost("/settings")
	networkRoute   = WithVirtualHost("/network")
	networkIdRoute = WithVirtualHost("/network/:id")
)
