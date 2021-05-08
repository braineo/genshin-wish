package parser

type GachaItem struct {
	ItemID   string `json:"item_id"`
	Name     string `json:"name"`
	ItemType string `json:"item_type"`
	RankType string `json:"rank_type"`
}

type GachaConfig struct {
	ID   string `gorm:"not null" json:"id"`
	Key  string `gorm:"not null" json:"key"`
	Name string `gorm:"not null" json:"name"`
}

type GachaConfigResponse struct {
	RetCode int    `json:"retcode"`
	Message string `json:"message"`
	Data    struct {
		GachaTypeList []GachaConfig `json:"gacha_type_list"`
	} `json:"data"`
}

type GachaLog struct {
	GachaType string `gorm:"not null" json:"gacha_type"`
	ID        string `gorm:"not null" json:"id"` // id for pagination
	UID       string `gorm:"not null" json:"uid"`
	Count     string `gorm:"not null" json:"conut"`
	ItemID    string `gorm:"not null" json:"item_id"`
	Time      string `gorm:"not null" json:"time"`
	ItemType  string `gorm:"not null" json:"item_type"`
	RankType  string `gorm:"not null" json:"rank_type"`
	Name      string `gorm:"not null" json:"name"`
	Lang      string `gorm:"not null" json:"lang"`
}

type GachaLogResponse struct {
	RetCode int    `json:"retcode"`
	Message string `json:"message"`
	Data    struct {
		Size         string     `json:"size"`
		Region       string     `json:"region"`
		GachaLogList []GachaLog `json:"list"`
	} `json:"data"`
}