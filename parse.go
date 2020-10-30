package main

import "encoding/json"

// Policy is the json structure that comes in from the front end
type Policy struct {
	ServiceName string     `json:"service_name"`
	PolicyName  string     `json:"policy_name"`
	Permissions Permission `json:"permissions"`
	Principals  Principal  `json:"principals"`
}

type Permission struct {
	Any   bool
	Type  string
	Value string
}

type Principal struct {
	Any   bool
	Type  string
	Value string
}

type Header struct {
	Name  string
	Type  string
	Match string
}

func rbacify(policy Policy) json.RawMessage {
	return json.RawMessage{}
}
