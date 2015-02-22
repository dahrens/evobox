<html>
<head>
    <!-- DataTables CSS -->
    <link rel="stylesheet" type="text/css" href="//netdna.bootstrapcdn.com/bootstrap/3.0.3/css/bootstrap.min.css">
    <link rel="stylesheet" type="text/css" href="//cdn.datatables.net/plug-ins/f2c75b7247b/integration/bootstrap/3/dataTables.bootstrap.css">
    <link href="https://gitcdn.github.io/bootstrap-toggle/2.2.0/css/bootstrap-toggle.min.css" rel="stylesheet">

    <style type="text/css">
        .dataTables_paginate {
            float: left !important;
        }
    </style>
    <!-- jQuery -->
    <script type="text/javascript" charset="utf8" src="//code.jquery.com/jquery-1.10.2.min.js"></script>
    <!-- DataTables -->
    <script type="text/javascript" charset="utf8" src="//cdn.datatables.net/1.10.5/js/jquery.dataTables.js"></script>
    <script src="https://gitcdn.github.io/bootstrap-toggle/2.2.0/js/bootstrap-toggle.min.js"></script>
    <script type="text/javascript" charset="utf8" src="//cdn.datatables.net/plug-ins/f2c75b7247b/integration/bootstrap/3/dataTables.bootstrap.js"></script>
    <script type="text/javascript">
        $(document).ready( function () {
            table = $('#evolvers').DataTable({
                "serverSide": true,
                "ajax": "/creatures",
                "searching": false,
                "aLengthMenu": [[30, 100, 500], [30, 100, 500]],
                "columns": [
                    { "data": "Name" },
                    { "data": "Gender" },
                    { "data": "Age" },
                    { "data": "Health" },
                    { "data": "Libido" },
                    { "data": "Hunger" },
                    { "data": "X" },
                    { "data": "Y" },
                ]
            });

            function intervalTrigger() {
                return window.setInterval( function() {
                    table.ajax.reload();
                }, 500 );
            };
            var id;

            $('#player').change(function() {
                if ($(this).prop('checked')) {
                    $.ajax({url: '/pause'});
                    window.clearInterval(id);
                } else {
                    $.ajax({url: '/start'});
                    id = intervalTrigger();
                }
            })
        });
    </script>
</head>
<body>
    <div style="float: left">
    <table id="evolvers" class="display" cellspacing="0" width="800px">
        <thead>
            <tr>
                <th>Name</th>
                <th>Gender</th>
                <th>Age</th>
                <th>Health</th>
                <th>Libido</th>
                <th>Hunger</th>
                <th>X</th>
                <th>Y</th>
            </tr>
        </thead>
    </table>
    </div>
    <div style="padding: 100px;">
        <input id="player" type="checkbox" checked data-toggle="toggle" data-on="<i class='fa fa-play'></i> Play" data-off="<i class='fa fa-pause'></i> Pause">
    </div>
</body>
</html>