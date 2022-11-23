package repository

import (
	"context"
	"encoding/json"
	"github.com/sornick01/UserAPI/internal/user"
	"github.com/sornick01/UserAPI/models"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"strconv"
	"sync"
	"time"
)

type (
	UserList map[string]*models.User

	UserStore struct {
		Increment int      `json:"increment"`
		List      UserList `json:"list"`
	}
)
type JsonRepo struct {
	//Increment int `json:"increment"`
	mutex *sync.Mutex
	//Users     map[string]*models.User `json:"users"`
	filename string
}

func NewUserStore() *UserStore {

}

func NewJsonRepo(filename string) *JsonRepo {
	return &JsonRepo{
		mutex: new(sync.Mutex),
		//Users:    make(map[string]*models.User),
		filename: filename,
	}
}

func (j *JsonRepo) CreateUser(ctx context.Context, user *models.User) (string, error) {
	j.mutex.Lock()
	defer j.mutex.Unlock() //TODO: need?
	file, err := os.Open(j.filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	//err := readFromFile(j.filename, j) //TODO: empty file err
	//if err != nil {
	//	return "", err
	//}
	//
	//j.Increment++
	//id := strconv.Itoa(j.Increment)
	//user.CreatedAt = time.Now()
	//j.Users[id] = user
	//
	//err = writeToFile(j.filename, j)
	//if err != nil {
	//	return "", err
	//}
	//
	//return id, nil

	userStore, err := userStoreFromFile(j.filename)
	if err != nil {
		return "", err
	}

	userStore.Increment++
	id := strconv.Itoa(userStore.Increment)
	user.CreatedAt = time.Now()
	userStore.List[id] = user

}

func (j *JsonRepo) DeleteUser(ctx context.Context, userId string) error {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	err := readFromFile(j.filename, j)
	if err != nil {
		return err
	}

	if _, ok := j.Users[userId]; !ok {
		return user.ErrUserNotFound
	}

	delete(j.Users, userId)

	err = writeToFile(j.filename, j)
	if err != nil {
		return err
	}

	return nil
}

func (j *JsonRepo) GetUser(ctx context.Context, userId string) (*models.User, error) {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	err := readFromFile(j.filename, j)
	if err != nil {
		return nil, err
	}

	if _, ok := j.Users[userId]; !ok {
		return nil, user.ErrUserNotFound
	}

	return j.Users[userId], nil
}

func (j *JsonRepo) SearchUser(ctx context.Context) (map[string]*models.User, error) {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	err := readFromFile(j.filename, j)
	if err != nil {
		return nil, err
	}

	return j.Users, nil
}

func (j *JsonRepo) UpdateUser(ctx context.Context, userId, newDisplayName string) error {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	err := readFromFile(j.filename, j)
	if err != nil {
		return err
	}

	if _, ok := j.Users[userId]; !ok {
		return user.ErrUserNotFound
	}

	u := j.Users[userId]
	u.DisplayName = newDisplayName

	err = writeToFile(j.filename, j)
	if err != nil {
		return err
	}
	return nil
}

func readFromFile(filename string, repo *JsonRepo) error {
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(f, repo)
	if err != nil {
		return err
	}

	return nil
}

func writeToFile(filename string, repo *JsonRepo) error {
	f, err := json.Marshal(repo)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, f, fs.ModePerm)
}

func userStoreFromFile(file *os.File) (*UserStore, error) {
	f, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	userStore := new(UserStore)
	if len(f) == 0 {
		userStore.List = make(UserList)
	} else {
		err = json.Unmarshal(f, userStore)
		if err != nil {
			return nil, err
		}
	}

	return userStore, nil
}

func userStoreToFile(file *os.File, userStore *UserStore) error {
	f, err := json.Marshal(us)
}