package config

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type HTTPConfig struct {
	Port                int        `property:"port,default=8080"`
	ReadTimeoutSeconds  int        `property:"readTimeout,default=5"`
	WriteTimeoutSeconds int        `property:"writeTimeout,default=10"`
	IdleTimeoutSeconds  int        `property:"idleTimeout,default=120"`
	MaxHeaderBytes      int        `property:"maxHeaderBytes,default=1048576"`
	TLS                 TLSConfig  `property:"tls"`
	CORS                CORSConfig `property:"cors"`
}

type CORSConfig struct {
	AllowedOrigins   []string `property:"allowedOrigins"`
	AllowedMethods   []string `property:"allowedMethods,default=['GET','POST','PUT','DELETE','OPTIONS']"`
	AllowedHeaders   []string `property:"allowedHeaders,default=['Authorization','Content-Type']"`
	AllowCredentials bool     `property:"allowCredentials,default=false"`
	MaxAgeSeconds    int      `property:"maxAgeSeconds,default=600"`
}

type TLSConfig struct {
	Enabled    bool   `property:"enabled,default=false"`
	CertFile   string `property:"certFile"`
	KeyFile    string `property:"keyFile"`
	MinVersion string `property:"minVersion,default='TLS1.2'"`
}

type DatabaseConfig struct {
	Host     string  `property:"host,default='localhost'"`
	Port     int     `property:"port,default=5432"`
	Username string  `property:"username"`
	Password string  `property:"password"`
	Timeout  float64 `property:"timeout,default=30.0"`

	NoPropertyTag string
	unExported    string
}

type ServiceConfig struct {
	Name   string            `property:"name"`
	Labels map[string]string `property:"labels,default={'env':'prod','team':'platform'}"`
}

type ServiceConfigWithInvalidMapDefault struct {
	Name   string         `property:"name"`
	Limits map[string]int `property:"labels,default={'cpu': 'high'}"`
}

type ServiceConfigWithInvalidSliceDefault struct {
	Name  string `property:"name"`
	Ports []int  `property:"ports,default=['8080','9090','notAnInt']"`
}

type ConfigWithInvalidPropertyTag struct {
	Labels map[string]string `property:"labels,default={env:'prod',team:'platform'}"`
}

type ConfigWithInvalidDefaultValue struct {
	Timeout int `property:"timeout,default='thirty'"`
}

func TestNewDefaultBinder(t *testing.T) {
	testCases := []struct {
		name            string
		propertySources *PropertySources
		wantPanic       error
	}{
		{
			name:            "nil property sources",
			propertySources: nil,
			wantPanic:       errors.New("nil property sources"),
		},
		{
			name:            "with property sources",
			propertySources: NewPropertySources(NewMapPropertySource("anyMapName", map[string]any{"anyKey": "anyValue"})),
			wantPanic:       nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given

			// when
			if tc.wantPanic != nil {
				require.PanicsWithValue(t, tc.wantPanic.Error(), func() {
					NewDefaultPropertyBinder(tc.propertySources)
				})
				return
			}

			resolver := NewDefaultPropertyBinder(tc.propertySources)

			// then
			require.NotNil(t, resolver)
		})
	}
}

func TestDefaultBinder_Bind(t *testing.T) {
	testCases := []struct {
		name       string
		propSource PropertySource

		propName    string
		targetType  reflect.Type
		targetValue any

		wantErr    error
		wantResult any
	}{
		{
			name: "nil target",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"app.name": "Procyon",
			}),
			wantErr: errors.New("nil target"),
		},
		{
			name: "non-pointer target",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"app.name": "Procyon",
			}),
			propName:    "app.name",
			targetValue: "notAPointer",
			wantErr:     errors.New("target must be a non-nil pointer"),
		},
		{
			name: "empty property name",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"app.name": "Procyon",
			}),
			propName:   "",
			targetType: reflect.TypeFor[string](),
			wantErr:    errors.New("empty or blank property name"),
		},
		{
			name: "blank property name",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"app.name": "Procyon",
			}),
			propName:   " ",
			targetType: reflect.TypeFor[string](),
			wantErr:    errors.New("empty or blank property name"),
		},
		{
			name: "any property not found",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"app.name": "Procyon",
			}),
			propName:   "app.version",
			targetType: reflect.TypeFor[any](),
			wantErr:    ErrNoPropertyFound,
		},
		{
			name: "string property",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"app.name": "Procyon",
			}),
			propName:   "app.name",
			targetType: reflect.TypeFor[string](),
			wantResult: "Procyon",
		},
		{
			name: "bool property from bool source",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"app.debug": true,
			}),
			propName:   "app.debug",
			targetType: reflect.TypeFor[bool](),
			wantResult: true,
		},
		{
			name: "bool property from string source",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"app.debug": "true",
			}),
			propName:   "app.debug",
			targetType: reflect.TypeFor[bool](),
			wantResult: true,
		},
		{
			name: "invalid bool property",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"app.debug": 8080,
			}),
			propName:   "app.debug",
			targetType: reflect.TypeFor[bool](),
			wantErr:    errors.New("strconv.ParseBool: parsing \"8080\": invalid syntax"),
		},
		{
			name: "int property from int source",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"app.port": 8080,
			}),
			propName:   "app.port",
			targetType: reflect.TypeFor[int](),
			wantResult: 8080,
		},
		{
			name: "int property from int8 source",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"app.port": int8(127),
			}),
			propName:   "app.port",
			targetType: reflect.TypeFor[int](),
			wantResult: 127,
		},
		{
			name: "int property from int16 source",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"app.port": int16(8080),
			}),
			propName:   "app.port",
			targetType: reflect.TypeFor[int](),
			wantResult: 8080,
		},
		{
			name: "int property from int32 source",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"app.port": int32(8080),
			}),
			propName:   "app.port",
			targetType: reflect.TypeFor[int](),
			wantResult: 8080,
		},
		{
			name: "int property from int64 source",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"app.port": int64(8080),
			}),
			propName:   "app.port",
			targetType: reflect.TypeFor[int](),
			wantResult: 8080,
		},
		{
			name: "int property from string source",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"app.port": "8080",
			}),
			propName:   "app.port",
			targetType: reflect.TypeFor[int](),
			wantResult: 8080,
		},
		{
			name: "int property from float32 source",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"app.port": float32(8080),
			}),
			propName:   "app.port",
			targetType: reflect.TypeFor[int](),
			wantResult: 8080,
		},
		{
			name: "int property from float64 source",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"app.port": float64(8080),
			}),
			propName:   "app.port",
			targetType: reflect.TypeFor[int](),
			wantResult: 8080,
		},
		{
			name: "int property from uint8 source",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"app.port": uint8(127),
			}),
			propName:   "app.port",
			targetType: reflect.TypeFor[int](),
			wantResult: 127,
		},
		{
			name: "int property from uint16 source",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"app.port": uint16(8080),
			}),
			propName:   "app.port",
			targetType: reflect.TypeFor[int](),
			wantResult: 8080,
		},
		{
			name: "int property from uint32 source",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"app.port": uint32(8080),
			}),
			propName:   "app.port",
			targetType: reflect.TypeFor[int](),
			wantResult: 8080,
		},
		{
			name: "int property from uint64 source",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"app.port": uint64(8080),
			}),
			propName:   "app.port",
			targetType: reflect.TypeFor[int](),
			wantResult: 8080,
		},
		{
			name: "invalid int property",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"weather.temperature": false,
			}),
			propName:   "weather.temperature",
			targetType: reflect.TypeFor[int](),
			wantErr:    fmt.Errorf("strconv.ParseInt: parsing \"false\": invalid syntax"),
		},
		{
			name: "uint property from uint source",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"app.port": 8080,
			}),
			propName:   "app.port",
			targetType: reflect.TypeFor[uint](),
			wantResult: uint(8080),
		},
		{
			name: "uint property from string source",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"app.port": "8080",
			}),
			propName:   "app.port",
			targetType: reflect.TypeFor[uint](),
			wantResult: uint(8080),
		},
		{
			name: "uint property from int8 source",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"app.port": int8(127),
			}),
			propName:   "app.port",
			targetType: reflect.TypeFor[uint](),
			wantResult: uint(127),
		},
		{
			name: "uint property from int16 source",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"app.port": int16(8080),
			}),
			propName:   "app.port",
			targetType: reflect.TypeFor[uint](),
			wantResult: uint(8080),
		},
		{
			name: "uint property from int32 source",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"app.port": int32(8080),
			}),
			propName:   "app.port",
			targetType: reflect.TypeFor[uint](),
			wantResult: uint(8080),
		},
		{
			name: "uint property from int64 source",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"app.port": int64(8080),
			}),
			propName:   "app.port",
			targetType: reflect.TypeFor[uint](),
			wantResult: uint(8080),
		},
		{
			name: "uint property from uint8 source",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"app.port": uint8(127),
			}),
			propName:   "app.port",
			targetType: reflect.TypeFor[uint](),
			wantResult: uint(127),
		},
		{
			name: "uint property from uint16 source",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"app.port": uint16(8080),
			}),
			propName:   "app.port",
			targetType: reflect.TypeFor[uint](),
			wantResult: uint(8080),
		},
		{
			name: "uint property from uint32 source",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"app.port": uint32(8080),
			}),
			propName:   "app.port",
			targetType: reflect.TypeFor[uint](),
			wantResult: uint(8080),
		},
		{
			name: "uint property from uint64 source",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"app.port": uint64(8080),
			}),
			propName:   "app.port",
			targetType: reflect.TypeFor[uint](),
			wantResult: uint(8080),
		},
		{
			name: "uint property from float32 source",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"app.port": float32(8080),
			}),
			propName:   "app.port",
			targetType: reflect.TypeFor[uint](),
			wantResult: uint(8080),
		},
		{
			name: "uint property from float64 source",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"app.port": float64(8080),
			}),
			propName:   "app.port",
			targetType: reflect.TypeFor[uint](),
			wantResult: uint(8080),
		},
		{
			name: "uint property from negative int source",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"app.port": -8080,
			}),
			propName:   "app.port",
			targetType: reflect.TypeFor[uint](),
			wantErr:    errors.New("cannot convert negative integer -8080 to uint"),
		},
		{
			name: "uint property from negative float source",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"app.port": float32(-8080),
			}),
			propName:   "app.port",
			targetType: reflect.TypeFor[uint](),
			wantErr:    errors.New("cannot convert negative float -8080.000000 to uint"),
		},
		{
			name: "invalid uint property",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"weather.temperature": false,
			}),
			propName:   "weather.temperature",
			targetType: reflect.TypeFor[uint](),
			wantErr:    fmt.Errorf("strconv.ParseUint: parsing \"false\": invalid syntax"),
		},
		{
			name: "float32 property from float source",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"weather.temperature": 23.5,
			}),
			propName:   "weather.temperature",
			targetType: reflect.TypeFor[float32](),
			wantResult: float32(23.5),
		},
		{
			name: "float32 property from string source",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"weather.temperature": "23.5",
			}),
			propName:   "weather.temperature",
			targetType: reflect.TypeFor[float32](),
			wantResult: float32(23.5),
		},
		{
			name: "float32 property from int8 source",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"weather.temperature": int8(23),
			}),
			propName:   "weather.temperature",
			targetType: reflect.TypeFor[float32](),
			wantResult: float32(23),
		},
		{
			name: "float32 property from int16 source",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"weather.temperature": int16(23),
			}),
			propName:   "weather.temperature",
			targetType: reflect.TypeFor[float32](),
			wantResult: float32(23),
		},
		{
			name: "float32 property from int32 source",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"weather.temperature": int32(23),
			}),
			propName:   "weather.temperature",
			targetType: reflect.TypeFor[float32](),
			wantResult: float32(23),
		},
		{
			name: "float32 property from int64 source",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"weather.temperature": int64(23),
			}),
			propName:   "weather.temperature",
			targetType: reflect.TypeFor[float32](),
			wantResult: float32(23),
		},
		{
			name: "float32 property from int source",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"weather.temperature": 23,
			}),
			propName:   "weather.temperature",
			targetType: reflect.TypeFor[float32](),
			wantResult: float32(23),
		},
		{
			name: "float32 property from uint8 source",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"weather.temperature": uint8(23),
			}),
			propName:   "weather.temperature",
			targetType: reflect.TypeFor[float32](),
			wantResult: float32(23),
		},
		{
			name: "float32 property from uint16 source",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"weather.temperature": uint16(23),
			}),
			propName:   "weather.temperature",
			targetType: reflect.TypeFor[float32](),
			wantResult: float32(23),
		},
		{
			name: "float32 property from uint32 source",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"weather.temperature": uint32(23),
			}),
			propName:   "weather.temperature",
			targetType: reflect.TypeFor[float32](),
			wantResult: float32(23),
		},
		{
			name: "float32 property from uint64 source",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"weather.temperature": uint64(23),
			}),
			propName:   "weather.temperature",
			targetType: reflect.TypeFor[float32](),
			wantResult: float32(23),
		},
		{
			name: "float32 property from int source",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"weather.temperature": 23,
			}),
			propName:   "weather.temperature",
			targetType: reflect.TypeFor[float32](),
			wantResult: float32(23),
		},
		{
			name: "invalid float32 property",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"weather.temperature": false,
			}),
			propName:   "weather.temperature",
			targetType: reflect.TypeFor[float32](),
			wantErr:    fmt.Errorf("strconv.ParseFloat: parsing \"false\": invalid syntax"),
		},
		{
			name: "float64 property from float source",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"weather.temperature": 23.5,
			}),
			propName:   "weather.temperature",
			targetType: reflect.TypeFor[float64](),
			wantResult: 23.5,
		},
		{
			name: "float64 property from string source",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"weather.temperature": "23.5",
			}),
			propName:   "weather.temperature",
			targetType: reflect.TypeFor[float64](),
			wantResult: 23.5,
		},
		{
			name: "map property",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"database.host":     "localhost",
				"database.port":     5432,
				"database.username": "admin",
				"database.password": "secret",
			}),
			propName:   "database",
			targetType: reflect.TypeFor[map[string]any](),
			wantResult: map[string]any{
				"host":     "localhost",
				"port":     5432,
				"username": "admin",
				"password": "secret",
			},
		},
		{
			name: "string slice property",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"servers[0]": "server1",
				"servers[1]": "server2",
				"servers[2]": "server3",
			}),
			propName:   "servers",
			targetType: reflect.TypeFor[[]string](),
			wantResult: []string{"server1", "server2", "server3"},
		},
		{
			name: "int slice property from string sources",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"ports[0]": "8080",
				"ports[1]": "9090",
				"ports[2]": "10080",
			}),
			propName:   "ports",
			targetType: reflect.TypeFor[[]int](),
			wantResult: []int{8080, 9090, 10080},
		},
		{
			name: "int slice property from string source with comma separation",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"ports": "8080, 9090, 10080",
			}),
			propName:   "ports",
			targetType: reflect.TypeFor[[]int](),
			wantResult: []int{8080, 9090, 10080},
		},
		{
			name: "invalid slice source with comma separation",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"ports": "true, false",
			}),
			propName:   "ports",
			targetType: reflect.TypeFor[[]int](),
			wantErr:    errors.New("cannot append element \"true\" to slice []int: strconv.ParseInt: parsing \"true\": invalid syntax"),
		},
		{
			name: "invalid slice source",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"ports": []any{true, false},
			}),
			propName:   "ports",
			targetType: reflect.TypeFor[[]int](),
			wantErr:    errors.New("cannot append property \"ports\" to slice []int: strconv.ParseInt: parsing \"true\": invalid syntax"),
		},
		{
			name: "invalid map source",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"ports": map[string]any{"first": true, "second": false},
			}),
			propName:   "ports",
			targetType: reflect.TypeFor[map[string]int](),
			wantErr:    errors.New("cannot bind map property \"ports\": strconv.ParseInt: parsing \"true\": invalid syntax"),
		},
		{
			name: "no values for map property",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"ports": map[string]any{},
			}),
			propName:   "ports",
			targetType: reflect.TypeFor[map[string]any](),
			wantResult: map[string]any{},
		},
		{
			name: "struct property",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"database.host":     "localhost",
				"database.port":     5432,
				"database.username": "admin",
				"database.password": "secret",
				"database.timeout":  30.5,
			}),
			propName:   "database",
			targetType: reflect.TypeFor[DatabaseConfig](),
			wantResult: DatabaseConfig{
				Host:     "localhost",
				Port:     5432,
				Username: "admin",
				Password: "secret",
				Timeout:  30.5,
			},
		},
		{
			name: "nested struct property",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"http.port":                   8443,
				"http.readTimeout":            10,
				"http.writeTimeout":           15,
				"http.idleTimeout":            180,
				"http.maxHeaderBytes":         2097152,
				"http.tls.enabled":            true,
				"http.tls.certFile":           "/path/to/cert.pem",
				"http.tls.keyFile":            "/path/to/key.pem",
				"http.tls.minVersion":         "TLS1.3",
				"http.cors.allowedOrigins[0]": "https://example.com",
			}),
			propName:   "http",
			targetType: reflect.TypeFor[HTTPConfig](),
			wantResult: HTTPConfig{
				Port:                8443,
				ReadTimeoutSeconds:  10,
				WriteTimeoutSeconds: 15,
				IdleTimeoutSeconds:  180,
				MaxHeaderBytes:      2097152,
				TLS: TLSConfig{
					Enabled:    true,
					CertFile:   "/path/to/cert.pem",
					KeyFile:    "/path/to/key.pem",
					MinVersion: "TLS1.3",
				},
				CORS: CORSConfig{
					AllowedOrigins:   []string{"https://example.com"},
					AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
					AllowedHeaders:   []string{"Authorization", "Content-Type"},
					AllowCredentials: false,
					MaxAgeSeconds:    600,
				},
			},
		},
		{
			name: "missing required property for struct",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"database.host":     "localhost",
				"database.port":     5432,
				"database.username": "admin",
				"database.timeout":  30.5,
			}),
			propName:   "database",
			targetType: reflect.TypeFor[DatabaseConfig](),
			wantErr:    errors.New("missing required property: \"database.password\""),
		},
		{
			name: "default value for struct",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"database.username": "admin",
				"database.password": "secret",
				"database.timeout":  30.5,
			}),
			propName:   "database",
			targetType: reflect.TypeFor[DatabaseConfig](),
			wantResult: DatabaseConfig{
				Host:     "localhost",
				Port:     5432,
				Username: "admin",
				Password: "secret",
				Timeout:  30.5,
			},
		},
		{
			name: "struct property with invalid value",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"database.host":     "localhost",
				"database.port":     5432,
				"database.username": "admin",
				"database.password": "secret",
				"database.timeout":  true,
			}),
			propName:   "database",
			targetType: reflect.TypeFor[DatabaseConfig](),
			wantErr:    errors.New("failed to bind property \"database.timeout\": strconv.ParseFloat: parsing \"true\": invalid syntax"),
		},
		{
			name: "unsupported target type",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"app.name": "Procyon",
			}),
			propName:   "app.name",
			targetType: reflect.TypeFor[chan int](),
			wantErr:    errors.New("unsupported target type: chan"),
		},
		{
			name: "struct property with map field",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"service.name":   "UserService",
				"service.labels": map[string]any{"env": "staging", "team": "backend"},
			}),
			propName:   "service",
			targetType: reflect.TypeFor[ServiceConfig](),
			wantResult: ServiceConfig{
				Name: "UserService",
				Labels: map[string]string{
					"env":  "staging",
					"team": "backend",
				},
			},
		},
		{
			name: "default value for map",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"service.name": "UserService",
			}),
			propName:   "service",
			targetType: reflect.TypeFor[ServiceConfig](),
			wantResult: ServiceConfig{
				Name: "UserService",
				Labels: map[string]string{
					"env":  "prod",
					"team": "platform",
				},
			},
		},
		{
			name: "struct property with invalid default map value",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"service.name": "UserService",
			}),
			propName:   "service",
			targetType: reflect.TypeFor[ServiceConfigWithInvalidMapDefault](),
			wantErr:    errors.New("failed to set default value for property \"service.labels\": cannot convert map value (string) to int: strconv.ParseInt: parsing \"high\": invalid syntax"),
		},
		{
			name: "struct property with invalid property tag",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"config.labels": map[string]any{"env": "prod"},
			}),
			propName:   "config",
			targetType: reflect.TypeFor[ConfigWithInvalidPropertyTag](),
			wantErr:    errors.New("failed to parse tags for field Labels: tag: failed to parse 'property' options \"labels,default={env:'prod',team:'platform'}\": failed to set option \"default\": invalid key: must start with '"),
		},
		{
			name:       "struct property with invalid default value",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{}),
			propName:   "config",
			targetType: reflect.TypeFor[ConfigWithInvalidDefaultValue](),
			wantErr:    errors.New("failed to set default value for property \"config.timeout\": strconv.ParseInt: parsing \"thirty\": invalid syntax"),
		},
		{
			name: "struct property with invalid default slice value",
			propSource: NewMapPropertySource("anyMapSource", map[string]any{
				"service.name": "UserService",
			}),
			propName:   "service",
			targetType: reflect.TypeFor[ServiceConfigWithInvalidSliceDefault](),
			wantErr:    errors.New("failed to set default value for property \"service.ports\": cannot append element \"notAnInt\" to slice []int: strconv.ParseInt: parsing \"notAnInt\": invalid syntax"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			propSources := NewPropertySources(tc.propSource)
			binder := NewDefaultPropertyBinder(propSources)

			var target any
			if tc.targetType != nil {
				target = reflect.New(tc.targetType).Interface()
			} else {
				target = tc.targetValue
			}

			// when
			err := binder.Bind(tc.propName, target)

			// then
			if tc.wantErr != nil {
				require.EqualError(t, err, tc.wantErr.Error())
				return
			}

			require.NoError(t, err)

			assert.Equal(t, tc.wantResult, reflect.ValueOf(target).Elem().Interface())
		})
	}
}
