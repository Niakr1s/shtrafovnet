package companyInfoGetter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetCompanyPageLink(t *testing.T) {
	link, err := getCompanyPageLink("7734123457")
	assert.Nil(t, err)
	assert.Equal(t, "https://www.rusprofile.ru/id/9019334", link)
}
