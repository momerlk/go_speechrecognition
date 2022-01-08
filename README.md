# go_speechrecognition
Ibm Watson Speech Recognition wrapper in go

extremely SIMPLE USAGE :


Import the package
`import (
  watson "github.com/dragonmaster101/go_speechrecognition/watson_sr"
)`


Inside of the main Function Create the Service
`service := watson.CreateService("IBM_WATSON_SPEECHRECOGNITION_KEY" , "URL")`

Transcribe a file and get the result as a string
`var transcription string = watson.TranscribeFile(service , "filename.mp3")`

Print the results
`fmt.Println(transcription)`
