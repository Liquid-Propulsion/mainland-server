package lockout

import (
	"log"
	"net"
	"time"

	"github.com/Liquid-Propulsion/mainland-server/config"
)

type Lockout struct {
	lastMessageRecieved time.Time
}

func New() *Lockout {
	lockout := new(Lockout)
	lockout.lastMessageRecieved = time.Unix(0, 0)
	return lockout
}

func (lockout *Lockout) Run() {
	socket, err := net.ListenPacket("udp", config.CurrentConfig.Lockout.ListenAddr)
	if err != nil {
		log.Fatalf("Lockout UDP Socket could not be opened: %s", err)
	}
	defer socket.Close()

	log.Printf("Successfully initialized lockout system, listening on %s", config.CurrentConfig.Lockout.ListenAddr)
	for {
		buf := make([]byte, 1024)
		len, addr, err := socket.ReadFrom(buf)
		if err != nil {
			continue
		}
		buf = buf[:len]
		// magic packet is [0xFF, 0xFA, 0xF3, 0xA3]
		if buf[0] == 0xFF && buf[1] == 0xFA && buf[2] == 0xF3 && buf[3] == 0xA3 {
			lockout.lastMessageRecieved = time.Now()
			socket.WriteTo(buf, addr)
		}
	}
}

func (lockout *Lockout) LockedOut() bool {
	if config.CurrentConfig.Lockout.Enabled {
		return time.Since(lockout.lastMessageRecieved) > time.Millisecond*40
	}
	return false
}
