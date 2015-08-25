var updateTypes = {
    NewRequest: 0,
    NewResponse: 1,
    DebuggingToggle: 2,
    InitialUpdate: 3
};

var commandTypes = {
    EnableDebugging: 0,
    DisableDebugging: 1,
    ContinueDebugging: 2
};

window.onload = function() {
    var toggleDebugging = function(enabled) {
        if(enabled === true) {
            console.log('debugger has started');
            $('#debug_stop').prop('disabled', false);
            $('#debug_continue').prop('disabled', false);
            $('#debug_start').prop('disabled', true);
        } else {
            console.log('debugger has stopped');
            $('#debug_stop').prop('disabled', true);
            $('#debug_continue').prop('disabled', true);
            $('#debug_start').prop('disabled', false);
        }
    };

    // listen on websocket
    var socket = new WebSocket("ws://" + window.location.host + "/_socket");
    socket.onmessage = function(msg) {
        console.log(msg.data);
        var receivedData = JSON.parse(msg.data);
        switch(receivedData.Type) {
        case updateTypes.NewRequest:
            $('#request').text(receivedData.Request);
            $('#response').text('');
            break;
        case updateTypes.NewResponse:
            $('#response').text(receivedData.Response);
            break;
        case updateTypes.DebuggingToggle:
            toggleDebugging(receivedData.DebuggingEnabled);
            break;
        case updateTypes.InitialUpdate:
            toggleDebugging(receivedData.DebuggingEnabled);
        default:
            console.log('Unknown update type: ' + receivedData.Type);
        }
    };

    $('#debug_continue').on('click', function() {
        socket.send(JSON.stringify({ type: commandTypes.ContinueDebugging, value: '' }));
    });

    $('#debug_start').on('click', function() {
        console.log('starting debugger');
        socket.send(JSON.stringify({ type: commandTypes.EnableDebugging, value: '' }));
    });

    $('#debug_stop').on('click', function() {
        console.log('stopping debugger');
        socket.send(JSON.stringify({ type: commandTypes.DisableDebugging, value: '' }));
    });
}
