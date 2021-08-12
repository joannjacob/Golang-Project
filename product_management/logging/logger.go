package logging

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func InitializeLogging() {
	logFile := os.Getenv("LOGGER_OUTPUT")
	logLevel := os.Getenv("LOGGER_LEVEL")

	file, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.WithFields(log.Fields{"method": "InitializeLogging()", "error": err.Error()}).
			Error("Error in opening log file")
	}
	log.SetOutput(file)

	level, err := log.ParseLevel(logLevel)
	if err != nil {
		log.WithFields(log.Fields{"method": "InitializeLogging()", "error": err.Error()}).
			Error("Could Not Set Log Level")
	}
	log.SetLevel(level)

	log.SetFormatter(&log.JSONFormatter{})

}


// Logger struct.
// type Logger struct {
// 	Log     *logrus.Logger
// 	Logfile *os.File
// }

// func NewLogger() (*Logger, error) {
// logFile = os.Getenv("app_logger_output")
// logLevel = os.Getenv("app_logger_level")

// log := logrus.New()

// log.Level, err := logrus.ParseLevel("debug")
// if err != nil {
// 	return nil, err
// }

// log.Out = file

// log.Formatter = &logrus.JSONFormatter{}

// return &Logger{Log: log, Logfile: file}, nil
// }

// var log = logrus.New()
// log.Formatter = new(logrus.JSONFormatter)
// log.Formatter = new(logrus.TextFormatter)                     //default
// log.Formatter.(*logrus.TextFormatter).DisableColors = true    // remove colors
// log.Formatter.(*logrus.TextFormatter).DisableTimestamp = true // remove timestamp from test output
// log.Level = logrus.TraceLevel
// log.Out = os.Stdout

// file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY, 0666)
// if err == nil {
// 	log.Out = file
// } else {
// 	log.Info("Failed to log to file, using default stderr")
// }

// defer func() {
//     err := recover()
//     if err != nil {
//         entry := err.(*logrus.Entry)
//         log.WithFields(logrus.Fields{
//             "omg":         true,
//             "err_animal":  entry.Data["animal"],
//             "err_size":    entry.Data["size"],
//             "err_level":   entry.Level,
//             "err_message": entry.Message,
//             "number":      100,
//         }).Error("The ice breaks!") // or use Fatal() to force the process to exit with a nonzero code
//     }
// }()

// log.WithFields(logrus.Fields{
//     "animal": "walrus",
//     "number": 0,
// }).Trace("Went to the beach")

// log.WithFields(logrus.Fields{
//     "animal": "walrus",
//     "number": 8,
// }).Debug("Started observing beach")

// log.WithFields(logrus.Fields{
//     "animal": "walrus",
//     "size":   10,
// }).Info("A group of walrus emerges from the ocean")

// log.WithFields(logrus.Fields{
//     "omg":    true,
//     "number": 122,
// }).Warn("The group's number increased tremendously!")

// log.WithFields(logrus.Fields{
//     "temperature": -4,
// }).Debug("Temperature changes")

// log.WithFields(logrus.Fields{
//     "animal": "orca",
//     "size":   9009,
// }).Panic("It's over 9000!")
