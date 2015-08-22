var updateTypes = {
    NewRequest: 0,
    NewResponse: 1,
    DebuggingEnabled: 2,
    DebuggingDisabled: 3
};

var commandTypes = {
    EnableDebugging: 0,
    DisableDebugging: 1,
    ContinueDebugging: 2
};

window.onload = function() {

    // listen on websocket
    var socket = new WebSocket("ws://" + window.location.host + "/_socket");
    socket.onmessage = function(msg) {
        console.log(msg.data);
        var receivedData = JSON.parse(msg.data);
        switch(receivedData.Type) {
        case updateTypes.NewRequest:
            $('#request').text(receivedData.Value);
            $('#response').text('');
            break;
        case updateTypes.NewResponse:
            $('#response').text(receivedData.Value);
            break;
        case updateTypes.DebuggingEnabled:
            console.log('debugger has started');
            $('#debug_stop').prop('disabled', false);
            $('#debug_continue').prop('disabled', false);
            $('#debug_start').prop('disabled', true);
            break;
        case updateTypes.DebuggingDisabled:
            console.log('debugger has stopped');
            $('#debug_stop').prop('disabled', true);
            $('#debug_continue').prop('disabled', true);
            $('#debug_start').prop('disabled', false);
            break;
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
