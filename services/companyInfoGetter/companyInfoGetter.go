package companyInfoGetter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type Server struct {
	UnimplementedCompanyInfoGetterServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) GetCompanyInfo(ctx context.Context, in *GetCompanyInfoRequest) (*GetCompanyInfoResponse, error) {
	urlStr, err := getCompanyPageLink(in.Inn)
	if err != nil {
		return nil, err
	}

	log.Printf("trying to visit url %s", urlStr)

	r, err := http.Get(urlStr)
	if err != nil {
		return nil, err
	}

	body := r.Body
	defer body.Close()

	info, err := parseCompanyInfoPage(body)
	if err != nil {
		return nil, err
	}

	response := GetCompanyInfoResponse{
		Inn:   info.Inn,
		Kpp:   info.Kpp,
		Name:  info.Name,
		Chief: info.Chief,
	}

	return &response, nil
}

func getCompanyPageLink(inn string) (string, error) {
	ajaxQueryUrl := fmt.Sprintf(`https://www.rusprofile.ru/ajax.php?query=%s&action=search`, inn)
	r, err := http.Get(ajaxQueryUrl)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	r.Body.Close()

	type Reply struct {
		Ul []struct {
			Link string `json:"link"`
			Inn  string `json:"inn"`
			URL  string `json:"url"`
		} `json:"ul"`
		Success bool   `json:"success"`
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	reply := Reply{}
	if err := json.Unmarshal(body, &reply); err != nil {
		return "", err
	}

	log.Println(reply)

	if !reply.Success {
		return "", fmt.Errorf("unsuccesfull query: code = %d, message = %s", reply.Code, reply.Message)
	}

	var pageLink string

	for _, ul := range reply.Ul {
		if strings.Contains(ul.Inn, inn) {
			pageLink = "https://www.rusprofile.ru" + ul.Link
		}
	}

	if pageLink == "" {
		return "", fmt.Errorf("no company with inn = %s found", inn)
	}

	return pageLink, nil
}
