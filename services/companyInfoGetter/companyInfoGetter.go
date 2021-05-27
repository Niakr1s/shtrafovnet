package companyInfoGetter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
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
