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
  </head>
  <body>

<div class="container">
<div class="row">
<div class="col-sm-3">
</div>
<div class="col-sm-9">
{{template "page_nav.tpl" PageNav .Slug "Edit"}}
<div class="content">
<form method="POST" action="{{Route .Slug "Update"}}">
<textarea class="form-control" name="data" rows="20">{{.Content}}</textarea>
<button type="submit" class="btn btn-default">Save</button>
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