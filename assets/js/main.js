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
var server = null;

$(document).ready( function () {

    evo = new Evobox();

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

        $('#tilemap').on("creatures-loaded", function(event, creatures) {
            creatures.forEach(function(v, k, m){
                tilemap.addCreature(v);
            });
        });

        $('#tilemap').on("creature-update", function(event, creature) {
            tilemap.updateCreature(creature);
        });

        $('#tilemap').on("creature-delete", function(event, creature) {
            console.log("delete requested");
            tilemap.removeCreature(creature);
        });

        table = $('#creatures').DataTable({
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

        table.on("creatures-loaded", function(event, creatures) {
            table.clear();
            creatures.forEach(function(v,k,m){
                table.row.add(v);
            })
            table.draw();
        });

        table.on("creature-update", function(event, creature) {
            table.row('#' + creature.DT_RowId).data(creature).draw();
        });

        table.on("creature-delete", function(event, creature) {
            table.row('#' + creature.DT_RowId).remove().draw();
        });

        $('#player').change(function() {
            if ($(this).prop('checked')) {
                evo.Pause()
            } else {
                evo.Start()
            }
        })

        $('#reset').click(function() {
            // broken
            evo.Reset()
            // $.ajax({
            //     url: '/reset',
            //     success: function() {
            //         Evobox.Init()
            //         $('#player').bootstrapToggle('on')
            //     }
            // })

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

    $("#tilemap").append(Main("", 800, 520));

});