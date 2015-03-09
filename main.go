package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/kardianos/service"
)

var logger service.Logger

// Program structures.
//  Define Start and Stop methods.
type program struct {
	exit chan struct{}
}

func (p *program) Start(s service.Service) error {
	if service.Interactive() {
		logger.Info("Running in terminal.")
	} else {
		logger.Info("Running under service manager.")
	}
	p.exit = make(chan struct{})

	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}

func (p *program) run() error {
	logger.Infof("I'm running %v.", service.Platform())
	httpServer()
	return nil
}

func (p *program) Stop(s service.Service) error {
	// Any work in Stop should be quick, usually a few seconds at most.
	logger.Info("I'm Stopping!")
	close(p.exit)
	return nil
}

func httpServer() {
	http.HandleFunc("/token", TokenHandle)
	if err := http.ListenAndServe(":3000", nil); err != nil {
		logger.Error(err.Error())
		panic(err.Error())
	}
}

func TokenHandle(w http.ResponseWriter, req *http.Request) {
	var rsp Rsp
	if req.Method == "POST" {
		str, err := ioutil.ReadAll(req.Body)
		if err != nil {
			logger.Info(err.Error())
			rsp.Code = "100"
			rsp.Msg = "系统错误"
		}
		fmt.Println(string(str))
		user := new(User)
		if err := json.Unmarshal(str, user); err != nil {
			logger.Info(err.Error())
			rsp.Code = "201"
			rsp.Msg = "数据解析错误"
		} else {
			if user.Username == "admin" && user.Password == "123456" {
				rsp.Code = "200"
				rsp.Msg = "登录成功"
			} else {
				rsp.Code = "202"
				rsp.Msg = "用户名或密码错误"
			}
		}

	} else {
		rsp.Code = "203"
		rsp.Msg = "请使用post访问"
	}

	if ss, err := json.Marshal(rsp); err != nil {
		logger.Error(err.Error())
		panic(err.Error())
	} else {
		io.WriteString(w, string(ss))
	}
}

type User struct {
	Username string
	Password string
}

type Rsp struct {
	Code   string
	Msg    string
	Object interface{}
}

// Service setup.
//   Define service config.
//   Create the service.
//   Setup the logger.
//   Handle service controls (optional).
//   Run the service.
func main() {
	//svcFlag := flag.String("service", "", "Control the system service.")
	flag.Parse()

	svcConfig := &service.Config{
		Name:        "daemon",
		DisplayName: "Go Service Test",
		Description: "This is a test Go service.",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	errs := make(chan error, 5)
	logger, err = s.Logger(errs)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			err := <-errs
			if err != nil {
				log.Print(err)
			}
		}
	}()
	args := flag.Args()

	if len(args) == 1 {
		err := service.Control(s, args[0])
		if err != nil {
			log.Printf("Valid actions: %q\n", service.ControlAction)
			log.Fatal(err)
		}
		return
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}

const usageTemplate = `usage: daemon command [arguments]

The commands are:
{{range .}}
    {{.Name | printf "%-11s"}} {{.Short}}{{end}}

Use "revel help [command]" for more information.
`

var helpTemplate = `usage: revel {{.UsageLine}}
{{.Long}}
`

type Command struct {
	Run                    func(args []string)
	UsageLine, Short, Long string
}

var commands = []*Command{}

func usage(exitCode int) {
	tmpl(os.Stderr, usageTemplate, commands)
	os.Exit(exitCode)
}

func tmpl(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	template.Must(t.Parse(text))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}
