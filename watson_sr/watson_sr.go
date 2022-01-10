package watson_sr

import (
	"fmt"
	"os"
	"strings"
	"log"

	"os/signal"
	"io/ioutil"
	"net/http"
	"net/url"
	"encoding/json"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/sacOO7/gowebsocket"
	api "github.com/watson-developer-cloud/go-sdk/v2/speechtotextv1"
)

type SpeechToTextService struct {
	Service *api.SpeechToTextV1
	API_KEY string 
	API_URL string 
	ACCESS_TOKEN string 
	LOCATION string
	INSTANCE_ID string
}

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

func CreateApiService(API_KEY string , API_URL string) *api.SpeechToTextV1 {
	authenticator := &core.IamAuthenticator{
		ApiKey: API_KEY,
	};


	apiService , err := api.NewSpeechToTextV1(&api.SpeechToTextV1Options{
		URL:       API_URL,
		Authenticator: authenticator,
	});
	check(err);

	return apiService;
}

func CreateService(API_KEY string , API_URL string , LOCATION string , 
	INSTANCE_ID string) *SpeechToTextService {

	authenticator := &core.IamAuthenticator{
		ApiKey: API_KEY,
	};


	apiService , err := api.NewSpeechToTextV1(&api.SpeechToTextV1Options{
		URL:       API_URL,
		Authenticator: authenticator,
	});
	check(err);
	
	accessToken , _ := GetAccessToken(API_KEY ,API_URL); 
	service := SpeechToTextService{Service : apiService , API_KEY: API_KEY,
		API_URL: API_URL, ACCESS_TOKEN: accessToken, LOCATION: LOCATION,
		INSTANCE_ID: INSTANCE_ID,
	};

	return &service;
}

func (s *SpeechToTextService) TranscribeFile(fileName string) (str string, err error){

	service := s.Service;
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

	if rawResponse.Result.Results == nil {
		err = fmt.Errorf("IBM Watson sent back nil Results (TRY AGAIN)");
		return "None" , err;
	}

	return rawResponse.Result.Results[0].Alternatives[0].Transcript , err;
}

type GetTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Expiration   int    `json:"expiration"`
	Scope        string `json:"scope"`
}

// returns the IAM access token for the given API KEY
func GetAccessToken(API_KEY , API_URL string) (result string , err error){
	params := url.Values{}
	params.Add("grant_type", `urn:ibm:params:oauth:grant-type:apikey`)
	params.Add("apikey", `djFD3P479m7_iCNgr0pEdBT78GM_Jzg-b0NUnEYwOPgP`)
	body := strings.NewReader(params.Encode())

	req, err := http.NewRequest("POST", "https://iam.cloud.ibm.com/identity/token", body)
	if err != nil {
		log.Fatal(err);
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err);	
	}
	defer resp.Body.Close()

	respBody := GetTokenResponse{};

	json.NewDecoder(resp.Body).Decode(&respBody);

	result = respBody.AccessToken;

	return result , err;
}

// type Callback struct {
// 	Status string `json:"status"`
// 	URL    string `json:"url"`
// }

// func RegisterCallback(service *api.SpeechToTextV1 , url string) Callback{
// 	result, response, responseErr := service.RegisterCallback(
	
// 		&api.RegisterCallbackOptions{
// 		CallbackURL: core.StringPtr(url),
// 		// UserSecret:  core.StringPtr(Secret),
// 	},
// 	)

// 	if responseErr != nil {
// 		fmt.Println("Registering callback Response unsuccessful");
// 		fmt.Println(responseErr);
// 		fmt.Println(response);
// 		os.Exit(1);
// 	}

// 	b, _ := json.MarshalIndent(result, "", "  ")
// 	fmt.Println(string(b))
// 	data := response.String();
	
// 	decodedResponse := Callback{};
// 	json.NewDecoder(strings.NewReader(data)).Decode(&decodedResponse);

// 	fmt.Println("Callback Response :" , data);

// 	return decodedResponse;
// }

// func CreateJob(service *api.SpeechToTextV1 , callback Callback , fileName string){


// 	var audioFile io.ReadCloser
// 	var audioFileErr error
// 	audioFile, audioFileErr = os.Open(fileName)
// 	check(audioFileErr);

// 	_, rawResponse, responseErr := service.CreateJob(
// 	&api.CreateJobOptions{
// 		Audio:                     audioFile,
// 		Timestamps:                core.BoolPtr(false),
// 		ProfanityFilter: 		   core.BoolPtr(false),
// 		CallbackURL:               &callback.URL,
// 		UserToken:                 core.StringPtr("job25"),
// 	},
// 	)
// 	check(responseErr);	

// 	fmt.Println("create Job response :" ,rawResponse.String());
// }

	
// type JobResponse struct {
// 	Created time.Time `json:"created"`
// 	ID      string    `json:"id"`
// 	Updated time.Time `json:"updated"`
// 	Results []struct {
// 		ResultIndex int `json:"result_index"`
// 		Results     []struct {
// 			Final        bool `json:"final"`
// 			Alternatives []struct {
// 				Transcript string          `json:"transcript"`
// 				Timestamps [][]interface{} `json:"timestamps"`
// 				Confidence float64         `json:"confidence"`
// 			} `json:"alternatives"`
// 		} `json:"results"`
// 	} `json:"results"`
// 	Status string `json:"status"`
// }

// func CheckJob(service *api.SpeechToTextV1) JobResponse{

// 	_, rawResponse, responseErr := service.CheckJob(
// 	&api.CheckJobOptions{
// 		ID: core.StringPtr("{id}"),
// 	},
// 	)
// 	check(responseErr);
// 	data := rawResponse.String();
// 	response := JobResponse{};

// 	json.NewDecoder(strings.NewReader(data)).Decode(&response);

// 	return response;
// }

type ResultType struct {
	ResultIndex int `json:"result_index"`
	Results     []struct {
		Final        bool `json:"final"`
		Alternatives []struct {
			Transcript string  `json:"transcript"`
			Confidence float64 `json:"confidence"`
		} `json:"alternatives"`
	} `json:"results"`
}

type StatusType struct {
	State string `json:"state"`
}

type Socket gowebsocket.Socket;

func SendSocketRecognitionRequest(socketParam Socket, source string, state *int){
	socket := gowebsocket.Socket(socketParam);
	socket.SendText(`{"action":"start"}`);
	var data []byte;
	var err error
	data , err = ioutil.ReadFile(source);
	if err != nil {
		log.Fatal(err);
	} 
	fmt.Println("bytes sent =" , len(data));
	socket.SendBinary(data);
	socket.SendText(`{"action":"stop"}`);
	*state = *state + 1;
}

func (s *SpeechToTextService) SocketsTranscribe(toExecute func(socket Socket , state *int) ,
	handleResult func(res string) , done chan bool){

	interrupt := make(chan os.Signal, 1);
	accessToken := s.ACCESS_TOKEN;
	signal.Notify(interrupt , os.Interrupt);

	wsUrl := fmt.Sprintf("wss://api.%s.speech-to-text.watson.cloud.ibm.com/instances/%s/v1/recognize",
	s.LOCATION , s.INSTANCE_ID);

	socket := gowebsocket.New(fmt.Sprintf(wsUrl + "?access_token=" + accessToken));


	var reqsVar = 0;
	requestsSent := &reqsVar;

	socket.OnConnected = func(socket gowebsocket.Socket) {
		log.Println("Connected to Watson Speech To Text Service WebSocket interface");
		toExecute(Socket(socket) , requestsSent);
	};
	
	socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
		log.Println("Could not connect to IBM Watson Web Socket interface", err)
	};
	
	socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
		log.Println("Recieved message " + message);
		decoder := json.NewDecoder(strings.NewReader(message));
		
		Result := ResultType{};
		status := StatusType{};
		err := decoder.Decode(&Result);
		if err != nil {
			log.Fatal(fmt.Errorf("result decoding error in websockets connection to IBM Watson"));
		}
		if len(Result.Results) == 0 {
			err := json.NewDecoder(strings.NewReader(message)).Decode(&status);
			if err != nil {
				log.Fatal(err);	
			}
		} else {
			handleResult(Result.Results[0].Alternatives[0].Transcript);
			fmt.Println("Socket Executing again");
			toExecute(Socket(socket) , requestsSent);
		}
	};
	
	socket.OnBinaryMessage = func(data [] byte, socket gowebsocket.Socket) {
		log.Println("Recieved binary data ", data)
	};
	
	socket.OnPingReceived = func(data string, socket gowebsocket.Socket) {
		log.Println("Recieved ping!")
		fmt.Println(data);
	};
	
	socket.OnPongReceived = func(data string, socket gowebsocket.Socket) {
		log.Println("Recieved pong " + data)
	};
	
	socket.OnDisconnected = func(err error, socket gowebsocket.Socket) {
		log.Println("Disconnected from IBM Watson Speech To Text WebSocket Server")
		return;
	};
	
	socket.Connect()
	
	loop :
		for {
			select {
			case <-done:
				log.Println("Done Socket work")
				log.Println("REQUESTS SENT :" , *requestsSent);
				socket.Close()
				break loop;
			}
		}
	return;
}


func check(err error){
	if err != nil {
		fmt.Println(err.Error());
		fmt.Println(err);
		os.Exit(1);
	}
}
