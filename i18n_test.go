package i18n

import (
	//"errors"
	//"flag"
	//"strings"
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {

	err := SetMessage("zh-CN", "locale_zh-CN.ini")
	if err != nil {
		fmt.Println(err)
	}

	lists := ListLangs()
	fmt.Println(lists)

	lang := new(Locale)
	lang.Lang = "zh-CN"
	fmt.Println(lang.Tr("admin.name"))
	fmt.Println(lang.Tr("admin.users.name"))
	fmt.Println(lang.Tr("admin.users.test.name"))
	fmt.Println(lang.Tr("admin.users.your.name", "forease"))

}
