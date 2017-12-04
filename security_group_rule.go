package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// GroupRule definition
type GroupRule struct {
	Direction    string `json:"direction"`
	Protocol     string `json:"protocol"`
	IPRange      string `json:"ip_range"`
	DestPortFrom int    `json:"dest_port_from,omitempty"`
	Action       string `json:"action"`
	Position     int    `json:"position"`
	DestPortTo   string `json:"dest_port_to"`
	Editable     bool   `json:"editable"`
	ID           string `json:"id"`
}

// GetGroupRules represents the response of a GET /_group/{groupID}/rules
type GetGroupRules struct {
	Rules []GroupRule `json:"rules"`
}

// GetGroupRule represents the response of a GET /_group/{groupID}/rules/{ruleID}
type GetGroupRule struct {
	Rules GroupRule `json:"rule"`
}

// NewGroupRule definition POST/PUT request /_group/{groupID}
type NewGroupRule struct {
	Action       string `json:"action"`
	Direction    string `json:"direction"`
	IPRange      string `json:"ip_range"`
	Protocol     string `json:"protocol"`
	DestPortFrom int    `json:"dest_port_from,omitempty"`
}

// GetGroupRules returns a GroupRules
func (s *API) GetGroupRules(groupID string) (*GetGroupRules, error) {
	resp, err := s.GetResponsePaginate(s.computeAPI, fmt.Sprintf("_groups/%s/rules", groupID), url.Values{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}
	var GroupRules GetGroupRules

	if err = json.Unmarshal(body, &GroupRules); err != nil {
		return nil, err
	}
	return &GroupRules, nil
}

// GetAGroupRule returns a GroupRule
func (s *API) GetAGroupRule(groupID string, rulesID string) (*GetGroupRule, error) {
	resp, err := s.GetResponsePaginate(s.computeAPI, fmt.Sprintf("_groups/%s/rules/%s", groupID, rulesID), url.Values{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}
	var GroupRules GetGroupRule

	if err = json.Unmarshal(body, &GroupRules); err != nil {
		return nil, err
	}
	return &GroupRules, nil
}

type postGroupRuleResponse struct {
	GroupRule GroupRule `json:"rule"`
}

// PostGroupRule posts a rule on a server
func (s *API) PostGroupRule(GroupID string, rules NewGroupRule) (*GroupRule, error) {
	resp, err := s.PostResponse(s.computeAPI, fmt.Sprintf("_groups/%s/rules", GroupID), rules)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := s.handleHTTPError([]int{http.StatusCreated}, resp)
	if err != nil {
		return nil, err
	}
	var res postGroupRuleResponse
	err = json.Unmarshal(data, &res)
	return &res.GroupRule, err
}

// PutGroupRule updates a GroupRule
func (s *API) PutGroupRule(rules NewGroupRule, GroupID, RuleID string) error {
	resp, err := s.PutResponse(s.computeAPI, fmt.Sprintf("_groups/%s/rules/%s", GroupID, RuleID), rules)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = s.handleHTTPError([]int{http.StatusOK}, resp)
	return err
}

// DeleteGroupRule deletes a GroupRule
func (s *API) DeleteGroupRule(GroupID, RuleID string) error {
	resp, err := s.DeleteResponse(s.computeAPI, fmt.Sprintf("_groups/%s/rules/%s", GroupID, RuleID))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = s.handleHTTPError([]int{http.StatusNoContent}, resp)
	return err
}
