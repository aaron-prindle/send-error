package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	// "github.com/go-errors/errors"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/clouderrorreporting/v1beta1"
)

func Crash() error {
	return errors.Errorf("this function is supposed to crash")
}

func main() {

	err := Crash()

	if err != nil {
		reportError(err)
	}

}

func reportError(err error) error {
	newError := errors.New(err.Error())
	errMsg := fmt.Sprintf("%+v\n", newError)
	errArray := strings.Split(errMsg, "\n")
	errOutput := []string{}
	// len check?
	errOutput = append(errOutput, errArray[0])
	for i := 1; i < len(errArray)-1; i += 2 {
		errOutput = append(errOutput, fmt.Sprintf("\tat %s (%s)", errArray[i],
			filepath.Base(errArray[i+1])))
	}
	errMsg = strings.Join(errOutput, "\n") + "\n"
	fmt.Printf(errMsg)
	client, err := google.DefaultClient(oauth2.NoContext, clouderrorreporting.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}

	s, err := clouderrorreporting.New(client)
	resp, err := s.Projects.Events.Report("projects/aprindle-vm",
		&clouderrorreporting.ReportedErrorEvent{
			Message: errMsg,
			ServiceContext: &clouderrorreporting.ServiceContext{
				Service: "default",
				Version: "v0.8.0",
			},
		}).Do()
	if err != nil {
		return err
	}
	fmt.Println(resp)
	return nil
}
