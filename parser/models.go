package parser

// Response of gacha config from mihoyo API
type GachaItem struct {
	ID       string `json:"item_id"`
	Name     string `json:"name"`
	ItemType string `gorm:"uniqueIndex" json:"item_type"`
	RankType string `json:"rank_type"`
}

// Response of gacha config from mihoyo API
type GachaConfig struct {
	ID   string `gorm:"not null" json:"id"`
	Key  string `gorm:"not null" json:"key"`
	Name string `gorm:"not null" json:"name"`
}

// Response of gacha config from mihoyo API
type GachaConfigResponse struct {
	RetCode int    `json:"retcode"`
	Message string `json:"message"`
	Data    struct {
		GachaTypeList []GachaConfig `json:"gacha_type_list"`
	} `json:"data"`
}

// Response of gacha log from mihoyo API
type GachaLog struct {
	GachaType string `gorm:"not null" json:"gacha_type"`
	ID        string `gorm:"primary_key" json:"id"` // id for pagination
	UID       string `gorm:"not null" json:"uid"`
	Count     string `gorm:"not null" json:"conut"`
	ItemID    string `gorm:"not null" json:"item_id"`
	// Pull time in format "2021-05-06 18:44:47", always in UTC-8
	Time     string `gorm:"not null;type:datetime" json:"time"`
	ItemType string `gorm:"not null" json:"item_type"`
	RankType string `gorm:"not null" json:"rank_type"`
	Name     string `gorm:"not null" json:"name"`
	Lang     string `gorm:"not null" json:"lang"`
}

// Response of gacha log from mihoyo API
type GachaLogResponse struct {
	RetCode int    `json:"retcode"`
	Message string `json:"message"`
	Data    struct {
		Size         string     `json:"size"`
		Region       string     `json:"region"`
		GachaLogList []GachaLog `json:"list"`
	} `json:"data"`
}
