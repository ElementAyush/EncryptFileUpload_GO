package resphandle

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-services/gomod/model"
)

/**
* This is the generic response handler function
*
* @parm w http.ResponseWriter used to send response
* @param statusCode response status code
* @param sucesss denoting if the request was successfull
* @param description detailed message about the response
* return null
 */
func HandleResp(w http.ResponseWriter, statusCode int, success bool, description string) {
	dataJson := &model.Error{Success: success, Description: description}
	dat, responseerr := json.MarshalIndent(dataJson, "", "  ")
	if responseerr != nil {
		log.Println("Not able to parse json")
	}
	w.WriteHeader(statusCode)
	w.Write(dat)
}
