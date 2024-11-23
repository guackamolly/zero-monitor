package http

var (
	rootRoute                           = WithVirtualHost("/")
	dashboardRoute                      = WithVirtualHost("/dashboard")
	settingsRoute                       = WithVirtualHost("/settings")
	networkRoute                        = WithVirtualHost("/network")
	networkIdRoute                      = WithVirtualHost("/network/:id")
	networkIdConnectionsRoute           = WithVirtualHost("/network/:id/connections")
	networkIdPackagesRoute              = WithVirtualHost("/network/:id/packages")
	networkIdProcessesRoute             = WithVirtualHost("/network/:id/processes")
	networkIdSpeedtestRoute             = WithVirtualHost("/network/:id/speedtest")
	networkIdSpeedtestHistoryRoute      = WithVirtualHost("/network/:id/speedtest/history")
	networkIdSpeedtestHistoryChartRoute = WithVirtualHost("/network/:id/speedtest/history/chart")
	networkIdSpeedtestIdRoute           = WithVirtualHost("/network/:id/speedtest/:id2")
	networkPublicKeyRoute               = WithVirtualHost("/network/public-key")
	networkConnectionEndpointRoute      = WithVirtualHost("/network/connection-endpoint")
)
