package kodbox

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Barton0403/btgo-pkg/util"
	"io"
	"net/http"
	"net/url"
)

type KodBoxApi struct {
	host       string
	client     *http.Client
	loginToken string
	token      string
}

func NewKodBoxApi(host string) *KodBoxApi {
	client := &http.Client{}

	return &KodBoxApi{
		host:   host,
		client: client,
	}
}

func (k *KodBoxApi) Login(loginToken string) error {
	resp, e := k.client.Get(k.host + "/index.php?user/index/loginSubmit&loginToken=" + loginToken)
	if e != nil {
		return e
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New("request fail")
	}

	body, e := io.ReadAll(resp.Body)
	var obj map[string]interface{}
	e = json.Unmarshal(body, &obj)
	if e != nil {
		return e
	}

	if !obj["code"].(bool) {
		return errors.New(obj["data"].(string))
	}

	k.loginToken = loginToken
	k.token = obj["info"].(string)

	return nil
}

func (k *KodBoxApi) reLogin() error {
	return k.Login(k.loginToken)
}

type callback func() error

func (k *KodBoxApi) call(f callback) error {
	e := f()
	if e == nil {
		return e
	}

	if e.Error() == "token expire" {
		// 尝试一次重新登陆
		e = k.reLogin()
		if e != nil {
			return e
		}

		return f()
	} else {
		return e
	}
}

func (k *KodBoxApi) add(path string, source string) (fileId int64, e error) {
	resp, e := k.client.PostForm(k.host+"/index.php?plugin/BTUpload/add&accessToken="+k.token, url.Values{"path": {path}, "source": {source}})
	if e != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		e = errors.New("request fail")
		return
	}

	body, e := io.ReadAll(resp.Body)
	var obj map[string]interface{}
	e = json.Unmarshal(body, &obj)
	if e != nil {
		return
	}

	code, e := util.ToString(obj["code"])
	if e != nil {
		return
	}

	if code == "10001" {
		e = errors.New("token expire")
		return
	} else if code != "200" {
		e = errors.New(obj["msg"].(string))
		return
	}

	data := obj["data"].(map[string]interface{})
	return util.ToInt64(data["file_id"])
}

func (k *KodBoxApi) Add(path string, source string) (fileId int64, e error) {
	f := func() error {
		fileId, e = k.add(path, source)
		return e
	}
	e = k.call(f)
	return
}

func (k *KodBoxApi) update(fileId int64, path string) (e error) {
	resp, e := k.client.PostForm(k.host+"/index.php?plugin/BTUpload/updatefile&accessToken="+k.token, url.Values{"path": {path}, "file_id": {fmt.Sprintf("%v", fileId)}})
	if e != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		e = errors.New("request fail")
		return
	}

	body, e := io.ReadAll(resp.Body)
	var obj map[string]interface{}
	e = json.Unmarshal(body, &obj)
	if e != nil {
		return
	}

	code, e := util.ToString(obj["code"])
	if e != nil {
		return
	}

	if code == "10001" {
		e = errors.New("token expire")
		return
	} else if code != "200" {
		e = errors.New(obj["msg"].(string))
		return
	}

	return
}

func (k *KodBoxApi) Update(fileId int64, path string) error {
	f := func() error {
		return k.update(fileId, path)
	}
	return k.call(f)
}

type Folder struct {
	Name     string
	Path     string
	SourceID int
}

type File struct {
	Type     string
	Ext      string
	SourceID int
	ParentID int
	Name     string
	FilePath string
}

type ListData struct {
	FileList   []*File
	FolderList []*Folder
}

func (k *KodBoxApi) list(path string, page int8, pageNum int16) (data *ListData, e error) {
	resp, e := k.client.PostForm(k.host+"/index.php?explorer/list/path&accessToken="+k.token,
		url.Values{"path": {path}, "page": {fmt.Sprintf("%v", page)}, "pageNum": {fmt.Sprintf("%v", pageNum)}, "withFilePath": {"1"}})
	if e != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		e = errors.New("request fail")
		return
	}

	body, e := io.ReadAll(resp.Body)
	var obj map[string]interface{}
	e = json.Unmarshal(body, &obj)
	if e != nil {
		return
	}

	code, e := util.ToString(obj["code"])
	if e != nil {
		return
	}

	if code == "10001" {
		e = errors.New("token expire")
		return
	} else if code != "true" {
		e = errors.New(obj["data"].(string))
		return
	}

	// 通过json动态转struct
	tmpJson, e := json.Marshal(obj["data"])
	if e != nil {
		return
	}
	e = json.Unmarshal(tmpJson, &data)
	if e != nil {
		return
	}

	return
}

func (k *KodBoxApi) List(path string, page int8, pageNum int16) (data *ListData, e error) {
	f := func() error {
		data, e = k.list(path, page, pageNum)
		return e
	}
	e = k.call(f)
	return
}

func (k *KodBoxApi) addServerFile(path string, sourceId int) (e error) {
	resp, e := k.client.PostForm(k.host+"/index.php?plugin/BTUpload/add&accessToken="+k.token,
		url.Values{"path": {path}, "source": {fmt.Sprintf("%v", sourceId)}})
	if e != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		e = errors.New("request fail")
		return
	}

	body, e := io.ReadAll(resp.Body)
	var obj map[string]interface{}
	e = json.Unmarshal(body, &obj)
	if e != nil {
		return
	}

	code, e := util.ToString(obj["code"])
	if e != nil {
		return
	}

	if code == "10001" {
		e = errors.New("token expire")
		return
	} else if code != "200" {
		e = errors.New(obj["msg"].(string))
		return
	}

	return
}

func (k *KodBoxApi) AddServerFile(path string, sourceId int) error {
	f := func() error {
		return k.addServerFile(path, sourceId)
	}
	return k.call(f)
}
