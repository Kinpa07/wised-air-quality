package sensor_readings_collector_pkg

import (
	"github.com/SintroSecurity/go-libraries/router"
	"github.com/go-playground/validator/v10"
)

type ClientTypeEnum uint

const (
	ClientTypeEnumDeviceV1 = iota + 1
	ClientTypeEnumDeviceV2
)

type ClientType string

const (
	ClientTypeDeviceV1 ClientType = "V1"
	ClientTypeDeviceV2            = "V2"
)

func init() {
	router.AddCustomValidator("client_type", func(fl validator.FieldLevel) bool {
		return ClientType(fl.Field().String()).IsValid()
	})
}

func (g *ClientType) ToEnum() ClientTypeEnum {
	switch *g {
	case ClientTypeDeviceV1:
		return ClientTypeEnumDeviceV1
	case ClientTypeDeviceV2:
		return ClientTypeEnumDeviceV2
	}
	panic("invalid ClientType enum")
}

func (g ClientType) IsValid() bool {
	return g == ClientTypeDeviceV1 || g == ClientTypeDeviceV2
}

func (g *ClientTypeEnum) ToString() ClientType {
	switch *g {
	case ClientTypeEnumDeviceV1:
		return ClientTypeDeviceV1
	case ClientTypeEnumDeviceV2:
		return ClientTypeDeviceV2
	}
	panic("invalid ClientType enum")
}
