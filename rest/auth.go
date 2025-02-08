package rest

// import (
// 	"encoding/json"
// 	"net/http"
// 	"strings"

// 	"backend-service/util"

// 	"github.com/julienschmidt/httprouter"
// )

// func Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
// 	ctx := r.Context()
// 	log := util.GetGlobalLogger(ctx)

// 	credentials := make(map[string]string)
// 	err := json.NewDecoder(r.Body).Decode(&credentials)
// 	if err != nil {
// 		log.Println("request json parsing error", err)
// 		http.Error(w, "request body not proper", http.StatusBadRequest)
// 		return
// 	}

// 	if strings.EqualFold(credentials["email"], "") {
// 		log.Println("email is empty")
// 		http.Error(w, "email not sent", http.StatusBadRequest)
// 		return
// 	}
// 	if strings.EqualFold(credentials["password"], "") {
// 		log.Println("password is empty")
// 		http.Error(w, "password not sent", http.StatusBadRequest)
// 		return
// 	}
// }
