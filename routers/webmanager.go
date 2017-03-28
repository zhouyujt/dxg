package routers

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/astaxie/beego/session"
	"github.com/zhouyujt/dxg/config"
)

type WebManager struct {
	sessionManager *session.Manager
	session        session.Store
	cfg            *config.Config
}

func NewWebManager(cfg *config.Config) *WebManager {
	mgr := new(WebManager)
	mgr.cfg = cfg

	config := `{"cookieName":"gosessionid","gclifetime":10, "enableSetCookie":true}`
	conf := new(session.ManagerConfig)
	json.Unmarshal([]byte(config), conf)
	mgr.sessionManager, _ = session.NewManager("memory", conf)
	go mgr.sessionManager.GC()

	return mgr
}

func (mgr *WebManager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mgr.session, _ = mgr.sessionManager.SessionStart(w, r)

	if r.Method == "GET" {
		mgr.get(w, r)
	} else if r.Method == "POST" {
		mgr.post(w, r)
	}
}

func (mgr *WebManager) get(w http.ResponseWriter, r *http.Request) {
	user := mgr.session.Get("user")
	if user == nil {
		mgr.showLoginPage(w)
	} else {
		mgr.showAdminPage(w)
	}
}

func (mgr *WebManager) post(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	uname := r.Form.Get("uname")
	pwd := r.Form.Get("pwd")

	if uname != mgr.cfg.WebManager.User {
		mgr.showLoginPage(w)
		return
	}

	h := md5.New()
	h.Write([]byte(mgr.cfg.WebManager.Pwd))
	pwd2 := hex.EncodeToString(h.Sum(nil))
	if pwd != string(pwd2) {
		mgr.showLoginPage(w)
		return
	}

	mgr.session.Set("user", uname)
	http.Redirect(w, r, "/admin", 302)
}

func (mgr *WebManager) showLoginPage(w http.ResponseWriter) {
	f, err := os.Open("webmanager/views/login.html")
	if err != nil {
		return
	}
	b, err := ioutil.ReadAll(f)
	w.Write(b)
	f.Close()
}

func (mgr *WebManager) showAdminPage(w http.ResponseWriter) {
	f, err := os.Open("webmanager/views/admin.html")
	if err != nil {
		return
	}
	b, err := ioutil.ReadAll(f)
	w.Write(b)
	f.Close()
}
