package i18n

import (
	"fmt"
	"testing"
)

func Test_i18n(t *testing.T) {
	err := LoadLocales("locales")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("trans:", Tr("zh-CN", "admin.users.your.name", "Google"))
}
