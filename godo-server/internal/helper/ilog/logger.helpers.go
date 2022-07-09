package ilog

func ErrorlnIf(err error, log StdLogger) {
	if err != nil {
		log.Errorln(err)
	}
}
