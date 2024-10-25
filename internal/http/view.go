package http

type SettingsView struct {
	Form  FormView
	Error error
}

func NewSettingsView(
	form FormView,
	error error,
) SettingsView {
	return SettingsView{
		Form:  form,
		Error: error,
	}
}
