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

func (j *JsonRepo) CreateUser(ctx context.Context, user *models.User) (string, error) {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	f, err := ioutil.ReadFile(j.filename)
	if err != nil {
		return "", err
	}

	us := UserStore{}

	if string(f) == "" {
		us.List = make(UserList)
	} else {
		if err = json.Unmarshal(f, &us); err != nil {
			return "", err
		}
	}

	us.Increment++
	id := strconv.Itoa(us.Increment)
	us.List[id] = user

	if f, err = json.Marshal(us); err != nil {
		return "", err
	}

	if err = ioutil.WriteFile(j.filename, f, fs.ModePerm); err != nil {
		return "", err
	}

	return id, nil
}

func (j *JsonRepo) DeleteUser(ctx context.Context, userId string) error {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	f, err := ioutil.ReadFile(j.filename)
	if err != nil {
		return err
	}

	us := UserStore{}

	if string(f) == "" {
		return user.ErrUserNotFound
	} else {
		if err = json.Unmarshal(f, &us); err != nil {
			return err
		}
	}

	if _, ok := us.List[userId]; !ok {
		return user.ErrUserNotFound
	}

	delete(us.List, userId)

	if f, err = json.Marshal(us); err != nil {
		return err
	}

	if err = ioutil.WriteFile(j.filename, f, fs.ModePerm); err != nil {
		return err
	}

	return nil
}

func (j *JsonRepo) GetUser(ctx context.Context, userId string) (*models.User, error) {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	f, err := ioutil.ReadFile(j.filename)
	if err != nil {
		return nil, err
	}

	us := UserStore{}

	if string(f) == "" {
		us.List = make(UserList)
	} else {
		if err = json.Unmarshal(f, &us); err != nil {
			return nil, err
		}
	}

	if _, ok := us.List[userId]; !ok {
		return nil, user.ErrUserNotFound
	}

	return us.List[userId], nil
}

func (j *JsonRepo) SearchUser(ctx context.Context) (map[string]*models.User, error) {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	f, err := ioutil.ReadFile(j.filename)
	if err != nil {
		return nil, err
	}

	us := UserStore{}

	if string(f) == "" {
		us.List = make(UserList)
	} else {
		if err = json.Unmarshal(f, &us); err != nil {
			return nil, err
		}
	}

	return us.List, nil
}

func (j *JsonRepo) UpdateUser(ctx context.Context, userId, newName string) error {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	f, err := ioutil.ReadFile(j.filename)
	if err != nil {
		return err
	}

	if string(f) == "" {
		return user.ErrUserNotFound
	}

	us := UserStore{}
	err = json.Unmarshal(f, &us)
	if err != nil {
		return err
	}

	if _, ok := us.List[userId]; !ok {
		return user.ErrUserNotFound
	}

	us.List[userId].DisplayName = newName

	if f, err = json.Marshal(us); err != nil {
		return err
	}

	if err = ioutil.WriteFile(j.filename, f, fs.ModePerm); err != nil {
		return err
	}

	return nil
}
