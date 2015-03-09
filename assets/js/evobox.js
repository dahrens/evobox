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

	this.connected = false
	this.creatures = new Map()
};

Evobox.prototype = {
	connect: function() {
		if (!this.connected) {
			this.server = new WebSocket("ws://localhost:8080/connect");
			this.server.evobox = this
			this.server.onopen = this.onOpen
			this.server.onerror = this.onError
			this.server.onmessage = this.onMessage
			this.server.onclose = this.onClose
			this.connected = true
			this.initialized = false
		}
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

	    this.animations = {
	    	'creature': {
	    		'move': []
	    	}
	    }
	    for (var i = 0; i <= 60; i++ ) {
	    	if (i < 10) {
	    		this.animations.creature.move.push('goat/move/lemmling-Cartoon-goat.000' + i + '.png')
	    	} else {
	    		this.animations.creature.move.push('goat/move/lemmling-Cartoon-goat.00' + i + '.png')
	    	}
	    }

	    this.connect()
	},
	render: function() {
		this.game.debug.cameraInfo(this.game.camera, 2, 32);
	 	this.game.debug.text(this.game.time.fps + 'FPS' || '-- FPS', 2, 14, "#ffffff");   
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
		    case "update-world":
				this.evobox.updateWorld(msg.Data)
				break;
			case "delete-creature":
				this.evobox.deleteCreature(msg.Data)
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

		for (var i=0; i<raw_world.Creatures.length; i++) {
			creature = raw_world.Creatures[i];
			self.addCreature(creature)
		}
		$('#player').change(function() {
            if ($(this).prop('checked')) {
                self.pause()
            } else {
                self.start()
            }
        });
		$('#zoomin').click(function() {
			self.world.world_map.zoomIn();
		});
		$('#zoomout').click(function() {
			self.world.world_map.zoomOut();
		});
        $('#reset').click(function() {
            self.reset()
        });
        $('#spawn').click(function() {
            self.spawn()
        });
	},
	updateWorld: function(raw_world) {
		for (var i=0; i<raw_world.Creatures.length; i++) {
			raw_creature = raw_world.Creatures[i];
			creature = this.creatures.get(raw_creature.Id);
			if (raw_creature.X != creature.X || raw_creature.Y != creature.Y) {
				this.moveCreature(creature, {x: raw_creature.X, y: raw_creature.Y});
			}
			for (var attrname in raw_creature) {
				creature[attrname] = raw_creature[attrname];
			}
			this.creatures.set(creature.Id, creature);
			this.table.row('#' + creature.DT_RowId).data(creature)
		}
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
		this.table.row.add(creature).draw();

		this.creatures.set(creature.Id, creature);
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
		this.table.row('#' + creature.DT_RowId).remove().draw();
	},
	deleteCreatures: function() {
		this.creatures.forEach(function(v,k,m){
			v.sprite.destroy();
		});
		this.table.clear().draw()
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
		this.deleteCreatures();
		msg = {"Action": "Reset", "Data": ReadSettings()}
		this.server.send(JSON.stringify(msg));
		$('#player').bootstrapToggle('on')
	},
	spawn: function() {
		msg = {"Action": "Spawn", "Data": ReadSettings()}
		this.server.send(JSON.stringify(msg));
	}
}