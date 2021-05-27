package companyInfoGetter

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)

func Test_GetCompanyPageLink(t *testing.T) {
	link, err := getCompanyPageLink("7734123457")
	assert.Nil(t, err)
	assert.Equal(t, "https://www.rusprofile.ru/id/9019334", link)
}

func Test_ParseCompanyInfoPage(t *testing.T) {
	page := `
<div class="company-name" itemprop="legalName">ОБЩЕСТВО ОГРАНИЧЕННОЙ ОТВЕТСТВЕННОСТИ ОБЩЕСТВО</div>
<span class="copy_target" id="clip_inn">12341234</span>
<span class="copy_target" id="clip_kpp">43214321</span>
<span class="founder-item__title">Босс</span>
`
	info, err := parseCompanyInfoPage(strings.NewReader(page))

	assert.NoError(t, err)
	assert.Equal(t, CompanyInfo{Inn: "12341234", Kpp: "43214321", Name: "ОБЩЕСТВО ОГРАНИЧЕННОЙ ОТВЕТСТВЕННОСТИ ОБЩЕСТВО", Chief: "Босс"}, info)
}

func Test_GetTextFromNode(t *testing.T) {
	node := `<div class="some_node">Node1</div><div class="some_node">Node2</div>`

	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(node))

	text, err := getTextFromFirstNode(".some_node", doc)

	assert.NoError(t, err)
	assert.Equal(t, "Node1", text)
}
