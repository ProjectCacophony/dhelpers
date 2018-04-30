package components

import (
	"crypto/tls"
	"net"
	"strings"

	"os"

	"github.com/globalsign/mgo"
	"gitlab.com/Cacophony/dhelpers/cache"
)

type mgoLogger struct {
}

func (mgol mgoLogger) Output(calldepth int, s string) error {
	cache.GetLogger().WithField("module", "mongodb").Info(s)
	return nil
}

// InitMongoDB initialises the MongoDB session
// reads the MongoDB URL from the environment variable MONGODB_URL
// reads the MongoDB Database from the environment variable MONGODB_DATABASE
func InitMongoDB() (err error) {
	mgoL := new(mgoLogger)
	mgo.SetLogger(mgoL)

	newURL := strings.TrimSuffix(os.Getenv("MONGODB_URL"), "?ssl=true")
	newURL = strings.Replace(newURL, "ssl=true&", "", -1)

	dialInfo, err := mgo.ParseURL(newURL)
	if err != nil {
		return err
	}

	// setup TLS if we use SSL
	if newURL != os.Getenv("MONGODB_URL") {
		tlsConfig := &tls.Config{}
		tlsConfig.InsecureSkipVerify = true

		dialInfo.DialServer = func(addr *mgo.ServerAddr) (conn net.Conn, err error) {
			conn, err = tls.Dial("tcp", addr.String(), tlsConfig)
			return conn, err
		}
	}

	mDbSession, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		return err
	}

	mDbSession.SetMode(mgo.Monotonic, true)
	mDbSession.SetSafe(&mgo.Safe{WMode: "majority"})

	cache.SetMgo(mDbSession.DB(os.Getenv("MONGODB_DATABASE")))
	return nil
}
