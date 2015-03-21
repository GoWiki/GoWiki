<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">



    <title>GoWiki - {{.Name}} - Edit</title>

    <link href="/static/css/bootstrap.min.css" rel="stylesheet">
    <!--[if lt IE 9]>
      <script src="https://oss.maxcdn.com/html5shiv/3.7.2/html5shiv.min.js"></script>
      <script src="https://oss.maxcdn.com/respond/1.4.2/respond.min.js"></script>
    <![endif]-->
		<style>
		body {
			background-color: #F4F4F4;
		}
		.content {
			background-color: #fff;
			padding: 20px;
			border: 1px solid #ddd;
			border-top: 0;
		}
		.content > h1, .content > h2, .content > h3 {
			margin-top: 10px;
		}
		.empty-link {
			color: red;
		}
	</style>
  </head>
  <body>
<nav class="navbar navbar-default">
  <div class="container-fluid">
    <div class="navbar-header">
      <button type="button" class="navbar-toggle collapsed" data-toggle="collapse" data-target="#bs-example-navbar-collapse-1">
        <span class="sr-only">Toggle navigation</span>
        <span class="icon-bar"></span>
        <span class="icon-bar"></span>
        <span class="icon-bar"></span>
      </button>
      <a class="navbar-brand" href="#">GoWiki</a>
    </div>

    <!-- Collect the nav links, forms, and other content for toggling -->
    <div class="collapse navbar-collapse" id="bs-example-navbar-collapse-1">
      <ul class="nav navbar-nav navbar-right">
		{{if .User.LoggedIn}}
        <li class="dropdown">
          <a href="#" class="dropdown-toggle" data-toggle="dropdown" role="button" aria-expanded="false">{{.User.Name}} <span class="caret"></span></a>
          <ul class="dropdown-menu" role="menu">
            <li><a href="{{Route "Logout"}}">Log Out</a></li>

          </ul>
        </li>
		{{else}}
		<li><a href="{{Route "LoginForm"}}">Log In</a></li>
		{{end}}
		
      </ul>
    </div><!-- /.navbar-collapse -->
  </div><!-- /.container-fluid -->
</nav>
<div class="container">
<div class="row">
<div class="col-sm-3">
{{GetContent "Logo"}}
</div>
<div class="col-sm-9">
{{GetContent "Header"}}
</div>
</div>
<div class="row">
<div class="col-sm-3">
{{GetContent "Sidebar"}}
</div>
<div class="col-sm-9">
{{template "page_nav.tpl" PageNav .Slug "Edit"}}
<div class="content">
<form method="POST" action="{{Route "Update" "page" .Slug}}">
<textarea class="form-control" name="data" rows="20">{{.Content}}</textarea>
<button type="submit" class="btn btn-primary">Save</button>
<a href="{{Route "Read" "page" .Slug}}" class="btn btn-danger">Cancel</a>
</form>
</div>
</div>

</div>
</div>
    <!-- jQuery (necessary for Bootstrap's JavaScript plugins) -->
    <script src="https://ajax.googleapis.com/ajax/libs/jquery/1.11.2/jquery.min.js"></script>
    <!-- Include all compiled plugins (below), or include individual files as needed -->
    <script src="/static/js/bootstrap.min.js"></script>
  </body>
</html>