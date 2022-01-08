package main 

import (
	"fmt"
	watson "github.com/dragonmaster101/go_speechrecognition/watson_sr"
)

func main(){
	service := watson.CreateService("{YOUR_API_KEY}" , "{YOUR_API_URL}");
	transcription := watson.TranscribeFile(service , "reply.mp3");

	fmt.Println("Transcription :" , transcription);
}