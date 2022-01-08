package watson_sr

import (
	// "encoding/json"
	// "bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	api "github.com/watson-developer-cloud/go-sdk/v2/speechtotextv1"
	// "github.com/dragonmaster101/go_audio/audio"
)

type SpeechToTextService *api.SpeechToTextV1

type RawResponse struct {
	StatusCode int `json:"StatusCode"`
	Headers    struct {
		Connection                   []string `json:"Connection"`
		ContentDisposition           []string `json:"Content-Disposition"`
		ContentLength                []string `json:"Content-Length"`
		ContentSecurityPolicy        []string `json:"Content-Security-Policy"`
		ContentType                  []string `json:"Content-Type"`
		Date                         []string `json:"Date"`
		Server                       []string `json:"Server"`
		SessionName                  []string `json:"Session-Name"`
		StrictTransportSecurity      []string `json:"Strict-Transport-Security"`
		XContentTypeOptions          []string `json:"X-Content-Type-Options"`
		XDpWatsonTranID              []string `json:"X-Dp-Watson-Tran-Id"`
		XEdgeconnectMidmileRtt       []string `json:"X-Edgeconnect-Midmile-Rtt"`
		XEdgeconnectOriginMexLatency []string `json:"X-Edgeconnect-Origin-Mex-Latency"`
		XFrameOptions                []string `json:"X-Frame-Options"`
		XGlobalTransactionID         []string `json:"X-Global-Transaction-Id"`
		XRequestID                   []string `json:"X-Request-Id"`
		XXSSProtection               []string `json:"X-Xss-Protection"`
	} `json:"Headers"`
	Result struct {
		Results []struct {
			Final        bool `json:"final"`
			Alternatives []TrueResult `json:"alternatives"`
		} `json:"results"`
		ResultIndex int `json:"result_index"`
	} `json:"Result"`
	RawResult interface{} `json:"RawResult"`
}

type TrueResult struct {
	Transcript string  `json:"transcript"`
	Confidence float64 `json:"confidence"`
}


func CreateService(API_KEY string , API_URL string) *api.SpeechToTextV1{
	authenticator := &core.IamAuthenticator{
		ApiKey: API_KEY,
	};


	service , err := api.NewSpeechToTextV1(&api.SpeechToTextV1Options{
		URL:       API_URL,
		Authenticator: authenticator,
	});
	check(err);

	return service;
}

func TranscribeFile(service *api.SpeechToTextV1, fileName string) string{

	file , err := os.Open(fileName);
	check(err);
	defer file.Close();

	reader := ioutil.NopCloser(file);

	recognitionOptions := service.NewRecognizeOptions(reader);
	_ , resp , err := service.Recognize(recognitionOptions);
	check(err);

	rawResponse := RawResponse{};

	data := resp.String();
	json.NewDecoder(strings.NewReader(data)).Decode(&rawResponse);

	return rawResponse.Result.Results[0].Alternatives[0].Transcript;
}
	
type Callback struct {
	Status string `json:"status"`
	URL    string `json:"url"`
}

func RegisterCallback(service *api.SpeechToTextV1 , url string) Callback{
	result, response, responseErr := service.RegisterCallback(
	
		&api.RegisterCallbackOptions{
		CallbackURL: core.StringPtr(url),
		// UserSecret:  core.StringPtr(Secret),
	},
	)

	if responseErr != nil {
		fmt.Println("Registering callback Response unsuccessful");
		fmt.Println(responseErr);
		fmt.Println(response);
		os.Exit(1);
	}

	b, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(b))
	data := response.String();
	
	decodedResponse := Callback{};
	json.NewDecoder(strings.NewReader(data)).Decode(&decodedResponse);

	fmt.Println("Callback Response :" , data);

	return decodedResponse;
}

func CreateJob(service *api.SpeechToTextV1 , callback Callback , fileName string){


	var audioFile io.ReadCloser
	var audioFileErr error
	audioFile, audioFileErr = os.Open(fileName)
	check(audioFileErr);

	_, rawResponse, responseErr := service.CreateJob(
	&api.CreateJobOptions{
		Audio:                     audioFile,
		Timestamps:                core.BoolPtr(false),
		ProfanityFilter: 		   core.BoolPtr(false),
		CallbackURL:               &callback.URL,
		UserToken:                 core.StringPtr("job25"),
	},
	)
	check(responseErr);	

	fmt.Println("create Job response :" ,rawResponse.String());
}

	
type JobResponse struct {
	Created time.Time `json:"created"`
	ID      string    `json:"id"`
	Updated time.Time `json:"updated"`
	Results []struct {
		ResultIndex int `json:"result_index"`
		Results     []struct {
			Final        bool `json:"final"`
			Alternatives []struct {
				Transcript string          `json:"transcript"`
				Timestamps [][]interface{} `json:"timestamps"`
				Confidence float64         `json:"confidence"`
			} `json:"alternatives"`
		} `json:"results"`
	} `json:"results"`
	Status string `json:"status"`
}

func CheckJob(service *api.SpeechToTextV1) JobResponse{

	_, rawResponse, responseErr := service.CheckJob(
	&api.CheckJobOptions{
		ID: core.StringPtr("{id}"),
	},
	)
	check(responseErr);
	data := rawResponse.String();
	response := JobResponse{};

	json.NewDecoder(strings.NewReader(data)).Decode(&response);

	return response;
}


func check(err error){
	if err != nil {
		fmt.Println(err.Error());
		fmt.Println(err);
		os.Exit(1);
	}
}
