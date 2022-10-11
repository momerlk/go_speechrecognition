package speechrecognition

import (
	"strings"


	watson "github.com/dragonmaster101/go_speechrecognition/watson_sr"
)

type recognizer interface {
	TranscribeFile(fileName string) (string , error)

}
type Recognizer recognizer



/*
<---------------------------------------------------------------------------------------------------------------------->
									<<<	Concrete Type Implementations >>>
*/


type RecognizerWatson struct {
	watson.SpeechToTextService
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

