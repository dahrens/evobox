onopen = function() {
	msg = {"action": "listCreatures"}
    server.send(JSON.stringify(msg));
};

onerror = function (error) {
  	console.log('WebSocket Error ' + error);
};

onmessage = function (e) {
	try {
	    var msg = JSON.parse(e.data);
	    switch(msg.action) {
	    	case "loadCreatures":
	        	data = msg.data;
	            me.map = new Map()
			    for (var i=0; i<data.length; i++) {
			    	o = data[i];
			    	o.DT_RowId = "creature-id-" + o.Id;
			    	me.map.set(o.Id, o);
			    }
			    $("#creatures").trigger("creatures-loaded", [me.map]);
			    $("#tilemap").trigger("creatures-loaded", [me.map]);
			    break;
		    case "update":
		    	o = msg.data;
		    	o.DT_RowId = "creature-id-" + o.Id;
		        me.map.set(o.Id, o);
		        $("#creatures").trigger("creature-update", [o], me.map);
		        $("#tilemap").trigger("creature-update", [o], me.map );
		        break;
		    case "delete":
		        o = msg.data;
		        o.DT_RowId = "creature-id-" + o.Id;
		        me.map.delete(o.Id);
		        $("#creatures").trigger("creature-delete", [o], me.map );
		        $("#tilemap").trigger("creature-delete", [o], me.map );
		        break;
		    default:
		        console.log("unknown message received:");
		        console.log(msg);
		}
		
	} catch(err) {
		console.log('can not handle this request:')
	    console.log(e);
	    console.log('error:')
	    console.log(err);
	}
};

var Evobox = function() {
	me = this
	server = new WebSocket("ws://localhost:8080/connect");
	server.onopen = function() {
		onopen(me);
	}
	server.onerror = function(err) {
		onerror(err, me);
	}
	server.onmessage = function(e) {
		onmessage(e, me);
	}
	server.onclose = function () {}; // disable onclose handler first
	this.creatures = new Map()
}
Evobox.prototype.constructor = Evobox

Evobox.prototype.Start = function() {
	msg = {"action": "Start"}
    server.send(JSON.stringify(msg));
}

Evobox.prototype.Pause = function() {
	msg = {"action": "Pause"}
    server.send(JSON.stringify(msg));
}

Evobox.prototype.Reset = function() {
    server.close()
    server = new WebSocket("ws://localhost:8080/connect");
}
