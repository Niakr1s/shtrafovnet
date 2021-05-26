package companyInfoGetter

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	UnimplementedCompanyInfoGetterServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) GetCompanyInfo(ctx context.Context, in *GetCompanyInfoRequest) (*GetCompanyInfoResponse, error) {
	urlStr := fmt.Sprintf("https://www.rusprofile.ru/search?query=%s", in.Inn)
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
