(function () {
    var constants = {
        unitsUpdateTime: 25
    };

    var url = "ws://" + window.location.hostname + ":7102/ws";
    var eventsHub = new EventsHub();
    var transport = new WebSocketTransport(url, eventsHub);
    var render = new Render(eventsHub, constants);

    setInterval(function() {
        eventsHub.trigger('ws:send', {'command': 'get:units'});
    }, constants.unitsUpdateTime);

    //setInterval(function() {
    //    eventsHub.trigger('ws:send', {'command': 'get:player'});
    //}, 25);
})();