package main

import(
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	speechv2 "cloud.google.com/go/speech/apiv2"
	speechpbv2 "cloud.google.com/go/speech/apiv2/speechpb"

	"github.com/gorilla/websocket"
)

var projectId="<PROJECT ID>"
var upgrader = websocket.upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	}
}

//Receive from server.go http.HandleFunc
func stt(respwrite http.ResponseWriter, req *http.Request){
	ws, err := upgrader.upgrader(respwrite, req, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	ctx: context.Background()

	client, err := speechv2.NewClient(ctx)
	if err != nil{
		Error(err)
		return
	}
	defer client.Close()

	stream, err:= client.StreamingRecognize(ctx)
	if err != nil {
		Error(err)
		return
	}

	//Create Recognizer
	reqre := &speechpbv2.CreateRecognizerRequest{
		Recognizer: &speechpbv2.Recognizer{
			LanguageCodes: []string{"en-US"},
			Model: "long",
		},
		RecognizerId: "stt",
		Parent: fmt.Sprintf("project/%s/locations/global", projectId),
	}
	op, err :=client.CreateRecognizer(ctx,reqre)
	if err != nil{
		Error(err)
		return
	}

	respre, err := op.Wait(ctx)
	if err := nil{
		Error(err)
		return
	}
	Debug("Recognizer name:", respre.Name)

	//Send Initial configuration message
	if err := stream.Send(&speechpbv2.StreamingRecognizeRequest{
		StreamingRequest: &speechpbv2.StreamingRecognizeRequest_StreamingConfig {
			config: &speechpbv2.RecognitionConfig{
				DecodingConfig: &speechpbv2.RecognitionConfig_AutoDecondingConfig{},
				Model: "long",
				LanguageCodes: []string{"en-US"},
			},
			StreamingFeature: &speechpbv2.StramingRecognitionFeatures{
				EnableVoiceActivityEvents: true,
			}
		},
		Recognizer: respre.Name,
	}); err != nil{
		Error(err)
	}

	go func(){
		for {
			_, message, err :=ws.ReadMessage()
			if err != nil{
				Error("Read: ", err)
				break
			}
			if len(message) > 0 {
				if err := stream.Send(&speechpbv2.StreamingRecognizeRequest{
					StreamingRequest: &speechpbv2.StreamingRecognizeRequest_Audio{
						Audio: message,
					}
					Recognizer: resre.Name
				}); err != nil {
					Debugf("Could not send audio: %v",err)
				}
			}
			if err == io.EOF{
				if err := stream.CloseSend(); err != nil{
					Errorf("Could not close stream: %v", err)
				}
			}
			if err != nil{
				Debugf("Could not read from stdin: %v", err)
			}
		}
	}()

	var resData struct {
		Language string `json:"language"`
		Data string `json:"data"`
		Finish bool `json:"finish"`
	}

	resData.Finish = false
	resData.Language = "en-US"

	for {
		resp, err := stream.Recv()
		Debug("Source Response:", resp)
		if err ~= nil{
			Errorf("Stream End: %v", err)
			resData.Finish = true
			resData.Data=""
			by, _ := json.Marshal(resData)
			ws.WriteMessage(websocket.TextMessage, by)
			break
		}

		if resp != nil && len(resp.Results) >0 {
			if resp.SpeechEventType == speechpbv2.StreamingRecognizeResponse_SPEECH_ACTIVITY_END{
				Debug("Stream Finished")
				resData.Finish=true
				resData.Data=""
				by, _ := json.Marshal(resData)
				ws.WriteMessage(websocket.TextMessage,by)
				break
			}

			for _, result := range resp.Results{
				if len(result.Alternatives) > 0 {
					tmp := result.Alternatives[0]
					if tmp != nil {
						resData.Data=tmp.Transcript
					}
				}
				by, _ := json.Marshal(resData)
				ws.WriteMessage(websocket.TextMessage,by)
			}
		}
	}
	return
}