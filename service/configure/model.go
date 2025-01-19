package configure

type Configure struct {
	Name string `gorm:"uniqueIndex"`
	Data string
}
