package http

var (
	rootRoute                 = WithVirtualHost("/")
	dashboardRoute            = WithVirtualHost("/dashboard")
	settingsRoute             = WithVirtualHost("/settings")
	networkRoute              = WithVirtualHost("/network")
	networkIdRoute            = WithVirtualHost("/network/:id")
	networkIdConnectionsRoute = WithVirtualHost("/network/:id/connections")
	networkIdProcessesRoute   = WithVirtualHost("/network/:id/processes")
	networkIdSpeedtestRoute   = WithVirtualHost("/network/:id/speedtest")
)
