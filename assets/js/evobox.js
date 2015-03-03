function ReadSettings() {
	settings = {}
	$("#settings .setting").each(function(index, child){
		settings[$(this).find(':input').attr("name")] = $(this).find(':input').val()
	})
	return settings;
}

function Evobox() {
	this.server = new WebSocket("ws://localhost:8080/connect");
	this.server.evobox = this
	this.server.onopen = this.onOpen
	this.server.onerror = this.onError
	this.server.onmessage = this.onMessage
	this.server.onclose = this.onClose
	this.world = new World()
	this.initialized = false
};

Evobox.prototype = {
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
				//this.evobox.world.updateCreature(msg.Data)
				this.evobox.world.update(msg.Data)
				break;
			case "delete-creature":
				this.evobox.world.deleteCreature(msg.Data)
				break;
			case "add-creature":
				this.evobox.world.addCreature(msg.Data)
				break;
			default:
				console.log("unknown message received:");
				console.log(msg)
		}
	},
	loadWorld: function(raw_world) {
		console.log(raw_world)
		var self = this
		if (self.initialized === false) {
			this.world.init(raw_world)
			$('#player').change(function() {
	            if ($(this).prop('checked')) {
	                self.pause()
	            } else {
	                self.start()
	            }
	        });
			$('#zoomin').click(function() {
				self.world.tilemap.zoomIn();
			});
			$('#zoomout').click(function() {
				self.world.tilemap.zoomOut();
			});
	        $('#reset').click(function() {
	            self.reset()
	        });
	        $('#spawn').click(function() {
	            self.spawn()
	        });
	        $('#tilemap').on('mousewheel', function(event) {
				if (event.deltaY === -1) {
					self.world.tilemap.zoomOut();
				} else {
					self.world.tilemap.zoomIn();
				}
			});
	        this.initialized = true
		} else {
			this.world.reload(raw_world);
		}
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
		msg = {"Action": "Reset", "Data": ReadSettings()}
		this.server.send(JSON.stringify(msg));
		$('#player').bootstrapToggle('on')
	},
	spawn: function() {
		msg = {"Action": "Spawn", "Data": ReadSettings()}
		this.server.send(JSON.stringify(msg));
	}
}