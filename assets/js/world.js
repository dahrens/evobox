function World(renderer_width, renderer_height) {
	this.renderer_width = renderer_width || $("#tilemap").width()
	this.renderer_height = renderer_height || $(document).height() - 100;
	this.stage = null
	this.renderer = null
	this.tilemap = null
	this.table = null
	this.creatures = new Map()
	this.raw_world = null
	this.stats = new Stats();
	this.stats.setMode(0); // 0: fps, 1: ms

	// align top-left
	this.stats.domElement.style.position = 'absolute';
	this.stats.domElement.style.left = '0px';
	this.stats.domElement.style.top = '50px';

	document.body.appendChild( this.stats.domElement );
}

World.prototype = {
	init: function(raw_world) {
		self = this
		self.interval = raw_world.Speed / 1000000;
		this.raw_world = raw_world
		this.tile_width = raw_world.W
		this.tile_height = raw_world.H
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
			var world = self
			self.stats.begin();
			self.creatures.forEach(function(v,k,m) {
				var creature = v
				upd = function(c) {
					var target_x = c.X * world.tilemap.tileSize
					var target_y = c.Y * world.tilemap.tileSize
					var cur_x = c.tile.position.x
					var cur_y = c.tile.position.y
					if (target_x == cur_x && target_y == target_y) {}
					if (target_x != cur_x) {
						c.tile.position.x = (target_x > cur_x ? cur_x + 2 : cur_x - 2);
					}
					if (target_y != cur_y) {
						c.tile.position.y = (target_y > cur_y ? cur_y + 2 : cur_y - 2);
					}
				};
				upd(creature);
			});
			self.renderer.render(self.stage);
			self.stats.end();
			requestAnimFrame(self.animate);
		}

		$("#tilemap").append(self.renderer.view);
	},
	onLoaded: function() {
		var world = this.world
		world.initTilemap()
		world.initTable()

		setInterval(function() {
		    world.table.draw();
		}, 1000);

		world.stage.addChild(world.tilemap);
		world.loadCreatures()
		// begin drawing
		world.animate();
    },
    initTilemap: function() {
		this.tilemap = new Tilemap(this.tile_width, this.tile_height, this.renderer_width, this.renderer_height, this.raw_world.Plan);
		this.tilemap.position.x = 0;
		this.tilemap.zoomIn();
	},
	initTable: function() {
		this.table = $('#creatures').DataTable({
			"searching": false,
			"aLengthMenu": [[10, 30, 50], [30, 30, 50]],
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
		this.tilemap.clear(raw_world.Plan)
		this.loadCreatures()
	},
	update: function(raw_world) {
		this.raw_world = raw_world;
		var self = this;
		for (var i = 0; i < raw_world.Creatures.length; i++) {
			self.updateCreature(raw_world.Creatures[i]);
		}
	},
    loadCreatures: function() {
		for (var i=0; i<this.raw_world.Creatures.length; i++) {
			creature = this.raw_world.Creatures[i];
			this.addCreature(creature)
	    }
    },
	addCreature: function(creature) {
		creature.DT_RowId = "creature-id-" + creature.Id;
		creature.tile = PIXI.Sprite.fromFrame("Entities/Characters/Bunny.png");
		creature.tile.anchor.x = 0;
		creature.tile.anchor.y = 0;
		creature.tile.position.x = creature.X * this.tilemap.tileSize;
		creature.tile.position.y = creature.Y * this.tilemap.tileSize;
		this.tilemap.addChild(creature.tile);
		this.table.row.add(creature).draw();
		this.creatures.set(creature.Id, creature);
	},
	updateCreature: function(raw_creature) {
		creature = this.creatures.get(raw_creature.Id);
		for (var attrname in raw_creature) {
			creature[attrname] = raw_creature[attrname];
		}
		this.creatures.set(creature.Id, creature);
		this.table.row('#' + creature.DT_RowId).data(creature)
	},
	deleteCreature: function(raw_creature) {
		creature = this.creatures.get(raw_creature.Id);
		this.tilemap.removeChild(creature.tile);
		this.table.row('#' + creature.DT_RowId).remove().draw();
    },
}