<form action="{{.Action}}" method="{{.Method}}" class="form-horizontal">
{{range .Fields}}
{{.}}
{{end}}
</form>