package main

import (
	"bytes"
	"fmt"
	"html/template"
	"strconv"
)

type StringField struct {
	Name        string
	Var         string
	Value       string
	Placeholder string
	Specialty   string
	form        *Form
}

type BoolField struct {
	Name  string
	Var   string
	Value bool
	form  *Form
}

type SelectOption struct {
	Name  string
	Value string
}

type SelectField struct {
	Name        string
	Var         string
	Value       string
	Options     []SelectOption
	Placeholder string
	form        *Form
}

type FormField interface {
	Render() template.HTML
	Parse(func(string) string)
}

type Form struct {
	fields map[string]FormField
	tpl    *template.Template
}

func NewForm(tpl *template.Template) *Form {
	return &Form{fields: make(map[string]FormField), tpl: tpl}
}

func (f *Form) NewBool(Name, Var string, Value bool) *Form {
	f.fields[Var] = &BoolField{Name: Name, Var: Var, Value: Value, form: f}
	return f
}

func (f *Form) NewString(Name, Var, Value, Placeholder string) *Form {
	f.fields[Var] = &StringField{Name: Name, Var: Var, Value: Value, Placeholder: Placeholder, form: f}
	return f
}

func (f *Form) NewPassword(Name, Var, Value, Placeholder string) *Form {
	f.fields[Var] = &StringField{Name: Name, Var: Var, Value: Value, Placeholder: Placeholder, form: f, Specialty: "password"}
	return f
}

func (f *Form) Parse(Values func(string) string) {
	for _, v := range f.fields {
		v.Parse(Values)
	}
}

func (f *Form) Render() template.HTML {
	fields := make([]template.HTML, 0, len(f.fields))
	for _, v := range f.fields {
		fields = append(fields, v.Render())
	}
	buf := &bytes.Buffer{}
	data := struct {
		Fields []template.HTML
		Action string
		Method string
	}{
		fields,
		"/Setup",
		"POST",
	}
	if err := f.tpl.ExecuteTemplate(buf, "form/form.tpl", data); err != nil {
		fmt.Println(err)
	}
	return template.HTML(buf.String())
}

func (b *BoolField) Parse(Values func(string) string) {
	val := Values(b.Var)
	b.Value, _ = strconv.ParseBool(val)
}

func (b *BoolField) Render() template.HTML {
	buf := &bytes.Buffer{}
	b.form.tpl.ExecuteTemplate(buf, "form/bool.tpl", b)
	return template.HTML(buf.String())
}

func (s *StringField) Parse(Values func(string) string) {
	s.Value = Values(s.Var)
}

func (s *StringField) Render() template.HTML {
	buf := &bytes.Buffer{}
	s.form.tpl.ExecuteTemplate(buf, "form/string.tpl", s)
	return template.HTML(buf.String())
}
