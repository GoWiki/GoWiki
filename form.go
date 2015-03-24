package main

import (
	"bytes"
	"fmt"
	"html/template"
	"reflect"
	"strconv"
)

type StringField struct {
	Name        string
	Var         string
	Placeholder string
	Specialty   string
	form        *Form
}

type BoolField struct {
	Name string
	Var  string
	form *Form
}

type SelectOption struct {
	Name  string
	Value string
}

type SelectField struct {
	Name        string
	Var         string
	Options     []SelectOption
	Placeholder string
	form        *Form
}

type SubmitButton struct {
	Var   string
	Value string
	Class string
}

type SubmitField struct {
	Buttons []*SubmitButton
	form    *Form
}

type FormField interface {
	Render(Values interface{}) template.HTML
	Parse(Values func(string) string, Dest interface{})
}

type Form struct {
	fields []FormField
	tpl    *template.Template
}

func NewForm(tpl *template.Template) *Form {
	return &Form{tpl: tpl}
}

func (f *Form) NewBool(Name, Var string) {
	f.fields = append(f.fields, &BoolField{Name: Name, Var: Var, form: f})
}

func (f *Form) NewString(Name, Var, Placeholder string) {
	f.fields = append(f.fields, &StringField{Name: Name, Var: Var, Placeholder: Placeholder, form: f})
}

func (f *Form) NewPassword(Name, Var, Placeholder string) {
	f.fields = append(f.fields, &StringField{Name: Name, Var: Var, Placeholder: Placeholder, form: f, Specialty: "password"})
}

func (f *Form) NewButtons() *SubmitField {
	submit := &SubmitField{form: f}
	f.fields = append(f.fields, submit)
	return submit
}

func (f *Form) Parse(Values func(string) string, Dest interface{}) {
	for _, v := range f.fields {
		v.Parse(Values, Dest)
	}
}

func (f *Form) Render(Values interface{}) template.HTML {
	fields := make([]template.HTML, 0, len(f.fields))
	for _, v := range f.fields {
		fields = append(fields, v.Render(Values))
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

func (b *BoolField) Parse(Values func(string) string, Dest interface{}) {
	ValueDest := reflect.ValueOf(Dest)
	if !ValueDest.IsValid() {
		return
	}
	ValueDest = ValueDest.Elem()
	VarDest := ValueDest.FieldByName(b.Var)
	if VarDest.IsValid() {
		boolval, _ := strconv.ParseBool(Values(b.Var))
		VarDest.SetBool(boolval)
	}
}

func (b *BoolField) Render(Values interface{}) template.HTML {
	val := false
	ValuesValue := reflect.ValueOf(Values)
	if ValuesValue.IsValid() {
		VarValue := ValuesValue.FieldByName(b.Var)
		if VarValue.IsValid() {
			val = VarValue.Bool()
		}
	}
	buf := &bytes.Buffer{}
	data := struct {
		Value bool
		Field *BoolField
	}{
		val,
		b,
	}
	b.form.tpl.ExecuteTemplate(buf, "form/string.tpl", data)
	return template.HTML(buf.String())
}

func (s *StringField) Parse(Values func(string) string, Dest interface{}) {
	ValueDest := reflect.ValueOf(Dest)
	if !ValueDest.IsValid() {
		return
	}
	ValueDest = ValueDest.Elem()
	VarDest := ValueDest.FieldByName(s.Var)
	if VarDest.IsValid() {
		VarDest.SetString(Values(s.Var))
	}
}

func (s *StringField) Render(Values interface{}) template.HTML {
	val := ""
	ValuesValue := reflect.ValueOf(Values)
	if ValuesValue.IsValid() {
		VarValue := ValuesValue.FieldByName(s.Var)
		if VarValue.IsValid() {
			val = VarValue.String()
		}
	}
	buf := &bytes.Buffer{}
	data := struct {
		Value string
		Field *StringField
	}{
		val,
		s,
	}
	s.form.tpl.ExecuteTemplate(buf, "form/string.tpl", data)
	return template.HTML(buf.String())
}

func (s *SubmitField) Parse(Values func(string) string, Dest interface{}) {
	ValueDest := reflect.ValueOf(Dest)
	if !ValueDest.IsValid() {
		return
	}
	ValueDest = ValueDest.Elem()
	for _, v := range s.Buttons {
		VarDest := ValueDest.FieldByName(v.Var)
		if VarDest.IsValid() {
			if Values(v.Var) != "" {
				VarDest.SetBool(true)
			}
		}
	}
}

func (s *SubmitField) Render(Values interface{}) template.HTML {
	buf := &bytes.Buffer{}
	data := struct {
		Field *SubmitField
	}{
		s,
	}
	s.form.tpl.ExecuteTemplate(buf, "form/submit.tpl", data)
	return template.HTML(buf.String())
}

func (s *SubmitField) AddButton(Value, Var, Class string) *SubmitField {
	s.Buttons = append(s.Buttons, &SubmitButton{Value: Value, Var: Var, Class: Class})
	return s
}
