package helpers

func GetValue(flg, cfg any, field string) string {
	if v := GetFieldString(flg, field); v != "" {
		return v
	}
	return GetFieldString(cfg, field)
}
