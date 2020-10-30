package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	apihttp "github.com/greymatter-io/gm-control-api/api/http/envelope"
	"github.com/spf13/viper"
)

func getPolicy(w http.ResponseWriter, r *http.Request) {
	//TODO make configurable
	proxyKey := getProxyKey(r)
	if proxyKey == "" {
		http.Error(w, "no proxy key specified", http.StatusBadRequest)
		return
	}

	logger.Info().Msgf("Generating get request to control-api service for proxy: %s", proxyKey)
	req := fmt.Sprintf("http://%s/v1.0/proxy/%s", viper.GetString("gm_control_api_address"), proxyKey)
	resp, err := http.Get(req)
	if err != nil || resp.StatusCode != 200 {
		logger.Error().Msgf("request to control-api for proxy key %s failed", proxyKey)
		http.Error(w, fmt.Sprintf("request to control-api for proxy key %s failed", proxyKey), http.StatusBadRequest)
		return
	}

	defer resp.Body.Close()
	var apiresp apihttp.Response
	err = json.NewDecoder(resp.Body).Decode(&apiresp)
	if err != nil {
		logger.Error().Msgf("could not decode response body: %s", err.Error())
	}
	result := (apiresp.Payload).(map[string]interface{})
	policies := parseRBACPolicies(result, proxyKey)

	json.NewEncoder(w).Encode(policies)
}

// getState returns a list of all current proxies in the API and their RBAC policies
func getState(w http.ResponseWriter, r *http.Request) {
	logger.Info().Msgf("Getting state of RBAC")

	req := fmt.Sprintf("http://%s/v1.0/proxy", viper.GetString("gm_control_api_address"))
	resp, err := http.Get(req)
	if err != nil {
		logger.Error().Msg(err.Error())
	}

	defer resp.Body.Close()
	var apiresp apihttp.Response
	err = json.NewDecoder(resp.Body).Decode(&apiresp)
	if err != nil {
		logger.Error().Msgf("could not decode response body: %s", err.Error())
	}
	result := (apiresp.Payload).([]interface{})
	state := getAllSvcPolicies(result)

	json.NewEncoder(w).Encode(state)

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

func deletePolicy(w http.ResponseWriter, r *http.Request) {
	logger.Info().Msg("Deleting Existing Policy")

}

func modifyPolicy(w http.ResponseWriter, r *http.Request) {
	logger.Info().Msg("Modifying Policy")

}

func getProxyKey(r *http.Request) string {
	vars := mux.Vars(r)
	raw, err := url.PathUnescape(vars["proxy"])
	if err != nil {
		logger.Error().Msgf(fmt.Sprintf("url.PathUnescape(%v)", vars["proxy"]), err)
		// if there's an error, just pass back the empty raw string
		// and let the caller sort it out
	}
	return raw
}
