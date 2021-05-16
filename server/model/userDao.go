package model

import (
	"encoding/json"
	"fmt"
	error2 "gin_example/chartroom/common/error"
	"gin_example/chartroom/common/message"
	"github.com/garyburd/redigo/redis"
)

var (
	MyUserDao *UserDao
)
//完成对user的操作

type UserDao struct {
	pool *redis.Pool
}

func NewUserDao(pool *redis.Pool)*UserDao{
	userDao := &UserDao{
		pool: pool,
	}
	return userDao
}

func (this *UserDao) getUserById(conn redis.Conn, id int) (user *message.User, errno error2.Errno){
	res, err := redis.String(conn.Do("HGet", "users", id))
	if err != nil {
		if err == redis.ErrNil{
			errno = error2.ERROR_USER_NOTEXISTS
		}
		return
	}
	user = &message.User{}
	err = json.Unmarshal([]byte(res), user)
	if err != nil {
		fmt.Println("[getUserById] json.Unmarshal user err=", err)
		errno = error2.ERROR_SEVER_ERR
		return
	}
	errno = error2.ERROR_OK
	return
}

func (this *UserDao) Login(userId int, userPwd string) (user *message.User, errno error2.Errno){
	conn := this.pool.Get()
	defer conn.Close()

	user, errno = this.getUserById(conn, userId)
	if errno != error2.ERROR_OK {
		return
	}

	if user.UserPwd != userPwd{
		errno = error2.ERROR_USER_PWD
		return
	}

	errno = error2.ERROR_OK
	return
}

func (this *UserDao) Register(user *message.User) (errno error2.Errno){
	conn := this.pool.Get()
	defer conn.Close()

	_, errno = this.getUserById(conn, user.UserId)
	if errno == error2.ERROR_OK {
		errno = error2.ERROR_USER_EXISTS
		return
	}

	data, err := json.Marshal(user)
	if err != nil {
		return
	}
	_, err = conn.Do("HSet", "users", user.UserId, string(data))
	if err != nil {
		fmt.Println("[Register] HSet users err=", err)
		errno = error2.ERROR_SEVER_ERR
		return
	}
	errno = error2.ERROR_OK
	return
}