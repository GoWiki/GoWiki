<div class="form-group">
	<label for="{{.Var}}" class="col-sm-2 control-label">{{.Name}}</label>
	<div class="col-sm-10">
		<input type="{{if .Specialty}}{{.Specialty}}{{else}}text{{end}}" class="form-control" id="{{.Var}}" name="{{.Var}}" placeholder="{{.Placeholder}}" value="{{.Value}}">
	</div>
</div>