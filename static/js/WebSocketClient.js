function ConnectToAssistantServer()
{
    if ("WebSocket" in window)
    {
        var ws = new WebSocket("ws://127.0.0.1:2080/ws_accept/");

        ws.onopen = function(evt){
            ws.send("Hello WebSockets!");
            alert("Sending...");
        };

        ws.onmessage = function(evt){
            alert("Recived Message:"+evt.data);
        };

        ws.onclose = function(evt){
            alert("Closed.");
        };
    }
    else
    {
        alert("Does not support websocket.");
    }
}
