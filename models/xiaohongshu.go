package models

type XhsBaseModel struct {
	ID           uint   `gorm:"primarykey"`
	UserID       string `gorm:"size:64;comment:用户ID"`
	NickName     string `gorm:"size:255;comment:用户昵称"`
	Avatar       string `gorm:"size:255;comment:用户头像地址"`
	IpLocation   string `gorm:"size:255;comment:用户的IP地理位置"`
	AddTs        int64  `gorm:"comment:记录添加时间戳"`
	LastModifyTs int64  `gorm:"comment:记录最后修改时间戳"`
}

type XhsNote struct {
	XhsBaseModel
	NoteId         string `gorm:"size:64;index;comment:笔记ID"`
	Type           string `gorm:"size:16;comment:笔记类型(normal | video)"`
	Title          string `gorm:"size:255;comment:笔记标题"`
	Desc           string `gorm:"size:255;comment:笔记描述"`
	Time           int64  `gorm:"size:255;comment:笔记发布时间戳"`
	LastUpdateTime int64  `gorm:"size:255;comment:笔记最后更新时间戳"`
	LikeCount      string `gorm:"size:16;comment:笔记点赞数"`
	CollectedCount string `gorm:"size:16;comment:笔记收藏数"`
	CommentCount   string `gorm:"size:16;comment:笔记评论数"`
	ShareCount     string `gorm:"size:16;comment:笔记分享数"`
	ImageList      string `gorm:"comment:笔记封面图片列表"`
}

func (XhsNote) TableName() string {
	return "xhs_note"
}

type XHSNoteComment struct {
	XhsBaseModel
	CommentID       string `gorm:"index;size:64;comment:评论ID"`
	CreateTime      int64  `gorm:"index;comment:评论时间戳"`
	NoteID          string `gorm:"size:64;comment:笔记ID"`
	Content         string `gorm:"comment:评论内容"`
	SubCommentCount int    `gorm:"comment:子评论数量"`
}

func (XHSNoteComment) TableName() string {
	return "xhs_note_comment"
}
