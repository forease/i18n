package i18n

import (
	"fmt"
	"github.com/forease/config"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

var (
	Default  = "zh_CN"
	language = Default
	locales  = make(map[string]*config.Config)
)

func IsExists(lang string) bool {
	if _, ok := locales[lang]; ok {
		return true
	}

	return false
}

func getLang(str string) string {
	s := strings.Replace(str, "locale_", "", -1)
	s = strings.Replace(s, ".ini", "", -1)

	return s
}

// load locales
func LoadLocales(dir string) (err error) {
	f, err := os.Open(dir)
	if err != nil {
		return
	}

	names, err := f.Readdirnames(-1)
	if err != nil {
		return
	}
	for _, name := range names {
		localeFile := filepath.Join(dir, name)
		fio, err := os.Lstat(localeFile)
		if fio.IsDir() {
			continue
		}

		lang := getLang(name)
		fmt.Println(name, localeFile, lang)
		if locale, ok := locales[lang]; ok {
			err = locale.ReloadConfig()
			if err != nil {
				return err
			}
		} else {
			message, err := config.NewConfig(localeFile, 24)
			if err != nil {
				return err
			}

			locales[lang] = message
		}

	}

	return
}

func SetLanguage(lang string) {
	language = lang
}

// Tr translate content to target language.
func Tr(lang, format string, args ...interface{}) string {

	locale, ok := locales[lang]
	if ok {
		value, _ := locale.String(format, format)
		format = value
	}

	if len(args) > 0 {
		params := make([]interface{}, 0, len(args))
		for _, arg := range args {
			if arg != nil {
				val := reflect.ValueOf(arg)
				if val.Kind() == reflect.Slice {
					for i := 0; i < val.Len(); i++ {
						params = append(params, val.Index(i).Interface())
					}
				} else {
					params = append(params, arg)
				}
			}
		}
		return fmt.Sprintf(format, params...)
	}
	return fmt.Sprintf(format)
}
