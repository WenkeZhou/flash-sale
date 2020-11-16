package setting

import "time"

type ServerSettings struct {
	RunMode      string
	HttpPort     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type AppSettingS struct {
	DefaultPageSize      int
	MaxPageSize          int
	LogSavePath          string
	LogFileName          string
	LogFileExt           string
	UploadSavePath       string
	UploadServerUrl      string
	UploadImageMaxSize   int
	UploadImageAllowExts []string
}

type DataBaseSettingS struct {
	DBType       string
	Username     string
	Password     string
	Host         string
	DBName       string
	TablePrefix  string
	CharSet      string
	ParseTime    bool
	MaxIdleConns int
	MaxOpenConns int
}

type RedisSettingS struct {
	Address     string
	Password    string
	LinkType    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

type VerifySettingS struct {
	VerifySalt           string
	UserHashKeyPrefix    string
	UserVisitCountPrefix string
	MaxUserBuyCount      int
}

type BusinessSettingS struct {
	StockCachePrefix string
}

//type JWTSettingS struct {
//	Secret string
//	Issuer string
//	Expire time.Duration
//}

//type EmailSettingS struct {
//	Host     string
//	Port     int
//	UserName string
//	Password string
//	IsSSL    bool
//	From     string
//	To       []string
//}

var sections = make(map[string]interface{})

func (s *Setting) ReadSection(k string, v interface{}) error {
	err := s.vp.UnmarshalKey(k, v)
	if err != nil {
		return err
	}

	if _, ok := sections[k]; !ok {
		sections[k] = v
	}

	return nil
}

func (s *Setting) ReloadAllSection() error {
	for k, v := range sections {
		err := s.ReadSection(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}
