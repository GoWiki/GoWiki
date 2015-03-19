<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">



    <title>GoWiki - {{.Name}}</title>

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
		.content.no-nav {
			border-top: 1px solid #ddd;
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
<div class="content no-nav">
<form action="{{Route "Login"}}" method="POST">
<input type="text" name="username" placeholder="Username" class="form-control">
<input type="password" name="password" placeholder="Password" class="form-control">
<input type="submit" class="btn btn-primary">
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