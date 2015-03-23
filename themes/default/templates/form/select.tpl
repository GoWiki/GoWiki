<div class="form-group">
	<label for="{{.Var}}" class="col-sm-2 control-label">{{.Name}}</label>
	<div class="col-sm-10">
		<select class="form-control" id="{{.Var}}" name="{{.Var}}">
			{{range .Options}}
			<option value="{{.Value}}"{{if eq $.Value .Value}} selected{{end}}>{{.Name}}</option>
			{{end}}
		</select>
	</div>
</div>
