package companyInfoGetter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const baseUrl = "https://www.rusprofile.ru"

type Server struct {
	UnimplementedCompanyInfoGetterServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) GetCompanyInfo(ctx context.Context, in *GetCompanyInfoRequest) (*GetCompanyInfoResponse, error) {
	log.Printf("GetCompanyInfo: got request: %+v", in)

	urlStr, err := getCompanyPageLink(in.Inn)
	if err != nil {
		log.Printf("GetCompanyInfo: fail, request: %+v, err: %v", in, err)
		if err == errInnNotFound {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, err
	}

	r, err := http.Get(urlStr)
	if err != nil {
		log.Printf("GetCompanyInfo: fail, request: %+v, err: %v", in, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	body := r.Body
	defer body.Close()

	info, err := parseCompanyInfoPage(body)
	if err != nil {
		log.Printf("GetCompanyInfo: fail, request: %+v, err: %v", in, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := &GetCompanyInfoResponse{
		Inn:   info.Inn,
		Kpp:   info.Kpp,
		Name:  info.Name,
		Chief: info.Chief,
	}
	log.Printf("GetCompanyInfo: success, request: %+v, response: %v", in, response)

	return response, nil
}

var errInnNotFound = errors.New("inn not found")

// getCompanyPageLink returns errInnNotFound if company with such inn wasn't found
func getCompanyPageLink(inn string) (string, error) {
	ajaxQueryUrl := baseUrl + fmt.Sprintf(`/ajax.php?query=%s&action=search`, inn)
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

	if !reply.Success {
		return "", fmt.Errorf("unsuccesfull query: code = %d, message = %s", reply.Code, reply.Message)
	}

	var pageLink string

	for _, ul := range reply.Ul {
		if strings.Contains(ul.Inn, inn) {
			pageLink = baseUrl + ul.Link
		}
	}

	if pageLink == "" {
		return "", errInnNotFound
	}

	return pageLink, nil
}

type CompanyInfo struct {
	Inn   string
	Kpp   string
	Name  string
	Chief string
}

func parseCompanyInfoPage(pageReader io.Reader) (CompanyInfo, error) {
	res := CompanyInfo{}

	doc, err := goquery.NewDocumentFromReader(pageReader)
	if err != nil {
		return res, err
	}

	res.Inn, err = getTextFromFirstNode("#clip_inn", doc)
	if err != nil {
		return res, err
	}

	res.Kpp, err = getTextFromFirstNode("#clip_kpp", doc)
	if err != nil {
		return res, err
	}

	res.Name, err = getTextFromFirstNode(".company-name", doc)
	if err != nil {
		return res, err
	}

	res.Chief, err = getTextFromFirstNode(".founder-item__title", doc)
	if err != nil {
		return res, err
	}

	return res, nil
}

func getTextFromFirstNode(sel string, doc *goquery.Document) (string, error) {
	var res string

	s := doc.Find(sel)
	if s.Length() == 0 {
		return res, fmt.Errorf("node with selector %s not found", sel)
	}

	s.EachWithBreak(func(i int, s *goquery.Selection) bool {
		res = s.Text()
		return false
	})

	return strings.TrimSpace(res), nil
}
