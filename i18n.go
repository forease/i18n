// Copyright 2013 beego authors
// Copyright 2015 forease authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

// Package i18n is for app Internationalization and Localization.
package i18n

import (
	"fmt"
	"github.com/forease/config"
	"reflect"
)

var (
	locales = &localeStore{store: make(map[string]*locale)}
)

type locale struct {
	id      int
	lang    string
	message *config.Config
}

type localeStore struct {
	langs []string
	store map[string]*locale
}

// Get locale from localeStore use specify lang string
func (d *localeStore) getLocale(lang string) (*locale, bool) {
	for _, l := range d.store {
		if l.lang == lang {
			return l, true
		}
	}
	return nil, false
}

// Get target language string
func (d *localeStore) Get(lang, format string) (string, bool) {
	if locale, ok := d.getLocale(lang); ok {
		value, err := locale.message.String(format, format)
		if err == nil {
			return value, true
		} else {
			fmt.Println(err)
		}

	}
	return "", false
}

func (d *localeStore) Add(lc *locale) bool {
	if _, ok := d.store[lc.lang]; ok {
		return false
	}
	lc.id = len(d.langs)
	d.langs = append(d.langs, lc.lang)
	d.store[lc.lang] = lc
	return true
}

func (d *localeStore) Reload(langs ...string) error {
	if len(langs) == 0 {
		for _, lc := range d.store {
			err := lc.message.ReloadConfig()
			if err != nil {
				return err
			}
		}
	} else {
		for _, lang := range langs {
			if lc, ok := d.getLocale(lang); ok {
				err := lc.message.ReloadConfig()
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// Reload locales
func ReloadLangs(langs ...string) error {
	return locales.Reload(langs...)
}

// List all locale languages
func ListLangs() []string {
	langs := make([]string, len(locales.langs))
	copy(langs, locales.langs)
	return langs
}

// Check language name if exist
func IsExist(lang string) bool {
	_, ok := locales.store[lang]
	return ok
}

// Check language name if exist
func IndexLang(lang string) int {
	if lc, ok := locales.store[lang]; ok {
		return lc.id
	}
	return -1
}

// Get language by index id
func GetLangByIndex(index int) string {
	if index < 0 || index >= len(locales.langs) {
		return ""
	}
	return locales.langs[index]
}

// SetMessage sets the message file for localization.
func SetMessage(lang, filePath string) error {
	message, err := config.NewConfig(filePath, 24)
	if err == nil {
		lc := new(locale)
		lc.lang = lang
		lc.message = message
		if locales.Add(lc) == false {
			return fmt.Errorf("Lang %s alread exist", lang)
		}
	}
	return err
}

// A Locale describles the information of localization.
type Locale struct {
	Lang string
}

// Tr translate content to target language.
func (l Locale) Tr(format string, args ...interface{}) string {
	return Tr(l.Lang, format, args...)
}

// Index get lang index of LangStore
func (l Locale) Index() int {
	return IndexLang(l.Lang)
}

// Tr translate content to target language.
func Tr(lang, format string, args ...interface{}) string {

	value, ok := locales.Get(lang, format)
	if ok {
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
