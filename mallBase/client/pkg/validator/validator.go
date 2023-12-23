package validator

import (
	"MyTestMall/mallBase/basics/pkg/app"
	"MyTestMall/mallBase/basics/pkg/config"
	"MyTestMall/mallBase/basics/pkg/log"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	translationsEn "github.com/go-playground/validator/v10/translations/en"
	translationsZh "github.com/go-playground/validator/v10/translations/zh"
	json "github.com/json-iterator/go"
	"reflect"
	"strings"
	"sync"
)

// Validator 验证器
type Validator struct {
	once     sync.Once
	validate *validator.Validate
}

var (
	_                   binding.StructValidator = &Validator{}
	validatorMessages   map[string]map[string]string
	langMapping         = map[string]string{"zh-cn": "zh", "en-us": "en"}
	enUs                = en.New()
	zhCn                = zh.New()
	universalTranslator = ut.New(enUs, zhCn)
	zhTrans, _          = universalTranslator.GetTranslator(zhCn.Locale())
	enTrans, _          = universalTranslator.GetTranslator(enUs.Locale())
)

// ValidateStruct 验证结构体
func (v *Validator) ValidateStruct(obj interface{}) error {
	if kindOfData(obj) == reflect.Struct {
		v.lazyInit()
		if err := v.validate.Struct(obj); err != nil {
			return err
		}
	}
	return nil
}

// Engine 获取验证器
func (v *Validator) Engine() interface{} {
	v.lazyInit()
	return v.validate
}

// lazyInit 延迟初始化
func (v *Validator) lazyInit() {
	v.once.Do(func() {
		v.validate = validator.New()
		// v.validate.SetTagName("binding")
		v.validate.RegisterTagNameFunc(func(field reflect.StructField) string {
			name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
		_ = translationsZh.RegisterDefaultTranslations(v.validate, zhTrans)
		_ = translationsEn.RegisterDefaultTranslations(v.validate, enTrans)
		if err := config.Config().Bind("validator-messages.toml", "", &validatorMessages, nil); err != nil {
			log.Error("验证配置读取错误", err.Error())
			panic(err)
		}
		for tag, languages := range validatorMessages {
			var trueTag, messages = tag, languages
			registerTranslation(v.validate, trueTag, messages)
		}
	})
}

func registerTranslation(validate *validator.Validate, tag string, languages map[string]string) {
	_ = validate.RegisterTranslation(tag, zhTrans, func(ut ut.Translator) error {
		return ut.Add(tag, languages["zh-cn"], true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(tag, fe.Param(), fe.Field())
		return t
	})
	_ = validate.RegisterTranslation(tag, enTrans, func(ut ut.Translator) error {
		return ut.Add(tag, languages["en-us"], true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(tag, fe.Param(), fe.Field())
		return t
	})
}

// ValidErrors 验证之后的错误信息
type ValidErrors struct {
	ErrorsInfo map[string]string
	triggered  bool
}

func (validErrors *ValidErrors) add(key, value string) {
	validErrors.ErrorsInfo[key] = value
	validErrors.triggered = true
}

// IsValid 是否验证成功
func (validErrors *ValidErrors) IsValid() bool {
	return !validErrors.triggered
}

func newValidErrors() *ValidErrors {
	return &ValidErrors{
		triggered:  false,
		ErrorsInfo: make(map[string]string),
	}
}

// Get
// @Description: 获取验证码对象
// @return *validator.Validate
// errs := validator.Get().Var(form.User, "required,email,gt=1,lt=10")
// fmt.Println(validator.GetValidErrors(ctx, errs))
func Get() *validator.Validate {
	v, _ := binding.Validator.Engine().(*validator.Validate)
	return v
}

func GetVal(c *gin.Context, field interface{}, tag string) ValidErrors {
	if err := Get().Var(field, tag); err != nil {
		return GetValidErrors(c, err)
	}
	return ValidErrors{}
}

func GetLang(c *gin.Context) string {
	return c.DefaultQuery("lang", "zh-cn")
}

// ValidateMap 验证Map
func ValidateMap(lang string, data map[string]interface{}, rules map[string]interface{}) (errors ValidErrors) {
	defer func() {
		if r := recover(); r != nil {
			errors = *newValidErrors()
			errors.add("err", fmt.Sprintf("%v", r))
			return
		}
	}()
	if lang == "" {
		lang = "zh-cn"
	}
	if errs := Get().ValidateMap(data, rules); errs != nil {
		trans, _ := universalTranslator.GetTranslator(langMapping[lang])
		errors = *newValidErrors()
		for key, err := range errs {
			if errs1, ok := err.(validator.ValidationErrors); ok {
				for _, value := range errs1 {
					errors.add(key, value.Translate(trans))
				}
			} else {
				errors.add(key, errs1.Error())
			}
		}
		return errors
	}
	return errors
}

// GetValidErrors
// @Description: 转换错误
// @param c
// @param err
// @return errors
func GetValidErrors(c *gin.Context, err error) (errors ValidErrors) {
	if err != nil {
		trans, _ := universalTranslator.GetTranslator(langMapping[GetLang(c)])
		errors = *newValidErrors()
		if errs, ok := err.(validator.ValidationErrors); ok {
			for _, value := range errs {
				errors.add(value.Field(), value.Translate(trans))
			}
		} else {
			errors.add("err", err.Error())
		}
	}
	return errors
}

// Bind 自定义错误信息, 如果没有匹配需要在 configs/validator-messages.yaml 中添加对应处理数据
func Bind(c *gin.Context, param interface{}) ValidErrors {
	var err error
	err = binding.JSON.Bind(c.Request, param)
	if err != nil && c.Request.Method == "DELETE" {
		var t interface{}
		err = app.Unmarshal(c.Request.URL.Query(), &t)
		if err == nil {
			jsons, _ := json.Marshal(t)
			err = binding.JSON.BindBody(jsons, param)
		}
	}
	if err != nil {
		return GetValidErrors(c, err)
	}
	return ValidErrors{}
}

func BindQuery(c *gin.Context, param interface{}) ValidErrors {
	if err := c.ShouldBindQuery(param); err != nil {
		return GetValidErrors(c, err)
	}
	return ValidErrors{}
}

func kindOfData(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()
	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}
