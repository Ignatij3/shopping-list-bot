package mylogger

import (
	"log"
	"os"
	"time"
)

const (
	ERROR string = "ERROR: "
	WARN  string = "WARNING: "
	INFO  string = "INFO: "
)

var logfile *os.File

type Mylogger struct {
	logger *log.Logger
}

// CloseResources must be closed before exiting program in order to ensure proper resourse freeing.
func CloseResources() {
	logfile.Close()
}

// init opens a file to write logs to.
func init() {
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		if err := os.Mkdir("logs", 0755); err != nil {
			log.Fatalf(ERROR+"Could not create \"logs\" directory: %v\n", err)
		}
	}

	var err1 error
	timestamp := time.Now().Format("Jan-_2-15.04.05.000000") // format based on time.StampMicro
	if logfile, err1 = os.Create("logs/shopping-helper_" + timestamp + ".log"); err1 != nil {
		log.Fatalf(ERROR+"Could not create logfile: %v\n", err1)
	}
}

// NewLogger returns a new Mylogger initialized with a specific prefix.
func NewLogger(prefix string) *Mylogger {
	return &Mylogger{logger: log.New(logfile, prefix, log.LstdFlags|log.Lmicroseconds)}
}

func (l *Mylogger) Println(v ...any) {
	l.logger.Println(v...)
}

func (l *Mylogger) Printf(format string, v ...any) {
	l.logger.Printf(format, v...)
}

func (l *Mylogger) Fatalln(v ...any) {
	l.logger.Fatalln(v...)
}

func (l *Mylogger) Fatalf(format string, v ...any) {
	l.logger.Fatalf(format, v...)
}
