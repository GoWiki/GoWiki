<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">



    <title>GoWiki - {{.Name}}</title>

    <link href="static/css/bootstrap.min.css" rel="stylesheet">
    <!--[if lt IE 9]>
      <script src="https://oss.maxcdn.com/html5shiv/3.7.2/html5shiv.min.js"></script>
      <script src="https://oss.maxcdn.com/respond/1.4.2/respond.min.js"></script>
    <![endif]-->
<style>
/* Tab Navigation */
.nav-tabs {
    margin: 0;
    padding: 0;
    border: 0;    
}
.nav-tabs > li > a {
    background: #DADADA;
    border-radius: 0;
    box-shadow: inset 0 -8px 7px -9px rgba(0,0,0,.4),-2px -2px 5px -2px rgba(0,0,0,.4);
}
.nav-tabs > li.active > a,
.nav-tabs > li.active > a:hover {
    background: #F5F5F5;
    box-shadow: inset 0 0 0 0 rgba(0,0,0,.4),-2px -3px 5px -2px rgba(0,0,0,.4);
}

/* Tab Content */
.content {
    background: #F5F5F5;
    box-shadow: 0 0 4px rgba(0,0,0,.4);
    border-radius: 0;
    text-align: center;
    padding: 10px;
}
</style>
  </head>
  <body>

<div class="container">
<div class="row">
<div class="col-sm-3">
</div>
<div class="col-sm-9">
<ul class="nav nav-tabs">
<li><a>Read</a></li>
<li><a href="edit">Edit</a></li>
</ul>
<div class="content">
{{.Content}}
</div>
</div>

</div>
</div>
    <!-- jQuery (necessary for Bootstrap's JavaScript plugins) -->
    <script src="https://ajax.googleapis.com/ajax/libs/jquery/1.11.2/jquery.min.js"></script>
    <!-- Include all compiled plugins (below), or include individual files as needed -->
    <script src="static/js/bootstrap.min.js"></script>
  </body>
</html>