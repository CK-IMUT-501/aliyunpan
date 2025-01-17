package syncdrive

import (
	"encoding/json"
	"fmt"
	"github.com/tickstep/aliyunpan-api/aliyunpan"
	"github.com/tickstep/aliyunpan/internal/utils"
	"github.com/tickstep/library-go/logger"
	"io/ioutil"
	"path"
)

type (
	// SyncTaskManager 同步任务管理器
	SyncTaskManager struct {
		syncDriveConfig      *SyncDriveConfig
		DriveId              string
		PanClient            *aliyunpan.PanClient
		SyncConfigFolderPath string
	}

	// SyncDriveConfig 同步盘配置文件
	SyncDriveConfig struct {
		ConfigVer    string      `json:"configVer"`
		SyncTaskList []*SyncTask `json:"syncTaskList"`
	}
)

var (
	ErrSyncTaskListEmpty error = fmt.Errorf("no sync task")
)

func NewSyncTaskManager(driveId string, panClient *aliyunpan.PanClient, syncConfigFolderPath string) *SyncTaskManager {
	return &SyncTaskManager{
		DriveId:              driveId,
		PanClient:            panClient,
		SyncConfigFolderPath: syncConfigFolderPath,
	}
}

func (m *SyncTaskManager) parseConfigFile() error {
	/** 样例
	{
	 "configVer": "1.0",
	 "syncTaskList": [
	  {
	   "name": "NS游戏备份",
	   "id": "5b2d7c10-e927-4e72-8f9d-5abb3bb04814",
	   "driveId": "19519111",
	   "localFolderPath": "D:\\smb\\datadisk\\game",
	   "panFolderPath": "/sync_drive/game",
	   "mode": "sync",
	   "lastSyncTime": ""
	  }
	 ]
	}
	*/
	configFilePath := m.configFilePath()
	r := &SyncDriveConfig{
		ConfigVer:    "1.0",
		SyncTaskList: []*SyncTask{},
	}
	m.syncDriveConfig = r

	if b, _ := utils.PathExists(configFilePath); b != true {
		text := utils.ObjectToJsonStr(r, true)
		ioutil.WriteFile(configFilePath, []byte(text), 0600)
		return nil
	}
	data, e := ioutil.ReadFile(configFilePath)
	if e != nil {
		return e
	}

	if len(data) > 0 {
		if err2 := json.Unmarshal(data, m.syncDriveConfig); err2 != nil {
			logger.Verboseln("parse sync drive config json error ", err2)
			return err2
		}
	}
	return nil
}

func (m *SyncTaskManager) configFilePath() string {
	return path.Join(m.SyncConfigFolderPath, "sync_drive_config.json")
}

// Start 启动同步进程
func (m *SyncTaskManager) Start() (bool, error) {
	if m.parseConfigFile() != nil {
		return false, nil
	}
	if m.syncDriveConfig.SyncTaskList == nil || len(m.syncDriveConfig.SyncTaskList) == 0 {
		return false, ErrSyncTaskListEmpty
	}

	// start the sync task one by one
	for _, task := range m.syncDriveConfig.SyncTaskList {
		if len(task.Id) == 0 {
			task.Id = utils.UuidStr()
		}
		task.DriveId = m.DriveId
		task.syncDbFolderPath = m.SyncConfigFolderPath
		task.panClient = m.PanClient
		if e := task.Start(); e != nil {
			logger.Verboseln(e)
			fmt.Println("start sync task error: {}", task.Id)
			continue
		}
		fmt.Println("\n启动同步任务")
		fmt.Println(task)
	}
	// save config file
	ioutil.WriteFile(m.configFilePath(), []byte(utils.ObjectToJsonStr(m.syncDriveConfig, true)), 0600)
	return true, nil
}

// Stop 停止同步进程
func (m *SyncTaskManager) Stop() (bool, error) {
	// stop task one by one
	for _, task := range m.syncDriveConfig.SyncTaskList {
		if e := task.Stop(); e != nil {
			logger.Verboseln(e)
			fmt.Println("stop sync task error: ", task.NameLabel())
			continue
		}
		fmt.Println("停止同步任务: ", task.NameLabel())
	}

	// save config file
	ioutil.WriteFile(m.configFilePath(), []byte(utils.ObjectToJsonStr(m.syncDriveConfig, true)), 0600)
	return true, nil
}
