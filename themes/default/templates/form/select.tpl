<div class="form-group">
	<label for="{{.Field.Var}}" class="col-sm-2 control-label">{{.Field.Name}}</label>
	<div class="col-sm-10">
		<select class="form-control" id="{{.Field.Var}}" name="{{.Field.Var}}">
			{{range .Field.Options}}
			<option value="{{.Value}}"{{if eq $.Value .Value}} selected{{end}}>{{.Name}}</option>
			{{end}}
		</select>
	</div>
</div>
