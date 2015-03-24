<div class="form-group">
	<label for="{{.Field.Var}}" class="col-sm-2 control-label">{{.Field.Name}}</label>
	<div class="col-sm-10">
		<input type="{{if .Field.Specialty}}{{.Field.Specialty}}{{else}}text{{end}}" class="form-control" id="{{.Field.Var}}" name="{{.Field.Var}}" placeholder="{{.Field.Placeholder}}" value="{{if ne .Field.Specialty "password"}}{{.Value}}{{end}}">
	</div>
</div>