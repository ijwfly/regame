function WebSocketTransport(url, eventsHub) {
    this.url = url;
    this.eventsHub = eventsHub;
    this.connect();

    this.eventsHub.on('ws:send', this.onsend, this);
}

WebSocketTransport.prototype.connect = function () {
    this.state = 'connect';
    this.socket = new WebSocket(this.url);
    this.socket.onopen = this.onopen.bind(this);
    this.socket.onclose = this.onclose.bind(this);
    this.socket.onmessage = this.onmessage.bind(this);
    this.socket.onerror = this.onerror.bind(this);
};

WebSocketTransport.prototype.onopen = function () {
    this.state = 'ready';
    this.eventsHub.trigger('ws:ready');
};

WebSocketTransport.prototype.onclose = function (event) {
    this.state = 'close';
    if (event.wasClean) {
        console.log('Соединение закрыто чисто');
    } else {
        console.log('Обрыв соединения');
    }
    console.log('Код: ' + event.code + ' причина: ' + event.reason);

    setTimeout(this.connect.bind(this), 1000);
};

WebSocketTransport.prototype.onmessage = function (event) {
    var data = JSON.parse(event.data);
    this.eventsHub.trigger('ws:received', data);
};

WebSocketTransport.prototype.onerror = function (error) {
    console.log("Ошибка " + error.message);
};

WebSocketTransport.prototype.onsend = function (event) {
    var data = event.data;
    if ('command' in data) {
        if (this.state == 'ready') {
            this.socket.send(JSON.stringify(data));
        }
    } else {
        throw "Not a message";
    }
};