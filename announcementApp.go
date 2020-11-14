package main

import (
	"fmt"
	"github.com/11Midnight/announcementApp/dbOp"
	"github.com/11Midnight/announcementApp/kingpinOp"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strings"
	"time"
)

type (
	Event struct {
		Data time.Duration
	}
	Observer interface {
		OnNotify(Event)
	}
	eventObserver struct {
		task string
		time string
	}
	eventNotifier struct {
		observers map[Observer]struct{}
	}
)

func stringToDuration(timeTask string) (dur time.Duration, err error) {
	timeLayout := "2006-01-02T15:04:05"
	callTime := strings.ReplaceAll(timeTask, " ", "T")
	ctime, err := time.Parse(timeLayout, callTime)
	if err != nil {
		return 0, err
	}
	//dur := time.Until(ctime)
	duration := ctime.Sub(time.Now().UTC())
	return duration, nil
}
func (o *eventObserver) OnNotify(e Event) {
	duration, err := stringToDuration(o.time)
	if err != nil {
		panic(err)
	}
	durationUser := duration - e.Data
	if duration < time.Second && duration > 0 {
		fmt.Println(o.task, "is right now")
	}
	if durationUser < time.Second && durationUser > 0 {
		fmt.Println(o.task, "is in "+e.Data.String())
	}
}

func (o *eventNotifier) Register(l Observer) {
	o.observers[l] = struct{}{}
}

func (o *eventNotifier) Deregister(l Observer) {
	delete(o.observers, l)
}

func (p *eventNotifier) Notify(e Event) {
	for o := range p.observers {
		o.OnNotify(e)
	}
}

//Обьявление аргументов и команд.
var (
	app            = kingpinOp.NewApp("notification", "A notification of tasks application.")
	register       = app.Command("register", "Register a new task.")
	registerDate   = register.Arg("date", "Date of the task. Format: 'yyyy-MM-dd hh:mm:ss'").String()
	registerTask   = register.Arg("task", "Description of task.").String()
	announcement   = app.Command("announcement", "Notification mode")
	announcementIn = announcement.Arg("announcementIn", "In how long make announcement").String()
)

func main() {
	fmt.Println(time.Now().UTC())
	db, err := dbOp.Connect("root:@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	var durMax time.Duration = 0
	switch kingpinOp.Read(app, os.Args[1:]) {
	case register.FullCommand():
		err := dbOp.InsertTask(db, *registerDate, *registerTask)
		if err != nil {
			panic(err)
		}
	case announcement.FullCommand():
		n := eventNotifier{
			observers: map[Observer]struct{}{},
		}
		tasks, err := dbOp.ReadTasks(db)
		if err != nil {
			panic(err)
		}
		hr, err := time.ParseDuration(*announcementIn)
		if err != nil {
			panic(err)
		}
		for _, t := range tasks {
			n.Register(&eventObserver{task: t.Task, time: t.Date})
			dur, err := stringToDuration(t.Date)
			if err != nil {
				panic(err)
			}
			if dur > durMax {
				durMax = dur
			}
		}
		stop := time.NewTimer(durMax).C
		tick := time.NewTicker(time.Second).C
		for {
			select {
			case <-stop:
				return
			case <-tick:
				n.Notify(Event{Data: hr})
			}
		}
	}
}
