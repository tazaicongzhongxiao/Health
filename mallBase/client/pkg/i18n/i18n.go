// Package i18n i18n 国际化
//
//	example:
//	bundle := i18n.NewBundle(language.Chinese).LoadFiles("./locales")
//	fmt.Println(i18n.Get().NewPrinter(language.Chinese).Translate("Hello", i18n.Data{"name": "medivh", "count": 156}, "one"))
//	# return 你好世界
package i18n

import (
	"bytes"
	"fmt"
	"github.com/pelletier/go-toml"
	"golang.org/x/text/language"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
)

var (
	bundle = NewBundle(language.English)
)

// Data 自定义的需要被模板解析的信息, 在不需要传入结构体时可以使用此类型
type Data map[string]interface{}

// Bundle 国际化对外操作结构体
type Bundle struct {
	mu         sync.Mutex
	defaultTag language.Tag
	messages   map[language.Tag]map[string]map[string]*template.Template
}

// Messages 所有的模板信息
type Messages map[language.Tag]map[string]map[string]*template.Template

// Printer 一个新的message translate 对象, 一般不需要调用他
type Printer struct {
	messages   Messages
	acceptTags []language.Tag
	defaultTag language.Tag
}

// Get
// @Description: Config 得到config对象
// @return *Bundle
func Get() *Bundle {
	return bundle
}

// NewBundle 得到一个写的国际化实例
func NewBundle(tag language.Tag) *Bundle {
	return &Bundle{
		defaultTag: tag,
		messages:   make(Messages),
	}
}

// SetMessage add message
func (bundle *Bundle) SetMessage(tag language.Tag, key string, message map[string]string) {
	bundle.mu.Lock()
	defer bundle.mu.Unlock()
	bundle.messages[tag][key] = make(map[string]*template.Template)
	for name, val := range message {
		bundle.messages[tag][key][name] = createMessageTemplate(name, val)
	}
}

func createMessageTemplate(messageID, text string) *template.Template {
	if text == "" {
		return nil
	}
	t, err := template.New(messageID).Parse(text)
	if err != nil {
		return nil
	}
	return t
}

// LoadFiles walk file dir and load messages
// language file like
//   - path
//     | -- zh.toml
//     | -- en.toml
func (bundle *Bundle) LoadFiles(path string) *Bundle {
	_ = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		fileInfo := strings.Split(info.Name(), ".")
		lang := fileInfo[0]
		content, _ := ioutil.ReadFile(path)
		var data map[string]map[string]string
		err = toml.Unmarshal(content, &data)
		if err != nil {
			return err
		}
		tag := language.MustParse(lang)
		bundle.messages[tag] = make(map[string]map[string]*template.Template)
		for messageID, content := range data {
			bundle.SetMessage(tag, messageID, content)
		}
		return nil
	})
	return bundle
}

// NewPrinter 根据传入语言tag获得具体翻译组件
func (bundle *Bundle) NewPrinter(tags ...language.Tag) *Printer {
	return &Printer{acceptTags: tags, defaultTag: bundle.defaultTag, messages: bundle.messages}
}

func (p *Printer) getAcceptTag() language.Tag {
	for _, tag := range p.acceptTags {
		if _, ok := p.messages[tag]; ok {
			return tag
		}
	}
	return p.defaultTag
}

// Translate 根据传入模板ID进行翻译, data can be nil
func (p *Printer) Translate(key string, data interface{}, code string) string {
	var rs bytes.Buffer
	var err error
	tag := p.getAcceptTag()
	messages := p.messages[tag]
	message, ok := messages[key]
	if !ok {
		message = p.messages[p.defaultTag][key]
	}
	if message == nil {
		return key
	}
	if message[code] == nil {
		return fmt.Sprintf("%v:没有语言配置", code)
	}
	err = message[code].Execute(&rs, data)
	if err != nil {
		return fmt.Sprintf("翻译错误: %v", err)
	}
	content, _ := ioutil.ReadAll(&rs)
	return string(content)
}
