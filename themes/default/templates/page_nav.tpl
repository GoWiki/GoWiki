<ul class="nav nav-tabs">
<li {{if eq .Section "Read"}}class="active"{{end}}><a href="{{.Read}}">Read</a></li>
<li {{if eq .Section "Edit"}}class="active"{{end}}><a href="{{.Edit}}">Edit</a></li>
<li {{if eq .Section "History"}}class="active"{{end}}><a href="{{.History}}">History</a></li>
</ul>