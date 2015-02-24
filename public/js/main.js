// define a few globals here
var stage = null;
var renderer = null;
var renderWidth = 800;
var renderHeight = 600;

var tilemap = null;
var menu = null;
var menuBarWidth = 120;
var table = null;

var assets_loaded = false;

$(document).ready( function () {

    function Main(tilesPath, w, h){
        // For zoomed-in pixel art, we want crisp pixels instead of fuzziness
        PIXI.scaleModes.DEFAULT = PIXI.scaleModes.NEAREST;

        // Create the stage. This will be used to add sprites (or sprite containers) to the screen.
        stage = new PIXI.Stage(0x888888);
        // Create the renderer and add it to the page.
        // (autoDetectRenderer will choose hardware accelerated if possible)
        if(w != 0 && h != 0){
          renderWidth = w;
          renderHeight = h;
        }
        renderer = PIXI.autoDetectRenderer(renderWidth, renderHeight);
        //document.body.appendChild(renderer.view);

        // Set up the asset loader for sprite images with the .json data and a callback
        loader = new PIXI.AssetLoader(["img/sprites.json"]);
        loader.onComplete = onLoaded;
        loader.load();

        return renderer.view;
    }

    function onLoaded() {
        assets_loaded = true
        // create tilemap
        tilemap = new Tilemap(32, 32);
        tilemap.position.x = menuBarWidth;

        tilemap.zoomIn();

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

        table.on( 'xhr', function ( e, settings, json ) {
          tilemap.loadCreatures(json.data)
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

        $('#reset').click(function() {
            $.ajax({url: '/reset'});
        })

        stage.addChild(tilemap);

        menu = new Menu();
        stage.addChild(menu);

        // begin drawing
        requestAnimFrame(animate);
    }

    function animate() {
        requestAnimFrame(animate);
        renderer.render(stage);
    }


    $("#demo").append(Main("", 640, 520));

});