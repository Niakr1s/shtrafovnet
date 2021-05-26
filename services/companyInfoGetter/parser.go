package companyInfoGetter

import (
	"fmt"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

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
