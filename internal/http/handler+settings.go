package http

import (
	"github.com/guackamolly/zero-monitor/internal/config"
	"github.com/guackamolly/zero-monitor/internal/di"
	"github.com/guackamolly/zero-monitor/internal/logging"
	"github.com/labstack/echo/v4"
)

const (
	settingsFormNodeStatsPollingId = "node-stats-polling"
	settingsFormNodeLastSeenId     = "node-last-seen"
	settingsFormAutoSavePeriodId   = "node-auto-save"
)

// Holds the current view for the settings page. If nil, that it
// means that the settings page hasn't been requested yet.
var settingsView *SettingsView

// GET /settings
func getSettingsHandler(ectx echo.Context) error {
	if settingsView != nil {
		return ectx.Render(200, "settings", settingsView)
	}

	return withServiceContainer(ectx, func(sc *di.ServiceContainer) error {
		cfg := sc.MasterConfiguration
		v := defaultSettingsView(cfg.Current())
		settingsView = &v

		return ectx.Render(200, "settings", settingsView)
	})
}

// POST /settings
func updateSettingsHandler(ectx echo.Context) error {
	if settingsView == nil {
		logging.LogError("really strange error. tried updating settings view, but settings view is not set yet.")
		return echo.ErrFailedDependency
	}

	form, err := ectx.FormParams()
	if err != nil || len(form) == 0 {
		logging.LogError("client send invalid or empty form (len: %d, err: %v)", len(form), err)
		return echo.ErrBadRequest
	}

	uf, err := settingsView.Form.Update(form)
	if err != nil {
		return ectx.Render(200, "settings", NewSettingsView(settingsView.Form, err))
	}

	v := NewSettingsView(uf, nil)
	settingsView = &v

	sp, sperr := settingsView.configurableValue(settingsFormNodeStatsPollingId)
	ls, lperr := settingsView.configurableValue(settingsFormNodeLastSeenId)
	as, aperr := settingsView.configurableValue(settingsFormAutoSavePeriodId)
	if sperr != nil || lperr != nil || aperr != nil {
		logging.LogError("really strange error. tried extract settings view form fields, but got: (sp: %v), (ls: %v), (as: %v)", sp, ls, as)
		return echo.ErrInternalServerError
	}

	return withServiceContainer(ectx, func(sc *di.ServiceContainer) error {
		cfg := sc.MasterConfiguration

		cfg.UpdateConfigurable(sp, ls, as)
		err = cfg.Save()
		if err != nil {
			logging.LogError("failed to save config after updating settings view configurable values, %v", err)
		}

		return ectx.Render(200, "settings", settingsView)
	})

}

func defaultSettingsView(cfg config.Config) SettingsView {
	return NewSettingsView(
		NewFormView(
			map[string][]FormFieldView{
				"Configuration": {
					NewRangeFormFieldView(
						settingsFormAutoSavePeriodId,
						"Auto Save (seconds)",
						"Period master node automatically saves last network changes in the configuration directory.",
						cfg.AutoSavePeriod.Value,
						cfg.AutoSavePeriod.Default,
						cfg.AutoSavePeriod.Min,
						cfg.AutoSavePeriod.Max,
					),
				},
				"Network": {
					NewRangeFormFieldView(
						settingsFormNodeStatsPollingId,
						"Stats Polling (seconds)",
						"Period nodes must wait until fetching system statistics.",
						cfg.NodeStatsPolling.Value,
						cfg.NodeStatsPolling.Default,
						cfg.NodeStatsPolling.Min,
						cfg.NodeStatsPolling.Max,
					),
					NewRangeFormFieldView(
						settingsFormNodeLastSeenId,
						"Missing node (seconds)",
						"Maximum duration that a node can go missing before reporting back to master node.",
						cfg.NodeLastSeenTimeout.Value,
						cfg.NodeLastSeenTimeout.Default,
						cfg.NodeLastSeenTimeout.Min,
						cfg.NodeLastSeenTimeout.Max,
					),
				},
			},
		),
		nil,
	)
}

func (v SettingsView) configurableValue(id string) (int, error) {
	f, err := v.Form.FieldById(id)
	if err != nil {
		return 0, err
	}

	return f.(RangeFormFieldView).Value, nil
}
