package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang-citcall-api/config"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type response struct {
	Error   bool        `json:"error"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

/*
 *
 * @desc Running Config and Server
 *
 */
func Run() {
	config.Load()
	listen()
}

/*
 *
 * @desc Running Router and Server
 *
 */
func listen() {
	r := gin.Default()
	r.POST("/misscall", misscallHandler)

	r.Run(":" + strconv.Itoa(config.PORT))
}

/*
 * POST : '/misscall'
 *
 * @desc Send request to citcall
 *
 * @param  {msisdn} string - No Handphone
 * @param  {gateway} int - Gateway
 *
 * @return {object} Request object
 */
func misscallHandler(c *gin.Context) {
	var data = make(map[string]interface{})
	var result = make(map[string]interface{})
	var w = c.Writer

	// bind data from form
	data["msisdn"] = c.Request.FormValue("msisdn")
	data["gateway"] = c.Request.FormValue("gateway")

	request, err := json.Marshal(data)
	if err != nil {
		errorCustomStatus(w, http.StatusInternalServerError, err.Error())
		return
	}

	// API Request
	url := fmt.Sprintf("%s%s", os.Getenv("BASE_URL_CITCAL"), os.Getenv("CITCALL_MISSCALL"))
	key := fmt.Sprintf("Apikey %s", config.APIKEY)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(request))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", key)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		errorCustomStatus(w, http.StatusInternalServerError, err.Error())
		return
	}

	defer resp.Body.Close()

	// get body from request
	body, _ := ioutil.ReadAll(resp.Body)

	// decode json and bind to result variable
	err = json.Unmarshal([]byte(body), &result)

	// display response
	responses(w, http.StatusOK, "Success", result)
	return
}

/*
 *
 * @desc Generate response
 *
 * @param  {w} http.ResponseWriter - Response Writer
 * @param  {code} int - Status Code
 * @param  {msg} string - Message status code
 * @param  {payload} interface{} - Data that display to response
 *
 */
func responses(w http.ResponseWriter, code int, msg string, payload interface{}) {
	var result response
	if code >= 400 {
		result.Error = true
	} else {
		result.Error = false
	}
	result.Code = code
	result.Message = msg
	result.Data = payload

	responses, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(responses)

}

/*
 *
 * @desc Generate Error Response
 *
 * @param  {w} http.ResponseWriter - Response Writer
 * @param  {code} int - Status Code
 * @param  {msg} string - Error Message
 *
 */
func errorCustomStatus(w http.ResponseWriter, code int, msg string) {
	responses(w, code, "Error", msg)
}
