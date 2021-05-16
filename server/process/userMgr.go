package process

//UserMgr 实例在服务器端有且只有一个
//TODO 需要使用rul链  定期清理长时间不在线的用户
var(
	userMgr *UserMgr
)
type UserMgr struct {
	onlineUsers map[int]*UserProcess
}

func init(){
	userMgr = &UserMgr{
		onlineUsers: make(map[int]*UserProcess, 1024),
	}
}

func (this *UserMgr) AddOnlineUser(up *UserProcess){
	this.onlineUsers[up.UserId] = up
}

func (this *UserMgr) DelOnlineUser(userId int){
	delete (this.onlineUsers, userId)
}

func (this *UserMgr) GetAllOnlineUser()map[int]*UserProcess{
	return this.onlineUsers
}

func (this *UserMgr) GetOnlineUserByUserId(userId int)(up *UserProcess){
	return this.onlineUsers[userId]
}