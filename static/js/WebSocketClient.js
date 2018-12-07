function ConnectToAssistantServer()
{
    if ("WebSocket" in window)
    {
        var ws = new WebSocket("ws://127.0.0.1:2080/ws_accept/");

        ws.onopen = function(evt){
            ws.send("{\"cmd\":0,\"verson\":\"1.0.0.0\",\"handle_data\":{\"msg\":\"Hello\"}}");
            console.log("Sending...");
        };

        ws.onmessage = function(evt){
            var msg = JSON.parse(evt.data);

            console.log(msg);
        };

        ws.onclose = function(evt){
            console.log("Closed.");
        };
    }
    else
    {
        alert("Does not support websocket.");
    }
}
