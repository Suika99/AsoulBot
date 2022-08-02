package minecraft

import (
	"github.com/FloatTech/zbputils/control"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"os"
)

type Mcserver struct {
	Name string `gorm:"primary_key"`
	IP   string
	Port string
}

type ServerInfo struct {
	Status  string `json:"status"`
	Online  bool   `json:"online"`
	Motd    string `json:"motd"`
	Players struct {
		Max    int `json:"max"`
		Now    int `json:"now"`
		Sample []struct {
			Name string `json:"name"`
			ID   string `json:"id"`
		} `json:"sample"`
	} `json:"players"`
	Server struct {
		Name     string `json:"name"`
		Protocol int    `json:"protocol"`
	} `json:"server"`
	LastUpdated string `json:"last_updated"`
	Duration    string `json:"duration"`
}

var (
	datapath = "data/minecraft/"
	dbfile   = datapath + "server.db"
	engine   = control.Register("minecraft", &control.Options{
		DisableOnDefault: false,
		Help:             "查mc服务器在线玩家信息",
	})
)

// 初始化数据库
func initDB() {
	var err error
	if _, err = os.Stat(dbfile); err != nil || os.IsNotExist(err) {
		_ = os.MkdirAll(datapath, 0755)
		f, err := os.Create(dbfile)
		if err != nil {
			log.Error("[minecraft]创建数据库文件失败:", err)
			return
		}
		log.Infof("[Minecraft]数据库文件(%v)创建成功", dbfile)
		db, err := gorm.Open("sqlite3", dbfile)
		if err != nil {
			log.Errorln("[Minecraft]打开数据库失败：", err)
			return
		}
		db.AutoMigrate(Mcserver{})
		defer f.Close()
		defer db.Close()
	}
}

// 插入数据
func insertData(server *Mcserver) error {
	db, err := gorm.Open("sqlite3", dbfile)
	if err != nil {
		log.Errorln("[Minecraft]打开数据库失败:", err)
		return err
	}
	err = db.Create(&Mcserver{
		Name: server.Name,
		IP:   server.IP,
		Port: server.Port,
	}).Error
	if err != nil {
		log.Errorln("[Minecraft]插入数据失败:", err)
		return err
	}
	defer db.Close()
	return nil
}

// 删除数据
func deleteData(name string) error {
	db, err := gorm.Open("sqlite3", dbfile)
	if err != nil {
		log.Errorln("[Minecraft]打开数据库失败:", err)
		return err
	}
	err = db.Delete(&Mcserver{}, name).Error
	if err != nil {
		log.Errorln("[Minecraft]插入数据失败:", err)
		return err
	}
	return nil
}

// 查找数据
func selectData(name string) (server Mcserver, err error) {
	db, err := gorm.Open("sqlite3", dbfile)
	if err != nil {
		log.Errorln("[Minecraft]打开数据库失败:", err)
		return
	}
	if err = db.Debug().Model(&Mcserver{}).Find(&server, "name = (?)", name).Error; err != nil {
		return server, err
	}
	return
}
