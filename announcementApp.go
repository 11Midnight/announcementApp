package main

import (
	"fmt"
	"github.com/11Midnight/announcementApp/dbOp"
	"github.com/11Midnight/announcementApp/kingpinOp"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"strings"
	"time"
)

//Обьявление аргументов и команд.
var (
	app            = kingpin.New("notification", "A notification of tasks application.")
	register       = app.Command("register", "Register a new task.")
	registerDate   = register.Arg("date", "Date of the task. Format: 'yyyy-MM-dd hh:mm:ss'").String()
	registerTask   = register.Arg("task", "Description of task.").String()
	announcement   = app.Command("announcement", "Notification mode")
	announcementIn = announcement.Arg("announcementIn", "In how long make announcement").String()
)
// Вывод уведомления о таске в указанное время.
func callAt(callTime string, task string, durationIn time.Duration) (time.Duration, error) {
	// Описание формата времени.
	timeLayout:= "2006-01-02T15:04:05"
	//Разбираем время таска, не придумал лучше способа как это сделать с тем форматом что приходит с базы данных.
	callTime = strings.ReplaceAll(callTime, " ", "T")
	ctime, err := time.Parse(timeLayout, callTime)
	if err != nil {
		return 0, err
	}
	//Вычисляем временной промежуток до таска.
	duration := ctime.Sub(time.Now().UTC())
	//Вычисляем временной промежуток до таска который задал юзер.
	durationUser := ctime.Sub(time.Now().UTC().Add(durationIn))
	//Создам горутин.
	go func() {
		//Выводим уведомление через столько таск относительно момента запуска программы.
		fmt.Println(task, "is in ", duration.String())
		//Если времени до таска больше чем задал юзер(то есть можно вывести уведомление, что таск за "время которое задал юзер") то слипаем горутин на это время.
		if durationUser > 0 {
			time.Sleep(durationUser)
			//Выводим уведомление что таск через "время которое задал юзер".
			fmt.Println(task, "is in "+durationIn.String())
			//Слипаем на время до самого таска.
			time.Sleep(duration - durationUser)
		} else {
			//Если времени меньше, то тоже слипаем но на время до самого таска.
			time.Sleep(duration)
		}
		//Выводим уведомление что таск прямо сейчас.
		fmt.Println(task, "is right now")
	}()
	return duration, nil
}
func main() {
	db, err := dbOp.Connect("root:@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	var durMax time.Duration = 0
	switch kingpinOp.Read(app,os.Args[1:]) {
	// Register user
	case register.FullCommand():
		err := dbOp.InsertTask(db,*registerDate, *registerTask)
		if err != nil {
			panic(err)
		}
	case announcement.FullCommand():
		tasks, err := dbOp.ReadTasks(db)
		if err != nil {
			panic(err)
		}
		for _, t := range tasks {
			hr, _ := time.ParseDuration(*announcementIn)
			dur, err := callAt(t.Date, t.Task, hr)
			if err != nil {
				panic(err)
			}
			//Вычесляем самый поздний таск(таск до которого самый большой промежуток времени).
			if dur > durMax {
				durMax = dur
			}
		}
		//Продолжаем работу программы до того момента пока не наступит самый поздний таск.
		time.Sleep(durMax)
	}
}
