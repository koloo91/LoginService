package logging

import "github.com/sirupsen/logrus"

var AppLogger *logrus.Entry

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	AppLogger = logrus.WithField("service", "lgn")
}
