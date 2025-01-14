package common

type UserRel struct {
	ChatId int64
	UserId int64
}

func UserRelHasher(rel UserRel) uint64 {
	a := rel.ChatId
	b := rel.UserId

	const prime = 0x9e3779b97f4a7c15 // 一个常用的素数
	return uint64(a) ^ (prime + uint64(b) + (uint64(a) << 6) + (uint64(a) >> 2))
}
