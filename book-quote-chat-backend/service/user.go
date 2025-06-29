package service

import (
	"book-quote-chat-backend/model"
	"book-quote-chat-backend/store"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"strings"
	"time"
)

var jwtSecret = []byte("your-secret-key") // 强烈建议用环境变量

// 随机昵称和头像
var nicknames = []string{"热心码友", "大前端", "Go小将", "吃瓜群众", "灵魂工程师", "AI玩家"}
var avatars = []string{
	"https://api.multiavatar.com/aric.svg",
	"https://api.multiavatar.com/go.svg",
	"https://api.multiavatar.com/vue.svg",
	"https://api.multiavatar.com/php.svg",
	"https://api.multiavatar.com/js.svg",
}

// 注册
func RegisterUser(name, password, email, avatar, group string) (model.User, error) {
	users, err := store.LoadUsers()
	if err != nil {
		return model.User{}, err
	}
	for _, u := range users {
		if u.Name == name {
			return model.User{}, errors.New("用户名已存在")
		}
		if u.Email == email {
			return model.User{}, errors.New("邮箱已存在")
		}
	}
	if name == "" {
		name = randomNickname()
	}
	if avatar == "" {
		avatar = randomAvatar()
	}
	if group == "" {
		group = "user"
	}
	// 加密密码
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, errors.New("密码加密失败")
	}
	user := model.User{
		ID:       uuid.New().String(),
		Name:     name,
		Password: string(hash),
		Avatar:   avatar,
		Email:    email,
		Group:    group,
		Status:   "normal",
		Created:  time.Now().Unix(),
	}
	// 打印 user 信息
	users = append(users, user)
	fmt.Printf("注册用户：%v\n", users)
	if err := store.SaveUsers(users); err != nil {
		return model.User{}, err
	}
	user.Password = "" // 返回时隐藏
	return user, nil
}

// 登录，返回token和用户
func LoginUser(name, password string) (string, model.User, error) {
	users, err := store.LoadUsers()
	if err != nil {
		return "", model.User{}, err
	}
	for _, u := range users {
		if u.Name == name {
			if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil {
				// 生成JWT token
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"userId": u.ID,
					"group":  u.Group,
					"exp":    time.Now().Add(24 * time.Hour).Unix(),
				})
				tokenStr, err := token.SignedString(jwtSecret)
				if err != nil {
					return "", model.User{}, err
				}
				u.Password = ""
				return tokenStr, u, nil
			}
			break
		}
	}
	return "", model.User{}, errors.New("用户名或密码错误")
}

// JWT验证
func ParseJWT(tokenStr string) (userID, group string, err error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return "", "", errors.New("无效token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errors.New("无效claims")
	}
	uid, _ := claims["userId"].(string)
	grp, _ := claims["group"].(string)
	return uid, grp, nil
}

// 查用户
func GetUserByID(id string) (model.User, error) {
	users, err := store.LoadUsers()
	if err != nil {
		return model.User{}, err
	}
	for _, u := range users {
		if u.ID == id {
			u.Password = ""
			return u, nil
		}
	}
	return model.User{}, errors.New("用户不存在")
}

func randomNickname() string {
	return nicknames[rand.Intn(len(nicknames))] + fmt.Sprintf("%d", rand.Intn(10000))
}

func randomAvatar() string {
	return avatars[rand.Intn(len(avatars))]
}

// ListUsers 支持分页、id、nickname筛选
// ListUsers 支持分页、id、nickname筛选
func ListUsers(offset, limit int, id, name string) ([]model.User, int, error) {
	users, err := store.LoadUsers()
	if err != nil {
		return nil, 0, err
	}
	var filtered []model.User
	for _, u := range users {
		if id != "" && u.ID != id {
			continue
		}
		if name != "" && !strings.Contains(u.Name, name) {
			continue
		}
		filtered = append(filtered, u)
	}
	total := len(filtered)
	start := offset
	end := offset + limit
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}
	return filtered[start:end], total, nil
}

// BanUser 封禁用户（将status设为banned）
func BanUser(id string) error {
	users, err := store.LoadUsers()
	if err != nil {
		return err
	}
	changed := false
	for i := range users {
		if users[i].ID == id {
			if users[i].Status == "banned" {
				return errors.New("用户已被封禁")
			}
			users[i].Status = "banned"
			changed = true
			break
		}
	}
	if !changed {
		return errors.New("用户不存在")
	}
	return store.SaveUsers(users)
}

// UnbanUser 解封用户
func UnbanUser(id string) error {
	users, err := store.LoadUsers()
	if err != nil {
		return err
	}
	changed := false
	for i := range users {
		if users[i].ID == id {
			if users[i].Status != "banned" {
				return errors.New("用户未被封禁")
			}
			users[i].Status = "normal"
			changed = true
			break
		}
	}
	if !changed {
		return errors.New("用户不存在")
	}
	return store.SaveUsers(users)
}
