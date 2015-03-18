function ReadSettings() {
	settings = {}
	$("#settings .setting").each(function(index, child){
		settings[$(this).find(':input').attr("name")] = $(this).find(':input').val()
	})
	return settings;
}

var Evobox = function(game){
	this.min_scale = 0.1
	this.max_scale = 10
	this.cur_scale = 1
	this.fragment_counter = 0

	this.creature_table = $('#creatures').DataTable({
		"searching": false,
		"aLengthMenu": [[10, 30, 50], [30, 30, 50]],
		"columns": [
			{ "data": "Name" },
			{ "data": "Gender" },
			{ "data": "Age" },
			{ "data": "Birth"},
			{ "data": "Health" },
			{ "data": "Libido" },
			{ "data": "Hunger" },
			{ "data": "X" },
			{ "data": "Y" },
		]
	});

	this.flower_table = $('#flowers').DataTable({
		"searching": false,
		"aLengthMenu": [[10, 30, 50], [30, 30, 50]],
		"columns": [
		 	{ "data": "NutritionalValue" },
		 	{ "data": "Age" },
			{ "data": "Birth"},
			{ "data": "X" },
			{ "data": "Y" },
		]
	});

	this.creatures = new Map()
	this.flowers = new Map()
	this.fragments = new Map()
};

Evobox.prototype = {
	connect: function() {
		this.server = new WebSocket("ws://localhost:8080/connect");
		this.server.evobox = this
		this.server.onopen = this.onOpen
		this.server.onerror = this.onError
		this.server.onmessage = this.onMessage
		this.server.onclose = this.onClose
		this.initialized = false
	},
	create: function() {
			// finally set up the game sceene
		this.game.time.advancedTiming = true;
		this.game.world.setBounds(0, 0, 2048, 2048);

			// change initial camera position
		this.game.camera.x = 100;
		this.game.camera.y = 100;

		back = this.game.add.sprite(0, 0, 'default', 'island.png');
		back.scale.set(this.cur_scale);

		this.flower_group = this.game.add.group();
		this.creature_group = this.game.add.group();

		this.game.input.mouse.mouseWheelCallback = mouseWheel;

		var self = this
		function mouseWheel(event) {
			if (self.game.input.mouse.wheelDelta == 1) {
				// zoom in
				self.cur_scale = self.cur_scale * 2
			} else {
				// zoom out
				self.cur_scale = self.cur_scale / 2
			}
			self.game.world.scale.set(self.cur_scale);
			self.game.world.setBounds(0, 0, self.game.world._width * self.cur_scale, self.game.world._height * self.cur_scale);
		}

		var canvas = window.document.getElementsByTagName('canvas')[0],
		prevX = 0, prevY = 0, mouseDown = false;

		canvas.addEventListener('touchstart',function(e){
			prevX = e.changedTouches[0].screenX;
			prevY = e.changedTouches[0].screenY;
		});

		canvas.addEventListener('mousedown',function(e){
			mouseDown = true;
			prevX = e.screenX;
			prevY = e.screenY;
		});

		canvas.addEventListener('touchmove',function(e){
			e.preventDefault();
			self.game.camera.x+= prevX - e.changedTouches[0].screenX;
			prevX = e.changedTouches[0].screenX;
			self.game.camera.y+= prevY - e.changedTouches[0].screenY;
			prevY = e.changedTouches[0].screenY;
		});

		canvas.addEventListener('mousemove',function(e){
			if(mouseDown){
				e.preventDefault();
				self.game.camera.x += prevX - e.screenX;
				prevX = e.screenX;
				self.game.camera.y += prevY - e.screenY;
				prevY = e.screenY;
			}
		});

		canvas.addEventListener('mouseup',function(e){
			mouseDown = false;
		});

		canvas.addEventListener('mouseleave',function(e){
			mouseDown = false;
		});

		this.connect()
	},
	render: function() {
		this.game.debug.text(this.game.time.fps + 'FPS' || '-- FPS', 2, 14, "#ffffff");
		this.game.debug.cameraInfo(this.game.camera, 2, 32);
	 	this.game.debug.inputInfo(2, 115);
	},
	onOpen: function()   {
		msg = {"Action": "Connect", "Data": ReadSettings()};
		this.send(JSON.stringify(msg));
	},
	onError: function(error) {
		console.log('WebSocket Error ' + error);
	},
	onMessage: function(raw_msg) {
		var msg = JSON.parse(raw_msg.data);
		switch(msg.Action) {
			case "load-world":
				this.evobox.loadWorld(msg.Data)
				break;
			case "add-creature":
				this.evobox.addCreature(msg.Data)
				break;
			case "add-flower":
				this.evobox.addFlower(msg.Data)
				break;
		    case "update-creature":
				this.evobox.updateCreature(msg.Data)
				break;
			case "update-flower":
				this.evobox.updateFlower(msg.Data)
				break;
			case "delete-creature":
				this.evobox.deleteCreature(msg.Data)
				break;
			case "delete-flower":
				this.evobox.deleteFlower(msg.Data)
				break;
			case "add-creature":
				this.evobox.addCreature(msg.Data)
				break;
			default:
				console.log("unknown message received:");
				console.log(msg)
		}
	},
	loadWorld: function(raw_world) {
		var self = this

		for (var i=0; i<raw_world.Plan.Fragments.length; i++) {
			fragment = raw_world.Plan.Fragments[i];
			self.addFragment(fragment)
		}
		$('#player').change(function() {
            if ($(this).prop('checked')) {
                self.pause()
            } else {
                self.start()
            }
        });
		// $('#zoomin').click(function() {
		// 	self.world.world_map.zoomIn();
		// });
		// $('#zoomout').click(function() {
		// 	self.world.world_map.zoomOut();
		// });
        $('#reset').click(function() {
            self.reset()
        });
        // $('#spawn').click(function() {
        //     self.spawn()
        // });
	},
	addFragment: function(fragment) {
		fragment.sprite = this.game.add.sprite(fragment.X, fragment.Y, fragment.Sheet, fragment.Sprite);
		fragment.Id = this.fragment_counter
		this.fragment_counter++
		this.fragments.set(fragment.Id, fragment)
	},
	updateCreature: function(raw_creature) {
		creature = this.creatures.get(raw_creature.Id);
		if (raw_creature.X != creature.X || raw_creature.Y != creature.Y) {
			this.moveCreature(creature, {x: raw_creature.X, y: raw_creature.Y});
		}
		for (var attrname in raw_creature) {
			creature[attrname] = raw_creature[attrname];
		}
		this.creatures.set(creature.Id, creature);
		this.creature_table.row('#' + creature.DT_RowId).data(creature)
	},
	updateFlower: function(raw_flower) {
		flower = this.flowers.get(raw_flower.Id);
		for (var attrname in raw_flower) {
			flower[attrname] = raw_flower[attrname];
		}
		this.flowers.set(flower.Id, flower);
		this.flower_table.row('#' + flower.DT_RowId).data(flower)
	},
	addCreature: function(creature) {
		creature.sprite = this.game.add.sprite(creature.X, creature.Y, 'default', 'creature/creature.png');

		move_frames = Phaser.Animation.generateFrameNames('creature/move/creature-move.', 0, 60, '.png', 4);
		creature.sprite.animations.add('move', move_frames, 60, true, false);
		idle_frames = ['creature/creature.png'];
		creature.sprite.animations.add('idle', idle_frames, 60, true, false);
		creature.sprite.animations.play('idle', 60);
		creature.sprite.anchor.setTo(0.5, 0.5);

		creature.DT_RowId = "creature-id-" + creature.Id;
		this.creature_table.row.add(creature).draw();
		this.creatures.set(creature.Id, creature);
		$("#creature-count").text(this.creatures.size.toString());
		this.creature_group.add(creature.sprite);
		creature.sprite.bringToTop();
	},
	addFlower: function(flower) {
		flower.sprite = this.game.add.sprite(flower.X, flower.Y, flower.Sheet, flower.Sprite);
		flower.sprite.anchor.setTo(0.5, 0.5);

		flower.DT_RowId = "flower-id-" + flower.Id;
		this.flower_table.row.add(flower).draw();
		this.flowers.set(flower.Id, flower);
		$("#flower-count").text(this.flowers.size.toString());
		this.flower_group.add(flower.sprite);
		flower.sprite.bringToTop()
	},
	moveCreature: function(creature, p) {
		if (creature.sprite.isTweening) {
			// stop active animations and tweens when new move starts...
			console.log("still tweening")
		}
		creature.sprite.animations.play('move');
		var tween = this.tweens.create(creature.sprite);
		//  creature.Speed = speed in pixels per second = the speed the sprite will move at, regardless of the distance it has to travel
		var duration = (this.game.physics.arcade.distanceToPointer(creature.sprite, p) / creature.Speed) * 1000;
		tween = tween.to({ x: p.x, y: p.y }, duration, Phaser.Easing.Linear.None, true);
		tween.onComplete.add(function() {
			creature.sprite.animations.play('idle');
		}, this);

		this.tweens.add(tween)
	},
	deleteCreature: function(raw_creature) {
		creature = this.creatures.get(raw_creature.Id);
		creature.sprite.destroy();
		this.creature_table.row('#' + creature.DT_RowId).remove().draw();
		this.creatures.delete(creature.Id);
		$("#creature-count").text(this.creatures.size.toString());
	},
	deleteFlower: function(raw_flower) {
		flower = this.flowers.get(raw_flower.Id);
		flower.sprite.destroy();
		this.flower_table.row('#' + flower.DT_RowId).remove().draw();
		this.flowers.delete(flower.Id);
		$("#flower-count").text(this.flowers.size.toString());
	},
	deleteEverything: function() {
		this.fragments.forEach(function(v,k,m){
			v.sprite.destroy();
		});
		this.creatures.forEach(function(v,k,m){
			v.sprite.destroy();
		});
		this.flowers.forEach(function(v,k,m){
			v.sprite.destroy();
		});
		this.creature_table.clear().draw()
		this.flower_table.clear().draw()
		this.creatures.clear()
		this.fragments.clear()
		this.flowers.clear()
		$("#creature-count").text(this.creatures.size.toString());
		$("#flower-count").text(this.flowers.size.toString());
	},
	start: function() {
		msg = {"Action": "Start", "Data": []}
		this.server.send(JSON.stringify(msg));
	},
	pause: function() {
		msg = {"Action": "Pause", "Data": []}
		this.server.send(JSON.stringify(msg));
	},
	reset: function() {
		this.deleteEverything();
		msg = {"Action": "Reset", "Data": ReadSettings()}
		this.server.send(JSON.stringify(msg));
		$('#player').bootstrapToggle('on')
	},
	spawn: function() {
		msg = {"Action": "Spawn", "Data": ReadSettings()}
		this.server.send(JSON.stringify(msg));
	}
}