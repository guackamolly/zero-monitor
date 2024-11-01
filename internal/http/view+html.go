package http

type ElementView struct {
	ID    string
	Value string
}

func NewElementView(
	id string,
	value string,
) ElementView {
	return ElementView{
		ID:    id,
		Value: value,
	}
}
