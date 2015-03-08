/*
 * Game preload
 */

var Preload = function(game) {
	console.log("Load State Initiated");
};

Preload.prototype = {

	preload: function () {
		// load resources here
		this.game.load.atlasJSONHash('default', 'sprites/default.png', 'sprites/default.json');

		// preload bar
		this.preloadBar = this.add.sprite((this.game.width) / 2.0, (this.game.height) / 2.0, 'preloaderBar');
		this.preloadBar.anchor.setTo(0.5, 0.5);
		this.load.setPreloadSprite(this.preloadBar);
		var style = { font: "32px Arial", fill: "#ababab", wordWrap: true, align: "center" };

     	this.text = this.game.add.text(this.game.world.centerX-60, this.game.world.centerY+60, "loading...", style);
	},

	create: function () {

		this.game.add.tween(this.preloadBar).to( { angle: 360 }, 2000, Phaser.Easing.Linear.None, true).loop(true);
    },

	update: function () {
		var self = this
		setTimeout(function() {
			self.state.start('Evobox');
		},100);
		
	}
}