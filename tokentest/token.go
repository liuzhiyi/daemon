package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"
)

const (
	version string = "1.0"
	proto   string = "tcp"
	addr    string = "127.0.0.1:3000"
)

type DaemonCli struct {
	proto     string
	addr      string
	scheme    string
	in        io.ReadCloser
	out       io.Writer
	err       io.Writer
	transport *http.Transport
}

func NewDaemonCli() *DaemonCli {
	scheme := "http"
	tr := &http.Transport{
	//TLSClientConfig: tlsConfig,
	}
	return &DaemonCli{
		proto:     proto,
		addr:      addr,
		scheme:    scheme,
		in:        os.Stdin,
		out:       os.Stdout,
		err:       os.Stderr,
		transport: tr,
	}
}

func (c *DaemonCli) getMethod(args ...string) (func(...string) error, bool) {
	camelArgs := make([]string, len(args))
	for i, s := range args {
		if len(s) == 0 {
			return nil, false
		}
		camelArgs[i] = strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
	}
	methodName := "Cmd" + strings.Join(camelArgs, "")
	method := reflect.ValueOf(c).MethodByName(methodName)
	if !method.IsValid() {
		return nil, false
	}
	return method.Interface().(func(...string) error), true
}

func (c *DaemonCli) Cmd(args ...string) error {
	if len(args) > 1 {
		method, exists := c.getMethod(args[:2]...)
		if exists {
			return method(args[2:]...)
		}
	}
	if len(args) > 0 {
		method, exists := c.getMethod(args[0])
		if !exists {
			fmt.Fprintf(c.err, "'%s' is not a command. See '--help'.\n", args[0])
			os.Exit(1)
		}
		return method(args[1:]...)
	}
	return c.CmdHelp()
}

func (c *DaemonCli) CmdHelp(args ...string) error {
	if len(args) > 1 {
		method, exists := c.getMethod(args[:2]...)
		if exists {
			method("--help")
			return nil
		}
	}
	if len(args) > 0 {
		method, exists := c.getMethod(args[0])
		if !exists {
			fmt.Fprintf(c.err, "'%s' is not a command. See '--help'.\n", args[0])
			os.Exit(1)
		} else {
			method("--help")
			return nil
		}
	}

	flag.Usage()

	return nil
}

func main() {
	flag.Parse()
	if len(*flHost) == 0 {
		*flHost = addr
	}
	cli := NewDaemonCli()
	if err := cli.Cmd(flag.Args()...); err != nil {
		fmt.Fprint(cli.err, err.Error())
	}
}

type User struct {
	Username string
	Password string
}

func postToken(user User) {
	if ss, err := json.Marshal(user); err != nil {
		panic(err.Error())
	} else {
		r := bytes.NewReader(ss)
		if rsp, err := http.Post("http://127.0.0.1:3000/token", "json", r); err == nil {
			processResult(rsp)
			rsp.Body.Close()
		} else {
			fmt.Println(err.Error())
		}
	}
}

func processResult(rsp *http.Response) {
	str, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(string(str))
	}
}

func getToken() {
	rsp, err := http.Get("http://127.0.0.1:3000/token")
	if err == nil {
		processResult(rsp)
		rsp.Body.Close()
	}
}
