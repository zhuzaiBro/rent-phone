package reservation

type Date struct {
	ClerkID string `gorm:"column:clerk_id;type:varchar(255)"`
	Time    string `gorm:"column:arrange_time;type"` // 11:00-12:00

}
