package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var hueConfigData HueConfigData

var infoLogger *log.Logger
var warnLogger *log.Logger
var errorLogger *log.Logger

func main() {
	configureLogging()
	loadConfig()

	http.HandleFunc("/huelights/changecolour", callHueBridge)
	http.HandleFunc("/huelights/code", requestClientKey)

	fmt.Println("Philips Hue Handler running, waiting for request...")

	if hueConfigData.Port != "" {
		fmt.Println("HUE Server listening on http://localhost:" + hueConfigData.Port)
		err := http.ListenAndServe("0.0.0.0:"+hueConfigData.Port, nil)
		if err != nil {
			errorLogger.Fatal("Could not start Philips Hue Handler server: ", err)
		}
	} else {
		infoLogger.Println("Hue server listening on http://localhost:8080")
		err := http.ListenAndServe("0.0.0.0:8080", nil)
		if err != nil {
			errorLogger.Fatal("Could not start Philips Hue Handler server: ", err)
		}
	}
}

func configureLogging() {
	file, _ := os.Create("hueApp.log")
	flags := log.Ldate | log.Lshortfile
	infoLogger = log.New(file, "INFO: ", flags)
	warnLogger = log.New(file, "WARN: ", flags)
	errorLogger = log.New(file, "ERROR: ", flags)
}

func loadConfig() {
	args := os.Args
	var configPath string
	if len(args) < 2 {
		configPath = "config.json"
	} else {
		configPath = args[1]
	}

	jsonFile, err := os.Open(configPath)
	if err != nil {
		errorLogger.Fatal("Config load error: ", err)
	}
	jsonByteValue, _ := io.ReadAll(jsonFile)
	err = json.Unmarshal(jsonByteValue, &hueConfigData)
	if err != nil {
		errorLogger.Fatal("Config load error: ", err)
	}
}

func callHueBridge(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		infoLogger.Println("Method not allowed: ", r.Method)
		return
	}

	reqBody, _ := io.ReadAll(r.Body)

	var hueLightValues HueLightValue

	err := json.Unmarshal(reqBody, &hueLightValues)
	if err != nil {
		warnLogger.Println("Cannot unmarshal json: ", reqBody)
		return
	}

	infoLogger.Println("light data type     : ", hueLightValues.LightDataType)
	infoLogger.Println("light x value     : ", hueLightValues.X_Data)
	infoLogger.Println("light y value     : ", hueLightValues.Y_Data)
	infoLogger.Println("brightness     : ", hueLightValues.Brightness)
	infoLogger.Println("mirek val     : ", hueLightValues.Mirek_value)

	var urlStr = "https://" + hueConfigData.Hue_IP_address + "/clip/v2/resource/grouped_light/" + hueConfigData.Hue_grouped_light_id
	var jsonLightRequest []byte

	if hueLightValues.LightDataType == "white" {
		var temperatureStruct LightTemperature
		temperatureStruct.ColourTemperature.Mirek = hueLightValues.Mirek_value
		temperatureStruct.Dimming_Struct.Brightness = hueLightValues.Brightness
		jsonLightRequest, err = json.Marshal(temperatureStruct)
		if err != nil {
			warnLogger.Println("Could not form white light request: ", temperatureStruct)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if hueLightValues.LightDataType == "colour" {
		var colourStruct LightColour
		colourStruct.Colour_Struct.ColourData.X_Val = hueLightValues.X_Data
		colourStruct.Colour_Struct.ColourData.Y_Val = hueLightValues.Y_Data
		colourStruct.Dimming_Struct.Brightness = hueLightValues.Brightness

		jsonLightRequest, err = json.Marshal(colourStruct)
		if err != nil {
			warnLogger.Println("Could not form coloured light request: ", colourStruct)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		warnLogger.Println("Cannot determine light data type, recerived: ", hueLightValues.LightDataType)
		http.Error(w, "Cannot determine Light data type, should be \"white\" or \"colour\"", http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest("PUT", urlStr, bytes.NewBuffer(jsonLightRequest))
	if err != nil {
		errorLogger.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("hue-application-key", hueConfigData.Hue_app_key)
	resp, err := sendHubRequest(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		errorLogger.Println("Error sending request:", err)
		return
	}

	n, err := w.Write(resp)
	if err != nil {
		warnLogger.Println("Cannot write hub response:", err, " Code: ", n)
		http.Error(w, "Cannot write coulour change response", http.StatusInternalServerError)
	}
}

func requestClientKey(w http.ResponseWriter, r *http.Request) {
	hueAuth := HueAuth{
		DeviceType:        "app_name#instance_name",
		GenerateClientKey: true,
	}

	requestBody, err := json.Marshal(hueAuth)
	if err != nil {
		errorLogger.Println("Error marshaling JSON:", err)
		return
	}

	urlStr := "https://" + hueConfigData.Hue_IP_address + "/api"
	req, err := http.NewRequest("POST", urlStr, bytes.NewBuffer(requestBody))
	if err != nil {
		errorLogger.Println("Error forming POST request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	hubResponse, err := sendHubRequest(req)
	if err != nil {
		http.Error(w, string(hubResponse), http.StatusInternalServerError)
	}

	n, err := w.Write(hubResponse)
	if err != nil {
		warnLogger.Println("Cannot write hub response:", err, " Code: ", n)
		http.Error(w, "Cannot write client key", http.StatusInternalServerError)
	}
}

func sendHubRequest(req *http.Request) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		errorLogger.Println("Error sending server request:", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		errorLogger.Println("Error reading server response:", err)
		return nil, err
	}
	infoLogger.Println("Response from", req.URL, ":", string(body))
	return body, nil
}