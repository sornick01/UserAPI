package repository

import (
	"context"
	"encoding/json"
	"github.com/sornick01/UserAPI/internal/user"
	"github.com/sornick01/UserAPI/models"
	"io/fs"
	"io/ioutil"
	"strconv"
	"sync"
)

type (
	UserList map[string]*models.User

	UserStore struct {
		Increment int      `json:"increment"`
		List      UserList `json:"list"`
	}
)
type JsonRepo struct {
	mutex    *sync.Mutex
	filename string
}

func NewJsonRepo(filename string) *JsonRepo {
	return &JsonRepo{
		mutex:    new(sync.Mutex),
		filename: filename,
	}
}

//TODO: empty json file, error managment, code repeating
func (j *JsonRepo) CreateUser(ctx context.Context, user *models.User) (string, error) {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	f, _ := ioutil.ReadFile(j.filename)
	us := UserStore{}
	_ = json.Unmarshal(f, &us)

	us.Increment++
	id := strconv.Itoa(us.Increment)
	us.List[id] = user

	f, _ = json.Marshal(us)
	ioutil.WriteFile(j.filename, f, fs.ModePerm)

	return id, nil
}

func (j *JsonRepo) DeleteUser(ctx context.Context, userId string) error {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	f, _ := ioutil.ReadFile(j.filename)
	us := UserStore{}
	_ = json.Unmarshal(f, &us)

	if _, ok := us.List[userId]; !ok {
		return user.ErrUserNotFound
	}

	delete(us.List, userId)

	f, _ = json.Marshal(us)
	ioutil.WriteFile(j.filename, f, fs.ModePerm)

	return nil
}

func (j *JsonRepo) GetUser(ctx context.Context, userId string) (*models.User, error) {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	f, _ := ioutil.ReadFile(j.filename)
	us := UserStore{}
	_ = json.Unmarshal(f, &us)

	if _, ok := us.List[userId]; !ok {
		return nil, user.ErrUserNotFound
	}

	return us.List[userId], nil
}

func (j *JsonRepo) SearchUser(ctx context.Context) (map[string]*models.User, error) {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	f, _ := ioutil.ReadFile(j.filename)
	us := UserStore{}
	_ = json.Unmarshal(f, &us)

	return us.List, nil
}

func (j *JsonRepo) UpdateUser(ctx context.Context, userId, newName string) error {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	f, _ := ioutil.ReadFile(j.filename)
	us := UserStore{}
	_ = json.Unmarshal(f, &us)

	if _, ok := us.List[userId]; !ok {
		return user.ErrUserNotFound
	}

	u := us.List[userId]
	u.DisplayName = newName

	f, _ = json.Marshal(us)
	ioutil.WriteFile(j.filename, f, fs.ModePerm)

	return nil
}
