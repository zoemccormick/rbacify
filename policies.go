package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/viper"
)

func getPolicy(w http.ResponseWriter, r *http.Request) {
	logger.Info().Msg("Generating get request to control-api service")
	resp, err := http.Get("http://" + viper.GetString("gm_control_api_address") + "/v1/cluster")
	if err != nil {
		logger.Error().Msg(err.Error())
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error().Msg(err.Error())
	}
	w.Write(body)
}

func addPolicy(w http.ResponseWriter, r *http.Request) {
	logger.Info().Msg("Creating policy to add")
	defer func() {
		if err := r.Body.Close(); err != nil {
			logger.Error().AnErr("request.Body.Close()", err).Msg("close error")
		}
	}()

	var policy Policy
	err := json.NewDecoder(r.Body).Decode(&policy)
	if err != nil {
		logger.Error().Msg("Received bad json policy")
		http.Error(w, fmt.Sprintf("unable to decode json: %v", err.Error()), http.StatusBadRequest)
		return
	}

	//rbacPolicy := rbacify(policy)

	json.NewEncoder(w).Encode(policy)
}

//func postPolicy(proxyKey string, rbacPolicy json.RawMessage) {
//	method := "PUT"
//	path := "/v1/proxy/" + proxyKey
//
//}

func deletePolicy(w http.ResponseWriter, r *http.Request) {
	logger.Info().Msg("Deleting Existing Policy")

}

func modifyPolicy(w http.ResponseWriter, r *http.Request) {
	logger.Info().Msg("Modifying Policy")

}
