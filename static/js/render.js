function Render(eventsHub) {
    this.renderer = PIXI.autoDetectRenderer(window.innerWidth, window.innerHeight, {backgroundColor : 0x1099bb});
    this.stage = new PIXI.Container();
    this.stage.hitArea = new PIXI.Rectangle(0, 0, window.innerWidth, window.innerHeight);
    this.units = {};
    this.turn = 0;
    this.eventsHub = eventsHub;
    this.playerWidth = 64;
    this.playerId = null;
    document.getElementById("map").appendChild(this.renderer.view);

    this.eventsHub.on('ws:received', this.onReceived.bind(this));
    this.lastStateUpdate = new Date().getTime();

    this.stage.interactive = true;
    //this.stage.buttonMode = true;
    this.stage.on("mousemove", this.onMouseMove.bind(this));
    this.stage.on("mousedown", this.onMouseDown.bind(this));
    this.stage.on("mouseup", this.onMouseUp.bind(this));

    this.animate();
}

Render.prototype.onReceived = function(event) {
    if (event.data['Command'] === 'get:units' && event.data['Turn'] > this.turn) {
        this.addUnits(event.data['Data']);
        this.turn = event.data['Turn'];
        this.render();
    }
};

Render.prototype.onMouseDown = function(mouseData) {
    this.eventsHub.trigger('ws:send', {'command': 'set:fire:on'});
};

Render.prototype.onMouseUp = function(mouseData) {
    this.eventsHub.trigger('ws:send', {'command': 'set:fire:off'});
};

Render.prototype.onMouseMove = function(mouseData) {
    this.playerX = mouseData.data.originalEvent.clientX + this.playerWidth/2;
    this.playerY = mouseData.data.originalEvent.clientY + this.playerWidth/2;
    this.eventsHub.trigger('ws:send', {'command': 'set:player', 'X': this.playerX, 'Y': this.playerY});
};

Render.prototype.animate = function() {
    requestAnimationFrame(this.animate.bind(this));
    var timeDiff = (new Date().getTime() - this.lastStateUpdate);
    this.moveUnits(timeDiff);
    this.render();
    this.lastStateUpdate = new Date().getTime();
};

Render.prototype.moveUnits = function (step) {
    for (var i in this.units) {
        var unit = this.units[i];
        unit.position.x += unit.parentUnit[4]*step/1000;
        unit.position.y += unit.parentUnit[5]*step/1000;
    }
};

Render.prototype.getView = function() {
    return this.renderer.view;
};

Render.prototype.setPlayerCursor = function () {
    if (this.playerId) {
        document.body.style.cursor = "url('/static/resources/spaceshooter/PNG/playerShip1_blue.png'), pointer"
    } else {
        document.body.style.cursor = "wait"
    }
};

Render.prototype.addUnits = function(units) {
    var unitsIdArray = [];
    if (units['Player']) {
        this.playerId = units['Player'][0];
    } else {
        this.playerId = null;
    }
    for(var i = 0; i < units['Units'].length; i++) {
        var parent = units['Units'][i];
        if(parent) {
            var parentId = parent[0];
            unitsIdArray.push(parentId);
            if (!(parentId in this.units)) {
                var unit = this.createUnit(parent);
                if(unit) {
                    this.updateUnit(unit, parent);
                    this.addUnit(unit)
                }
            } else {
                this.updateUnit(this.units[parentId], parent);
            }
        }

    }

    this.removeDisappeared(unitsIdArray);
    this.setPlayerCursor()
};

Render.prototype.createUnit = function(parent) {
    var parentId = parent[0];
    var type = parent[1];
    var unit = null;
    switch (type) {
        case 1:
            if (parentId != this.playerId) {
                unit = BluePlayerFactory(parent);
            }
            break;
        case 10:
            unit = BlackEnemyFactory(0.5, parentId, parent);
            break;
        case 20:
            unit = BlueLaserFactory(0.5, parent);
            break;
        case 30:
            var radius = parent[6];
            unit = ExplosionFactory(radius*2);
            break;
    }
    return unit;
};

Render.prototype.updateUnit = function(unit, parent) {
    unit.parentUnit = parent;
    unit.position.x = parent[2] - unit.width/2;
    unit.position.y = parent[3] - unit.height/2;
};

Render.prototype.addUnit = function(unit) {
    this.units[unit.parentUnit[0]] = unit;
    this.stage.addChild(unit);
};

Render.prototype.removeDisappeared = function(existedIds) {
    var existedMap = {};
    for(var index = 0; index < existedIds.length; index++) {
        existedMap[existedIds[index]] = true;
    }

    for (var id in this.units) {
        if (!(id in existedMap)) {
            this.stage.removeChild(this.units[id]);
            delete this.units[id];
        }
    }
};

Render.prototype.render = function() {
    this.renderer.render(this.stage);
};