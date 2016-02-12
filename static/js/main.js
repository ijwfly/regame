(function () {
    var url = "ws://" + window.location.hostname + ":7102/ws";
    var eventsHub = new EventsHub();
    var transport = new WebSocketTransport(url, eventsHub);
    var render = new Render(eventsHub);

    eventsHub.on('ws:ready', function(event) {
        eventsHub.trigger('ws:send', {'command': 'ping'});
    });

    setInterval(function() {
        eventsHub.trigger('ws:send', {'command': 'get:units'});
    }, 50)
})();