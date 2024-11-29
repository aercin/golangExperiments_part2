package data_validation

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type echoDataValidator struct {
	validator *validator.Validate
}

func NewEchoDataValidator() *echoDataValidator {
	return &echoDataValidator{
		validator: validator.New(),
	}
}

func (cv *echoDataValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func bindPathParams(c echo.Context, dest interface{}) error {
	val := reflect.ValueOf(dest).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		paramTag := fieldType.Tag.Get("param")

		if paramTag != "" {
			// Path parametresini al
			paramValue := c.Param(paramTag)
			if paramValue == "" {
				continue
			}

			// Değerin türüne göre işlem yap
			switch field.Kind() {
			case reflect.String:
				field.SetString(paramValue)
			case reflect.Int, reflect.Int64:
				intValue, err := strconv.Atoi(paramValue)
				if err != nil {
					return err
				}
				field.SetInt(int64(intValue))
				// Diğer türler eklenebilir
			}
		}
	}
	return nil
}

// Middleware fonksiyonu
func ValidationMiddleware(structType interface{}) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Gelen struct türüne göre yeni bir instance oluştur
			req := reflect.New(reflect.TypeOf(structType).Elem()).Interface()

			// Path parametrelerini struct'a yansıt
			if err := bindPathParams(c, req); err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid path parameters")
			}

			// Body'den alınacak verileri struct'a yansıt
			bodyBytes, _ := io.ReadAll(c.Request().Body)
			// Body'yi tekrar kullanılabilir hale getir
			c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			if len(bodyBytes) != 0 {
				if err := json.Unmarshal(bodyBytes, req); err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
				}
			}

			// Validasyonu çalıştır
			if err := c.Validate(req); err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}

			// Bir sonraki handler'a geç
			return next(c)
		}
	}
}
