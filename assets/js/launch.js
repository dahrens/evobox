$(document).ready( function () {

    $('li.dropdown.mega-dropdown a').on('click', function (event) {
	    $(this).parent().toggleClass("open");
	});

	$('body').on('click', function (e) {
	    if (!$('li.dropdown.mega-dropdown').is(e.target) && $('li.dropdown.mega-dropdown').has(e.target).length === 0 && $('.open').has(e.target).length === 0) {
	        $('li.dropdown.mega-dropdown').removeClass('open');
	    }
	});

    $(".toggle").addClass("navbar-btn");

    width = $("#sceene").width()
	height = $(document).height() - 100;

	var game = new Phaser.Game(width, height, Phaser.Auto, 'sceene');
	game.state.add('Boot', Boot);
	game.state.add('Preload', Preload);
	game.state.add('Evobox', Evobox);
	game.state.start('Boot');
});