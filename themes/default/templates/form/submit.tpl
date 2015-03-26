<div class="form-group">
	<div class="col-sm-offset-2 col-sm-10">
		{{range .Field.Buttons}}
		{{if eq .Type "submit"}}
		<input type="submit" class="btn btn-{{.Class}}" name="{{.Var}}" value="{{.Value}}"/>
		{{else if eq .Type "link"}}
		<a href="{{.Href}}" class="btn btn-{{.Class}}">{{.Value}}</a>
		{{end}}
		{{end}}
	</div>
</div>