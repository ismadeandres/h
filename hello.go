package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"log/syslog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

//#import "github.com/ismadeandres/goini"

func main() {

	//path search
	flag.Parse()
	root := flag.Arg(0)
	err3 := filepath.Walk(root, visit)
	fmt.Printf("filepath.Walk() returned %v\n", err3)

	// Now from anywhere else in your program, you can use this:
	log.Print("Hello Logs!")

	//cmdRun
	var (
		cmdOut []byte
		err    error
	)
	fmt.Println("Hello, new gopher!")
	cmd := "ls"
	args := []string{"-l", "-a"}
	if cmdOut, err = exec.Command(cmd, args...).Output(); err != nil {
		fmt.Println(err)
		fmt.Println("Ha habido algun error")
		os.Exit(1)
	}
	fmt.Println(string(cmdOut))
	//#	fini,_ := ini.Load("/tmp/example.ini")
	//#fmt.Println( ini.GetString("Wine","Grape" ) )

	cmd2 := exec.Command("ls", "/tmp")
	var waitStatus syscall.WaitStatus
	if err := cmd2.Run(); err != nil {
		printError(err)
		// Did the command fail because of an unsuccessful exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus = exitError.Sys().(syscall.WaitStatus)
			printOutput([]byte(fmt.Sprintf("%d", waitStatus.ExitStatus())))
		}
	} else {
		// Command was successful
		waitStatus = cmd2.ProcessState.Sys().(syscall.WaitStatus)
		printOutput([]byte(fmt.Sprintf("%d", waitStatus.ExitStatus())))
	}
}

func init() {

	f, err := os.OpenFile("/var/log/hello.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println("This is a test log entry")

	//log2Syslog
	logwriter, e := syslog.New(syslog.LOG_NOTICE, "hello")
	if e == nil {
		log.SetOutput(logwriter)
		log.Print("Initialization Complete!")
	}
	mWriter := io.MultiWriter(f, logwriter)

	data := []byte("Hello World MultiWriter!")

	n, err := mWriter.Write(data)

	if err == nil {
		fmt.Printf("Multi write %d bytes to two files simultaneously.\n", n)
	}

}

func printCommand(cmd *exec.Cmd) {
	fmt.Printf("==> Executing: %s\n", strings.Join(cmd.Args, " "))
}

func printError(err error) {
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))
	}
}

func printOutput(outs []byte) {
	if len(outs) > 0 {
		fmt.Printf("==> Output: %s\n", string(outs))
	}
}

func visit(path string, f os.FileInfo, err error) error {
	fmt.Printf("Visited: %s\n", path)
	return nil
}
