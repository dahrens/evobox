function World(renderer_width, renderer_height) {
	this.renderer_width = renderer_width || 800
	this.renderer_height = renderer_height || 600
	this.stage = null
	this.renderer = null
	this.tilemap = null
	this.table = null
	this.creatures = new Map()
	this.raw_world = null
}

World.prototype = {
	init: function(world) {
		self = this
		this.raw_world = world
		this.tile_width = world.W
		this.tile_height = world.H
        // For zoomed-in pixel art, we want crisp pixels instead of fuzziness
        PIXI.scaleModes.DEFAULT = PIXI.scaleModes.NEAREST;

        // Create the stage. This will be used to add sprites (or sprite containers) to the renderer.
        this.stage = new PIXI.Stage(0x888888);
        // Create the renderer and add it to the page.
        // (autoDetectRenderer will choose hardware accelerated if possible)
        this.renderer = PIXI.autoDetectRenderer(this.renderer_width, this.renderer_height);
        //document.body.appendChild(renderer.view);

        // Set up the asset loader for sprite images with the .json data and a callback
        loader = new PIXI.AssetLoader(["img/sprites.json"]);
        loader.world = this
        loader.onComplete = this.onLoaded;
        loader.load();

        self.animate = function() {
	        requestAnimFrame(self.animate);
	        self.renderer.render(self.stage);
	    }

        $("#tilemap").append(self.renderer.view);
    },
    onLoaded: function() {
		world = this.world
		world.initTilemap()
		world.initTable()

		world.stage.addChild(world.tilemap);
		world.loadCreatures()
		// begin drawing
		world.animate();
    },
    initTilemap: function() {
		this.tilemap = new Tilemap(this.tile_width, this.tile_height, this.renderer_width, this.renderer_height);
		this.tilemap.position.x = 0;
		this.tilemap.zoomIn();
	},
	initTable: function() {
		this.table = $('#creatures').DataTable({
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
    },
    reload: function(raw_world) {
		var self = this
		this.creatures.forEach(function(v, k, m) {
			self.tilemap.removeChild(v.tile);
		});
		this.table.clear().draw();
		this.raw_world = raw_world
		this.loadCreatures()
    },
    loadCreatures: function() {
		for (var i=0; i<this.raw_world.Creatures.length; i++) {
			creature = this.raw_world.Creatures[i];
			creature.DT_RowId = "creature-id-" + creature.Id;
			creature.tile = PIXI.Sprite.fromFrame("Entities/Characters/Bunny.png");
			creature.tile.position.x = creature.X * this.tilemap.tileSize;
			creature.tile.position.y = creature.Y * this.tilemap.tileSize;
			this.addCreature(creature)
	    }
    },
	addCreature: function(creature) {
		this.creatures.set(creature.Id, creature);
		this.tilemap.addChild(creature.tile);
		this.table.row.add(creature).draw();
    },
	updateCreature: function(raw_creature) {
		creature = this.creatures.get(raw_creature.Id);
		for (var attrname in raw_creature) {
			creature[attrname] = raw_creature[attrname];
		}
		creature.tile.position.x = creature.X * this.tilemap.tileSize;
		creature.tile.position.y = creature.Y * this.tilemap.tileSize;
		this.creatures.set(creature.Id, creature);
		this.table.row('#' + creature.DT_RowId).data(creature).draw();
	},
	deleteCreature: function(raw_creature) {
		creature = this.creatures.get(raw_creature.Id);
		this.tilemap.removeChild(creature.tile);
		this.table.row('#' + creature.DT_RowId).remove().draw();
    },
}