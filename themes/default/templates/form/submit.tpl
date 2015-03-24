<div class="form-group">
	<div class="col-sm-offset-2 col-sm-10">
		{{range .Field.Buttons}}
		<input type="submit" class="btn btn-{{.Class}}" name="{{.Var}}" value="{{.Value}}"/>
		{{end}}
	</div>
</div>