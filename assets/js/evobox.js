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
        msg = {"action": "connect"};
		this.send(JSON.stringify(msg));
	},
	onError: function(error) {
		console.log('WebSocket Error ' + error);
	},
	onMessage: function(raw_msg) {
		try {
			var msg = JSON.parse(raw_msg.data);
			switch(msg.action) {
				case "load-world":
					this.evobox.loadWorld(msg.data)
					break;
			    case "update":
					this.evobox.world.updateCreature(msg.data)
					break;
				case "delete":
					this.evobox.world.deleteCreature(msg.data)
					break;
				default:
					console.log("unknown message received:");
				console.log(msg);
			}
		} catch(err) {
			console.log('can not handle this request:')
		    console.log(raw_msg);
		    console.log('error:')
		    console.log(err);
		}
	},
	loadWorld: function(raw_world) {
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

	        $('#reset').click(function() {
	            self.reset()
	        });
	        this.initialized = true
		} else {
			this.world.reload(raw_world);
		}
	},
	start: function() {
		msg = {"action": "Start"}
		this.server.send(JSON.stringify(msg));
	},
	pause: function() {
		msg = {"action": "Pause"}
		this.server.send(JSON.stringify(msg));
	},
	reset: function() {
		msg = {"action": "Reset"}
		this.server.send(JSON.stringify(msg));
		$('#player').bootstrapToggle('on')
	}
}
