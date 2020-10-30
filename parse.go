package main

import (
	"encoding/json"
)

type SvcPolicies struct {
	ProxyName string   `json:"proxy_key"`
	Policies  []Policy `json:"rbac_policies"`
}

// Policy is the json structure that comes in from the front end
type Policy struct {
	PolicyName  string       `json:"policy_name"`
	Permissions []Permission `json:"permissions"`
	Principals  []Principal  `json:"principals"`
}

type Permission struct {
	Any            bool        `json:"any"`
	PermissionType string      `json:"permission"`
	Value          interface{} `json:"value"`
}

type Principal struct {
	Any           bool        `json:"any"`
	PrincipalType string      `json:"principal"`
	Value         interface{} `json:"value"`
}

type Header struct {
	Name  string `json:"name"`
	Value string `json":"value"`
}

// TODO we need ordering
func parseRBACPolicies(result map[string]interface{}, proxyKey string) SvcPolicies {
	if result["proxy_filters"] == nil {
		logger.Info().Msgf("No proxy filters enabled on proxy: %s", proxyKey)
		return SvcPolicies{ProxyName: proxyKey}
	}

	proxyfilters := result["proxy_filters"].(map[string]interface{})

	rbaccfg, ok := proxyfilters["envoy_rbac"]
	if !ok {
		logger.Info().Msgf("rbac filter not enabled for proxy: %s", proxyKey)
		return SvcPolicies{ProxyName: proxyKey}
	}

	logger.Info().Msgf("rbac filter enabled for proxy: %s, parsing policies", proxyKey)

	rbacCfg := rbaccfg.(map[string]interface{})
	policyList := []Policy{}
	rules := rbacCfg["rules"].(map[string]interface{})
	policies := rules["policies"].(map[string]interface{})

	for p, cfg := range policies {
		cfgPerm := cfg.(map[string]interface{})["permissions"]
		cfgPrinc := cfg.(map[string]interface{})["principals"]
		permissions := parsePermissions(cfgPerm.([]interface{}))
		principals := parsePrincipals(cfgPrinc.([]interface{}))
		pc := Policy{
			PolicyName:  p,
			Permissions: permissions,
			Principals:  principals,
		}
		policyList = append(policyList, pc)
	}

	sPolicies := SvcPolicies{
		ProxyName: proxyKey,
		Policies:  policyList,
	}

	return sPolicies
}

func parsePermissions(perm []interface{}) []Permission {
	permission := []Permission{}
	for _, p := range perm {
		pt := p.(map[string]interface{})
		for ptype, cfg := range pt {
			perm := Permission{}
			if ptype == "any" {
				perm.Any = true
			} else {
				perm.PermissionType = ptype
				perm.Value = permissionValue(ptype, cfg)
			}
			permission = append(permission, perm)
		}
	}
	return permission
}

// currently only uses exact match and header TODO
func permissionValue(ptype string, cfg interface{}) interface{} {
	pcfg := cfg.(map[string]interface{})
	switch ptype {
	case "header":
		pValue := Header{
			Name:  pcfg["name"].(string),
			Value: pcfg["exact_match"].(string),
		}
		return pValue
	}

	return nil
}

func parsePrincipals(princ []interface{}) []Principal {
	principals := []Principal{}
	for _, p := range princ {
		pt := p.(map[string]interface{})
		for ptype, cfg := range pt {
			pr := Principal{}
			if ptype == "any" {
				pr.Any = true
			} else {
				pr.PrincipalType = ptype
				pr.Value = principalValue(ptype, cfg)
			}
			principals = append(principals, pr)
		}
	}
	return []Principal{}
}

func principalValue(ptype string, cfg interface{}) interface{} {
	pcfg := cfg.(map[string]interface{})
	switch ptype {
	case "header":
		pValue := Header{
			Name:  pcfg["name"].(string),
			Value: pcfg["exact_match"].(string),
		}
		return pValue
	}

	return nil
}

func rbacify(policy Policy) json.RawMessage {
	return json.RawMessage{}
}

func getAllSvcPolicies(proxies []interface{}) []SvcPolicies {
	allSvc := []SvcPolicies{}
	for _, p := range proxies {
		proxycfg := p.(map[string]interface{})
		pKey := proxycfg["proxy_key"].(string)
		policies := parseRBACPolicies(proxycfg, pKey)
		allSvc = append(allSvc, policies)

	}
	return allSvc
}
