package taxcloud

import (
	"encoding/json"
	"errors"
)

type TIC struct {
	TICID int
	Description string
}

type GetTICsResponse struct {
	ResponseBase

	TICs []TIC
}

func (g *Gateway) GetTICs() (*GetTICsResponse, error) {
	responsePtr, err := g.httpPost("GetTICs", nil)
	if err != nil {
		return nil, err
	}

	var getTICsResponse GetTICsResponse
	err = json.Unmarshal(*responsePtr, &getTICsResponse)
	if err != nil {
		return nil, errors.New("51a262ea-89cc-4ebd-9db3-339467783034: " + err.Error())
	}

	if err = getTICsResponse.ResponseBase.check(*responsePtr); err != nil {
		return nil, err
	}

	return &getTICsResponse, nil
}

type TICGroup struct {
	GroupID int
	Description string
}

type GetTICGroupsResponse struct {
	ResponseBase

	TICGroups []TICGroup
}

func (g *Gateway) GetTICGroups() (*GetTICGroupsResponse, error) {
	responsePtr, err := g.httpPost("GetTICGroups", nil)
	if err != nil {
		return nil, err
	}

	var getTICGroupsResponse GetTICGroupsResponse
	err = json.Unmarshal(*responsePtr, &getTICGroupsResponse)
	if err != nil {
		return nil, errors.New("6dbb7558-5e70-4b82-b888-429cc3e50c7f: " + err.Error())
	}

	if err = getTICGroupsResponse.ResponseBase.check(responsePtr); err != nil {
		return nil, err
	}

	return &getTICGroupsResponse, nil
}

type GetTICsByGroupParams struct {
	GroupID int `json:"ticGroup"`
}

func (g *Gateway) GetTICsByGroup(getTICsByGroupParams GetTICsByGroupParams) (*GetTICsResponse, error) {
	responsePtr, err := g.httpPost("GetTICsByGroup", getTICsByGroupParams)
	if err != nil {
		return nil, err
	}

	var getTICsResponse GetTICsResponse
	err = json.Unmarshal(*responsePtr, &getTICsResponse)
	if err != nil {
		return nil, errors.New("bd86f72a-f980-435f-84ec-94bbb6f78da8: " + err.Error())
	}

	if err = getTICsResponse.ResponseBase.check(responsePtr); err != nil {
		return nil, err
	}

	return &getTICsResponse, nil
}







