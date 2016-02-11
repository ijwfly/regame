var blackEnemyImage = '/static/resources/spaceshooter/PNG/Enemies/enemyBlack1.png';
var blackEnemyImageTexture = PIXI.Texture.fromImage(blackEnemyImage);
var blackEnemyTexture = new PIXI.Texture(blackEnemyImageTexture, new PIXI.Rectangle(0, 0, 93, 84));

function scaleTexture(texture, scale) {
    if (scale === undefined) {
        scale = 1;
    }
    texture.width = texture.width*scale;
    texture.height = texture.height*scale;
}

function BlackEnemyFactory(scale, zIndex, parent) {
    var spaceShip = new PIXI.Sprite(blackEnemyTexture);
    scaleTexture(spaceShip, scale);
    spaceShip.zIndex = zIndex;
    spaceShip.parentUnit = parent;
    return spaceShip
}

var blueLaserImage = '/static/resources/spaceshooter/PNG/Lasers/laserBlue08.png';
var blueLaserImageTexture = PIXI.Texture.fromImage(blueLaserImage);
var blueLaserTexture = new PIXI.Texture(blueLaserImageTexture, new PIXI.Rectangle(0, 0, 48, 46));

function BlueLaserFactory(scale, parent) {
    var laser = new PIXI.Sprite(blueLaserTexture);
    scaleTexture(laser, scale);
    laser.zIndex = 100000;
    laser.parentUnit = parent;
    return laser
}

var meteorImage = '/static/resources/spaceshooter/PNG/Meteors/meteorBrown_big1.png';
var meteorImageTexture = PIXI.Texture.fromImage(meteorImage);
var meteorTexture = new PIXI.Texture(meteorImageTexture, new PIXI.Rectangle(0, 0, 101, 84));

function MeteorFactory(scale, parent) {
    var meteor = new PIXI.Sprite(meteorTexture);
    scaleTexture(meteor, scale);
    meteor.zIndex = 10000;
    meteor.parentUnit = parent;
    return meteor
}

PIXI.loader
    .add('/static/resources/explosion/mc.json')
    .load(onAssetsLoaded);
var frames = [];

function onAssetsLoaded() {
    for (var i = 1; i < 28; i++) {
        frames.push(PIXI.Texture.fromFrame('Explosion_Sequence_A ' + i + '.png'));
    }
}

function ExplosionFactory(scale) {
    var movie = new PIXI.extras.MovieClip(frames);
    movie.animationSpeed = 0.5;
    scaleTexture(movie, scale);
    movie.play();
    return movie;
}