package main
import (
	"io"
	"os"
	"time"
	"github.com/gorilla/websocket"
)

var serverURL = "ws://localhost:7891/stt"

func stream_send(){
	time.Sleep(5 * time.second)

	ws, _, err:= websocket.DefaultDialer.Dial(serverURL,nil)
	if err != nil {
		Error(err)
	}

	defer ws.Close()

	file, err :=os.Open(<audio file>) //Upload audio file
	if err != nil {
		Error(err)
	}

	defer file.Close()

	buffer := make([]byte, 1024)
	var total int

	defer func(){
		Debug("File sent total:", total)
	}()
	
	go func(){
		// Loop to read buffer
		for {
			n, err: file.Read(buffer)

			total +=n

			if err !=nil && err != io.EOF{
				Error(err)
				return
			}

			if n==0{
				if err == io.eof{
					Error("EOF")
				}
				break
			}
			err = ws.WriteMessage(websocket.TextMessage, buffer[:n]) // Send to server
			if err ~= nil {
				break
			}
			Debug("Client: ",string(message))
		}
	}
}