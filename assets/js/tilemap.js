Tilemap.prototype = new PIXI.DisplayObjectContainer();
Tilemap.prototype.constructor = Tilemap;

function Tilemap(width, height, renderWidth, renderHeight) {
  PIXI.DisplayObjectContainer.call(this);
  this.interactive = true;

  this.renderWidth = renderWidth;
  this.renderHeight = renderHeight;

  this.tilesWidth = width;
  this.tilesHeight = height;

  this.tileSize = 32;
  this.zoom = 0.5;
  this.scale.x = this.scale.y = this.zoom;

  this.startLocation = { x: 0, y: 0 };

  // fill the map with tiles
  this.generateMap();

  // variables and functions for moving the map
  this.mouseoverTileCoords = [0, 0];
  this.selectedTileCoords = [0, 0];
  this.mousePressPoint = [0, 0];
  this.selectedGraphics = new PIXI.Graphics();
  this.mouseoverGraphics = new PIXI.Graphics();
  this.addChild(this.selectedGraphics);
  this.addChild(this.mouseoverGraphics);

  this.mousedown = this.touchstart = function(data) {
    if(data.getLocalPosition(this.parent).x > 0) {
      this.dragging = true;
      this.mousePressPoint[0] = data.getLocalPosition(this.parent).x - this.position.x;
      this.mousePressPoint[1] = data.getLocalPosition(this.parent).y - this.position.y;

      this.selectTile(Math.floor(this.mousePressPoint[0] / (this.tileSize * this.zoom)),
                 Math.floor(this.mousePressPoint[1] / (this.tileSize * this.zoom)));
    }
  };
  this.mouseup = this.mouseupoutside =
    this.touchend = this.touchendoutside = function(data) {
    this.dragging = false;
  };
  this.mousemove = this.touchmove = function(data)
  {
    if(this.dragging)
    {
      var position = data.getLocalPosition(this.parent);
      this.position.x = position.x - this.mousePressPoint[0];
      this.position.y = position.y - this.mousePressPoint[1];

      this.constrainTilemap();
    }
    else{
      var mouseOverPoint = [0, 0];
      mouseOverPoint[0] = data.getLocalPosition(this.parent).x - this.position.x;
      mouseOverPoint[1] = data.getLocalPosition(this.parent).y - this.position.y;

      var mouseoverTileCoords = [Math.floor(mouseOverPoint[0] / (this.tileSize * this.zoom)),
                            Math.floor(mouseOverPoint[1] / (this.tileSize * this.zoom))];
      this.mouseoverGraphics.clear();
      this.mouseoverGraphics.lineStyle(1, 0xFFFFFF, 1);
      this.mouseoverGraphics.beginFill(0x000000, 0);
      this.mouseoverGraphics.drawRect(mouseoverTileCoords[0] * this.tileSize,
                            mouseoverTileCoords[1] * this.tileSize,
                            this.tileSize - 1,
                            this.tileSize - 1);
      this.mouseoverGraphics.endFill();
    }
  };
}

Tilemap.prototype.addTile = function(x, y, terrain) {
  var tile = PIXI.Sprite.fromFrame("Tiles/grass1.png");
  tile.position.x = x * this.tileSize;
  tile.position.y = y * this.tileSize;
  tile.tileX = x;
  tile.tileY = y;
  tile.terrain = terrain;
  this.addChildAt(tile, x * this.tilesHeight + y);
}

Tilemap.prototype.clear = function() {
  this.removeChildren()
  this.generateMap()
}

Tilemap.prototype.changeTile = function(x, y, terrain) {
  this.removeChild(this.getTile(x, y));
  this.addTile(x, y, terrain);
}

Tilemap.prototype.getTile = function(x, y) {
  return this.getChildAt(x * this.tilesHeight + y);
}

Tilemap.prototype.generateMap = function() {
  // just grass bg atm...
  for(var i = 0; i < this.tilesWidth; ++i){
    var currentRow = [];
    for(var j=0; j < this.tilesHeight; j++){
      this.addTile(i, j, 0);
    }
  }

}

Tilemap.prototype.selectTile = function(x, y) {
  this.selectedTileCoords = [x, y];
  this.selectedGraphics.clear();
  this.selectedGraphics.lineStyle(2, 0xFFFF00, 1);
  this.selectedGraphics.beginFill(0x000000, 0);
  this.selectedGraphics.drawRect(this.selectedTileCoords[0] * this.tileSize,
                         this.selectedTileCoords[1] * this.tileSize,
                         this.tileSize,
                         this.tileSize);
  this.selectedGraphics.endFill();
}

Tilemap.prototype.zoomIn = function(){
  this.zoom = Math.min(this.zoom * 2, 8);
  this.scale.x = this.scale.y = this.zoom;

  this.centerOnSelectedTile();
  this.constrainTilemap();
}

Tilemap.prototype.zoomOut = function(){
  this.mouseoverGraphics.clear();

  this.zoom = Math.max(this.zoom / 2, 0.2);
  this.scale.x = this.scale.y = this.zoom;

  this.centerOnSelectedTile();
  this.constrainTilemap();
}

Tilemap.prototype.centerOnSelectedTile = function() {
  this.position.x = this.renderWidth / 2 -
    this.selectedTileCoords[0] * this.zoom * this.tileSize -
    this.tileSize * this.zoom / 2;
    this.position.y = this.renderHeight / 2 -
    this.selectedTileCoords[1] * this.zoom * this.tileSize -
    this.tileSize * this.zoom / 2;
}

Tilemap.prototype.constrainTilemap = function() {
  this.position.x = Math.max(this.position.x, -1 * this.tileSize * this.tilesWidth * this.zoom + this.renderWidth);
  this.position.x = Math.min(this.position.x, 0);
  this.position.y = Math.max(this.position.y, -1 * this.tileSize * this.tilesHeight * this.zoom + this.renderHeight);
  this.position.y = Math.min(this.position.y, 0);
}