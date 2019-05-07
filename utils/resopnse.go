package utils

import (
	"encoding/json"
	"net/http"
)

type result struct {
	Data   interface{} `json:"data"`
	Msg    string      `json:"msg"`
	Status int         `json:"status"`
}

func (r *result) return_json() string {
	b, _ := json.MarshalIndent(r, "", "	")
	return string(b)
}

func Return_JSON_200_For_Default(w http.ResponseWriter) {
	rs := new(result)
	rs.Data = make(map[string]interface{})
	rs.Msg = "success"
	rs.Status = 200

	w.Write([]byte(rs.return_json()))
	// stop_run()
}

func Return_JSON_200(w http.ResponseWriter, data interface{}) {
	rs := new(result)
	rs.Data = data
	rs.Msg = "success"
	rs.Status = 200

	w.Write([]byte(rs.return_json()))
	// stop_run()
}

func Return_JSON_220(w http.ResponseWriter) {
	rs := new(result)
	rs.Data = make(map[string]interface{})
	rs.Msg = "success"
	rs.Status = 220

	w.Write([]byte(rs.return_json()))
	// stop_run()
}

func Return_JSON_400(w http.ResponseWriter, msg string) {
	rs := new(result)
	rs.Data = make(map[string]interface{})
	rs.Msg = msg
	rs.Status = 400

	w.Write([]byte(rs.return_json()))
	// stop_run()
}

func Return_JSON_412(w http.ResponseWriter) {
	rs := new(result)
	rs.Data = make(map[string]interface{})
	rs.Msg = "authentication failed"
	rs.Status = 412

	w.Write([]byte(rs.return_json()))
	// stop_run()
}

func Return_Content_500(w http.ResponseWriter, msg string) {
	w.WriteHeader(500)
	w.Write([]byte(msg))
}
