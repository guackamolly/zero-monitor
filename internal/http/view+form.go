package http

import (
	"fmt"
	"net/url"
	"strconv"
)

type FormView struct {
	Groups map[string][]FormFieldView

	// represents all FormFields of this Form, present in [Groups], but mapped by their ID.
	// this field is initialized at FormView initialization.
	cacheFields map[string]FormFieldView
}

type FormFieldView interface {
	ID() string
	Label() string
	Tooltip() string
	Type() string
	IsRanged() bool
}

type MetaFormFieldView struct {
	id       string
	label    string
	tooltip  string
	ttype    string
	isRanged bool
}

type RangeFormFieldView struct {
	FormFieldView
	Value   int
	Default int
	Min     int
	Max     int
}

func (v MetaFormFieldView) ID() string {
	return v.id
}

func (v MetaFormFieldView) Label() string {
	return v.label
}

func (v MetaFormFieldView) Tooltip() string {
	return v.tooltip
}

func (v MetaFormFieldView) Type() string {
	return v.ttype
}

func (v MetaFormFieldView) IsRanged() bool {
	return v.isRanged
}

func (v RangeFormFieldView) Accepts(d int) bool {
	return d >= v.Min && d <= v.Max
}

func (v FormView) FieldById(id string) (FormFieldView, error) {
	if len(v.cacheFields) == 0 {
		return MetaFormFieldView{}, fmt.Errorf("form wasn't initialized using the New function")
	}

	f, ok := v.cacheFields[id]
	if !ok {
		return MetaFormFieldView{}, fmt.Errorf("couldn't find form field with id '%s'", id)
	}

	return f, nil
}

func (v FormView) Update(data url.Values) (FormView, error) {
	c := map[string]FormFieldView{}
	for _, fs := range v.Groups {
		for _, f := range fs {
			c[f.ID()] = f
		}
	}

	for id := range data {
		f, ok := c[id]
		if !ok {
			return FormView{}, fmt.Errorf("could not process unknown field '%s'", id)
		}

		val := data.Get(id)
		switch tf := f.(type) {
		case RangeFormFieldView:
			if len(val) == 0 {
				return FormView{}, fmt.Errorf("field '%s' cannot be empty", id)
			}

			vv, cerr := strconv.Atoi(val)
			if cerr != nil {
				return FormView{}, fmt.Errorf("field '%s' needs to be a number", id)
			}

			if !tf.Accepts(vv) {
				return FormView{}, fmt.Errorf("field '%s' is on in the expected range (>= %d and <= %d)", id, tf.Min, tf.Max)
			}

			tf.Value = vv
			c[id] = tf
			continue
		default:
			return FormView{}, fmt.Errorf("field '%s' is not supported on the server", id)
		}
	}

	uv := map[string][]FormFieldView{}
	for k, v := range v.Groups {
		fs := make([]FormFieldView, len(v))
		for i, f := range v {
			fs[i] = c[f.ID()]
		}

		uv[k] = fs
	}

	return NewFormView(uv), nil
}

func NewFormView(
	groups map[string][]FormFieldView,
) FormView {
	c := map[string]FormFieldView{}
	for _, fs := range groups {
		for _, f := range fs {
			c[f.ID()] = f
		}
	}

	return FormView{
		cacheFields: c,
		Groups:      groups,
	}
}

func NewRangeFormFieldView(
	id string,
	label string,
	tooltip string,
	value int,
	defaultValue int,
	min int,
	max int,
) RangeFormFieldView {
	return RangeFormFieldView{
		FormFieldView: MetaFormFieldView{
			id:       id,
			label:    label,
			tooltip:  tooltip,
			ttype:    "number",
			isRanged: true,
		},
		Value:   value,
		Default: defaultValue,
		Min:     min,
		Max:     max,
	}
}
