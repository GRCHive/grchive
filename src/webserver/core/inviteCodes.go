package core

type InviteCode struct {
	Id         int64    `db:"id"`
	FromUserId int64    `db:"from_user_id"`
	FromOrgId  int32    `db:"from_org_id"`
	ToEmail    string   `db:"to_email"`
	SentTime   NullTime `db:"sent_time"`
	UsedTime   NullTime `db:"used_time"`
}
