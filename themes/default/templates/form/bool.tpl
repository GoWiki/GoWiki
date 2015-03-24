<div class="form-group">
	<div class="col-sm-10">
		<label>
			<input type="checkbox" class="form-control" id="{{.Field.Var}}" name="{{.Field.Var}}"{{if .Value}} checked{{end}}>
			{{.Field.Name}}
		</label>
	</div>
</div>