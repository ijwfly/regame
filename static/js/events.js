function Event(name, data){
    this.name = name;
    this.data = data;
}

function EventsHub(){
    this.events = {};

    this.on = function(event, callback, context) {
        if (!(callback instanceof Function)){
            throw "Not a function";
        }
        if (context) {
            callback = callback.bind(context);
        }

        var listeners = this.events[event];
        if (listeners && listeners.length) {
            listeners.push(callback);
        } else {
            this.events[event] = [callback];
        }
    };

    this.stop = function(event, callback) {

    };

    this.trigger = function(name, data) {
        var event = new Event(name, data);
        var listeners = this.events[event.name];
        for (var i in listeners) {
            listeners[i](event);
        }
    };
}