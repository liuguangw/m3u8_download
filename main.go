package main

import (
	"fmt"
	"github.com/gosuri/uiprogress"
	"github.com/liuguangw/m3u8_download/download"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	configPath := "F:\\movie\\config.json"
	if len(os.Args) < 2 {
		fmt.Println("Use: " + filepath.Base(os.Args[0]) + " <configPath>")
		//return
	} else {
		configPath = os.Args[1]
	}
	taskConfig, err := download.ReadTaskConfig(configPath)
	if err != nil {
		log.Fatalln("read config error: ", err)
	}
	taskStatus, err := download.GetTaskStatus(taskConfig)
	if err != nil {
		log.Fatalln(err)
	}
	saveDir := filepath.Join(taskConfig.SaveDir, "tmp", taskConfig.FileName)
	_, err = os.Stat(saveDir)
	if os.IsNotExist(err) {
		err = os.Mkdir(saveDir, 0644)
		if err != nil {
			log.Fatalln("make dir (" + saveDir + ") error: " + err.Error())
		}
	}
	err = download.WriteTsList(filepath.Join(saveDir, "000list.txt"), len(taskStatus.Nodes))
	if err != nil {
		log.Fatalln("write list.txt error: " + err.Error())
	}
	downloadSuccessCount := 0
	dataFile := filepath.Join(saveDir, "000task.data")
	_, err = os.Stat(dataFile)
	if !os.IsNotExist(err) {
		taskDataBits, err := ioutil.ReadFile(dataFile)
		if err != nil {
			log.Fatalln("read file (" + dataFile + ") error: " + err.Error())
		}
		for nodeIndex, nodeStatus := range taskDataBits {
			if nodeStatus == download.STATUS_SUCCESS {
				taskStatus.Nodes[nodeIndex].Status = nodeStatus
				downloadSuccessCount++
			}
		}
	}
	uiprogress.Start()
	mainDownTask := &download.MainDownloadTask{
		TaskStatus:       taskStatus,
		Config:           taskConfig,
		NewTaskIndex:     make(chan int),
		DownloadComplete: make(chan bool),
		Bar:              uiprogress.AddBar(len(taskStatus.Nodes) - downloadSuccessCount),
	}
	//_ = mainDownTask.Bar.Set(downloadSuccessCount)
	mainDownTask.Bar.AppendCompleted()
	mainDownTask.Bar.PrependElapsed()
	mainDownTask.Bar.PrependFunc(func(b *uiprogress.Bar) string {
		return "download: " + strconv.Itoa(b.Current()) + "/" + strconv.Itoa(b.Total)
	})
	go mainDownTask.RunBackend()
	for i := 0; i < taskConfig.MaxTask; i++ {
		go mainDownTask.RunDownload()
	}
	<-mainDownTask.DownloadComplete
}
