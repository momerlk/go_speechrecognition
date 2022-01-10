package speechrecognition

import (
	"fmt"
	"log"
	"os"
	"strings"

	audio "github.com/dragonmaster101/go_audio/audio"
	capture "github.com/dragonmaster101/go_speechrecognition/capture"
	watson "github.com/dragonmaster101/go_speechrecognition/watson_sr"
)

type recognizer interface {
	TranscribeFile(fileName string) (string , error)
	AsyncTranscribe(res chan string , mic Microphone , done chan bool)
}
type Recognizer recognizer

type Source interface {
	GetData() *os.File
}

type Microphone interface {
	Record(fileName string , timeLimit int)
}

func NewMicrophone() Microphone{
	mic := &portAudioMicrophone{};
	mic.Init();

	return mic;
}

// uses a microphone to Recognize audio and then transcribes it using the recognizer
func Recognize(recognizer Recognizer , mic Microphone) string {
	recordingFile := "recording101.aiff";
	mic.Record(recordingFile , 4);
	
	res , err := recognizer.TranscribeFile(recordingFile);
	if err != nil {
		log.Fatal(err);
	}

	return res;
}

func AsyncRecognize(recognizer Recognizer , mic Microphone , res chan string , done chan bool){
	recognizer.AsyncTranscribe(res , mic , done);
}

/*
<---------------------------------------------------------------------------------------------------------------------->
									<<<	Concrete Type Implementations >>>
*/


type RecognizerWatson struct {
	watson.SpeechToTextService
}

func (r *RecognizerWatson) AsyncTranscribe(resChannel chan string , mic Microphone , done chan bool) {
	r.SocketsTranscribe(
		func(socket watson.Socket , state *int) {

			fmt.Println("Recording! Say Something");
			sourceFileName := "recording.aiff"
			capture.CaptureWithCGo(func(){mic.Record(sourceFileName , 4)});
			fmt.Println("Recording Complete!");
			watson.SendSocketRecognitionRequest(socket , sourceFileName , state);

		}, 
		func (result string) {
			resChannel <- result;
		} , done);
}

// Creates A Service instance from the Watson Speech To Text V1 sdk
func RecognizerIBM(API_KEY , API_URL string) Recognizer {

	dotSplit := strings.Split(API_URL , ".");
	location := dotSplit[1];
	
	slashSplit := strings.Split(API_URL , "/");
	instanceId := slashSplit[len(slashSplit)-1];

	service := watson.CreateService(API_KEY , API_URL , location , instanceId);

	recognizer := &RecognizerWatson{
		watson.SpeechToTextService{
			Service : service.Service,
			API_KEY : service.API_KEY,
			API_URL : service.API_URL,
			ACCESS_TOKEN : service.ACCESS_TOKEN,
			LOCATION : service.LOCATION,
			INSTANCE_ID : service.INSTANCE_ID,
		},
	};

	return recognizer;
}

// uses the portaudio library to record '.aiff' files from the microphone PORTAUDIO SHOULD BE INSTALLED
type portAudioMicrophone struct {
	recorded bool
	fileName string
}

func (p *portAudioMicrophone) Init(){
	p.recorded = false;
}

func (p *portAudioMicrophone) Record(fileName string , timeLimit int){
	p.recorded = true;
	p.fileName = fileName;
	audio.Record(fileName , timeLimit);
} 

func (p *portAudioMicrophone) GetData() *os.File{
	
	if !p.recorded {
		p.fileName = "mic101.aiff"
		p.Record(p.fileName , 4);
	}

	file , err := os.Open(p.fileName);
	if err != nil {
		log.Fatal(err);
	}
	defer file.Close();

	return file;
}