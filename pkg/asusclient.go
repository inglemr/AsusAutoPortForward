package pkg

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type AsusRouterClient struct {
	routerAddress string
	username      string
	password      string
	httpClient    *http.Client
}

type AsusTokenResponse struct {
	AsusToken string `json:"asus_token"`
}

type AsusPortForwardRulesResponse struct {
	VtsRulelist string `json:"vts_rulelist"`
}

func NewAsusRouterClient(routerAddress, username, password string) *AsusRouterClient {
	return &AsusRouterClient{
		routerAddress: routerAddress,
		username:      username,
		password:      password,
		httpClient:    http.DefaultClient,
	}
}

func (c *AsusRouterClient) GetAuthToken() AsusTokenResponse {
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(c.username + ":" + c.password))
	payload := strings.NewReader("login_authorization=" + encodedAuth)
	req, _ := http.NewRequest("POST", c.routerAddress+"/login.cgi", payload)
	req.Header.Add("user-agent", "asusrouter-Android-DUTUtil-1.0.0.245")
	res, err := c.httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	tokenResponse := AsusTokenResponse{}
	err = json.NewDecoder(res.Body).Decode(&tokenResponse)
	if err != nil {
		panic(err)
	}
	return tokenResponse
}

func (c *AsusRouterClient) GetPortForwardRules() map[string]PortForwardRule {
	asus_token := c.GetAuthToken().AsusToken
	req, _ := http.NewRequest("GET", c.routerAddress+"/appGet.cgi?hook=nvram_get(vts_rulelist)", nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Cookie", "asus_token="+asus_token)
	req.Header.Add("user-agent", "asusrouter-Android-DUTUtil-1.0.0.245")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	rules := AsusPortForwardRulesResponse{}
	err = json.NewDecoder(resp.Body).Decode(&rules)
	if err != nil {
		log.Printf("Failure Getting Port Forward Rules: %v", err)
	}
	unescaped := rules.VtsRulelist
	unescaped = strings.ReplaceAll(unescaped, "&#60", "<")
	unescaped = strings.ReplaceAll(unescaped, "&#62", ">")
	ruleId := 0
	forwardRules := make(map[string]PortForwardRule)
	rule := PortForwardRule{}
	for _, c := range unescaped {
		if c == '<' {
			continue
		}
		if c == '>' {
			ruleId += 1
			if ruleId > 4 {
				ruleId = 0
				forwardRules[rule.RuleName] = rule
				rule = PortForwardRule{}
			}
			continue
		}
		switch ruleId {
		case 0:
			rule.RuleName += string(c)
		case 1:
			rule.SourcePort += string(c)
		case 2:
			rule.TargetIP += string(c)
		case 3:
			rule.TargetPort += string(c)
		case 4:
			rule.Protocol += string(c)
		}
	}
	return forwardRules
}

func (c *AsusRouterClient) UpdatePortForwardRules(rules []PortForwardRule) {
	asus_token := c.GetAuthToken().AsusToken
	form := url.Values{}
	form.Add("action_mode", "apply")
	form.Add("rc_service", "restart_firewall")
	portForwardRules := ""
	for _, rule := range rules {
		portForwardRules += rule.RouterString()
	}
	form.Add("vts_rulelist", portForwardRules)
	reqbody := strings.NewReader(form.Encode())
	fmt.Printf("Encoded Form: %v\n", form.Encode())
	req, _ := http.NewRequest("POST", c.routerAddress+"/applyapp.cgi", reqbody)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Cookie", "asus_token="+asus_token)
	req.Header.Add("user-agent", "asusrouter-Android-DUTUtil-1.0.0.245")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var body []byte
	resp.Body.Read(body)
	fmt.Println(string(body))
	fmt.Println("Port Forward Rules Updated")
}
