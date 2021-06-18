package client

import (
	"encoding/json"
)

type RuleEngineClient interface {
	// Stream
	CreateStream(stream string) (string, error)
	ShowStreams() ([]string, error)
	DescribeStream(name string) (string, error)
	UpdateStream(name string, stream string) (string, error)
	DropStream(name string) (string, error)

	// Rule
	CreateRule(rule string) (string, error)
	ShowRules() ([]map[string]interface{}, error)
	DescribeRule(name string) (string, error)
	UpdateRule(name string, rule string) (string, error)
	DropRule(name string) (string, error)
	StatusRule(name string) (string, error)
	StartRule(name string) (string, error)
	StopRule(name string) (string, error)
	RestartRule(name string) (string, error)
}

type KuiperRuleClient struct {
	baseUrl string
}

func NewKuiperRuleClient(baseUrl string) RuleEngineClient {
	return &KuiperRuleClient{
		baseUrl: baseUrl,
	}
}

func (c *KuiperRuleClient) CreateStream(stream string) (string, error) {
	dataResponse, err := SendRequest(c.baseUrl, "streams", "POST", []byte(stream))
	return string(dataResponse), err
}

func (c *KuiperRuleClient) ShowStreams() ([]string, error) {
	dataResponse, err := SendRequest(c.baseUrl, "streams", "GET", nil)
	if err != nil {
		return nil, err
	}

	var response []string
	err = json.Unmarshal(dataResponse, &response)
	if err != nil {
		return nil, err
	}
	return response, err
}

func (c *KuiperRuleClient) DescribeStream(name string) (string, error) {
	dataResponse, err := SendRequest(c.baseUrl, "streams/"+name, "GET", nil)
	return string(dataResponse), err
}

func (c *KuiperRuleClient) UpdateStream(name string, stream string) (string, error) {
	dataResponse, err := SendRequest(c.baseUrl, "streams/"+name, "PUT", []byte(stream))
	return string(dataResponse), err
}

func (c *KuiperRuleClient) DropStream(name string) (string, error) {
	dataResponse, err := SendRequest(c.baseUrl, "streams/"+name, "DELETE", nil)
	return string(dataResponse), err
}

func (c *KuiperRuleClient) CreateRule(rule string) (string, error) {
	dataResponse, err := SendRequest(c.baseUrl, "rules", "POST", []byte(rule))
	return string(dataResponse), err
}

func (c *KuiperRuleClient) ShowRules() ([]map[string]interface{}, error) {
	dataResponse, err := SendRequest(c.baseUrl, "rules", "GET", nil)
	if err != nil {
		return nil, err
	}

	var response []map[string]interface{}
	err = json.Unmarshal(dataResponse, &response)
	if err != nil {
		return nil, err
	}
	return response, err
}

func (c *KuiperRuleClient) DescribeRule(name string) (string, error) {
	dataResponse, err := SendRequest(c.baseUrl, "rules/"+name, "GET", nil)
	return string(dataResponse), err
}

func (c *KuiperRuleClient) UpdateRule(name string, rule string) (string, error) {
	dataResponse, err := SendRequest(c.baseUrl, "rules/"+name, "PUT", []byte(rule))
	return string(dataResponse), err
}

func (c *KuiperRuleClient) DropRule(name string) (string, error) {
	dataResponse, err := SendRequest(c.baseUrl, "rules/"+name, "DELETE", nil)
	return string(dataResponse), err
}

func (c *KuiperRuleClient) StatusRule(name string) (string, error) {
	dataResponse, err := SendRequest(c.baseUrl, "rules/"+name+"/status", "GET", nil)
	return string(dataResponse), err
}

func (c *KuiperRuleClient) StartRule(name string) (string, error) {
	dataResponse, err := SendRequest(c.baseUrl, "rules/"+name+"/start", "POST", nil)
	return string(dataResponse), err
}

func (c *KuiperRuleClient) StopRule(name string) (string, error) {
	dataResponse, err := SendRequest(c.baseUrl, "rules/"+name+"/stop", "POST", nil)
	return string(dataResponse), err
}

func (c *KuiperRuleClient) RestartRule(name string) (string, error) {
	dataResponse, err := SendRequest(c.baseUrl, "rules/"+name+"/restart", "POST", nil)
	return string(dataResponse), err
}
