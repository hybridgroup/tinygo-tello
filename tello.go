// package tello is a client for the Tello drone that works with TinyGo.
package tello

import (
	"encoding/binary"
	"strconv"
	"sync"
	"time"

	"tinygo.org/x/drivers/net"
)

// Tello represents a client to the DJI Tello drone.
type Tello struct {
	adaptor   net.Adapter
	reqAddr   string
	reqPort   string
	respPort  string
	videoPort string
	conn      net.Conn

	cmdMutex  sync.Mutex
	cmdPacket [22]byte

	seq            int16
	rx, ry, lx, ly float32
	throttle       int

	Flying bool
}

func New(a net.Adapter, port string) *Tello {
	n := &Tello{
		adaptor:   a,
		reqAddr:   "192.168.10.1",
		reqPort:   "8889",
		respPort:  port,
		videoPort: "11111",
	}

	return n
}

func (t *Tello) Start() (err error) {
	reqAddr, err := net.ResolveUDPAddr("udp", t.reqAddr+":"+t.reqPort)
	if err != nil {
		return err
	}

	p, err := strconv.Atoi(t.respPort)
	if err != nil {
		return err
	}
	respPort := &net.UDPAddr{Port: p}

	t.conn, err = net.DialUDP("udp", respPort, reqAddr)
	if err != nil {
		return err
	}

	// send connection request using video port
	t.conn.Write([]byte(t.connectionString()))

	go func() {
		for {
			err := t.SendStickCommand()
			if err != nil {
				println("stick command error:", err)
			}
			time.Sleep(20 * time.Millisecond)
		}
	}()

	return err
}

// TakeOff tells the Tello to takeoff
func (t *Tello) TakeOff() (err error) {
	t.cmdMutex.Lock()
	defer t.cmdMutex.Unlock()

	t.createPacketHeader(takeoffCommand, 0x68, 0)
	t.seq++
	binary.LittleEndian.PutUint16(t.cmdPacket[7:], uint16(t.seq))
	binary.LittleEndian.PutUint16(t.cmdPacket[9:], CalculateCRC16(t.cmdPacket[:9]))

	_, err = t.conn.Write(t.cmdPacket[:11])

	return err
}

// Land tells the Tello to land
func (t *Tello) Land() (err error) {
	t.cmdMutex.Lock()
	defer t.cmdMutex.Unlock()

	t.createPacketHeader(landCommand, 0x68, 1)
	t.seq++
	binary.LittleEndian.PutUint16(t.cmdPacket[7:], uint16(t.seq))
	t.cmdPacket[9] = 0x00
	binary.LittleEndian.PutUint16(t.cmdPacket[10:], CalculateCRC16(t.cmdPacket[:10]))

	_, err = t.conn.Write(t.cmdPacket[:12])

	return err
}

// Up tells the drone to ascend. Pass in an int from 0-100.
func (t *Tello) Up(val int) error {
	t.cmdMutex.Lock()
	defer t.cmdMutex.Unlock()

	t.ly = float32(val) / 100.0
	return nil
}

// Down tells the drone to descend. Pass in an int from 0-100.
func (t *Tello) Down(val int) error {
	t.cmdMutex.Lock()
	defer t.cmdMutex.Unlock()

	t.ly = float32(val) / 100.0 * -1
	return nil
}

// Forward tells the drone to go forward. Pass in an int from 0-100.
func (t *Tello) Forward(val int) error {
	t.cmdMutex.Lock()
	defer t.cmdMutex.Unlock()

	t.ry = float32(val) / 100.0
	return nil
}

// Backward tells drone to go in reverse. Pass in an int from 0-100.
func (t *Tello) Backward(val int) error {
	t.cmdMutex.Lock()
	defer t.cmdMutex.Unlock()

	t.ry = float32(val) / 100.0 * -1
	return nil
}

// Right tells drone to go right. Pass in an int from 0-100.
func (t *Tello) Right(val int) error {
	t.cmdMutex.Lock()
	defer t.cmdMutex.Unlock()

	t.rx = float32(val) / 100.0
	return nil
}

// Left tells drone to go left. Pass in an int from 0-100.
func (t *Tello) Left(val int) error {
	t.cmdMutex.Lock()
	defer t.cmdMutex.Unlock()

	t.rx = float32(val) / 100.0 * -1
	return nil
}

// Clockwise tells drone to rotate in a clockwise direction. Pass in an int from 0-100.
func (t *Tello) Clockwise(val int) error {
	t.cmdMutex.Lock()
	defer t.cmdMutex.Unlock()

	t.lx = float32(val) / 100.0
	return nil
}

// CounterClockwise tells drone to rotate in a counter-clockwise direction.
// Pass in an int from 0-100.
func (t *Tello) CounterClockwise(val int) error {
	t.cmdMutex.Lock()
	defer t.cmdMutex.Unlock()

	t.lx = float32(val) / 100.0 * -1
	return nil
}

func (t *Tello) SendStickCommand() (err error) {
	t.cmdMutex.Lock()
	defer t.cmdMutex.Unlock()

	t.createPacketHeader(stickCommand, 0x60, 11)
	binary.LittleEndian.PutUint16(t.cmdPacket[7:], 0x00) // seq = 0

	// RightX center=1024 left =364 right =-364
	axis1 := int16(660.0*t.rx + 1024.0)

	// RightY down =364 up =-364
	axis2 := int16(660.0*t.ry + 1024.0)

	// LeftY down =364 up =-364
	axis3 := int16(660.0*t.ly + 1024.0)

	// LeftX left =364 right =-364
	axis4 := int16(660.0*t.lx + 1024.0)

	// speed control
	axis5 := int16(t.throttle)

	packedAxis := int64(axis1)&0x7FF | int64(axis2&0x7FF)<<11 | 0x7FF&int64(axis3)<<22 | 0x7FF&int64(axis4)<<33 | int64(axis5)<<44
	t.cmdPacket[9] = byte(0xFF & packedAxis)
	t.cmdPacket[10] = byte(packedAxis >> 8 & 0xFF)
	t.cmdPacket[11] = byte(packedAxis >> 16 & 0xFF)
	t.cmdPacket[12] = byte(packedAxis >> 24 & 0xFF)
	t.cmdPacket[13] = byte(packedAxis >> 32 & 0xFF)
	t.cmdPacket[14] = byte(packedAxis >> 40 & 0xFF)

	now := time.Now()
	t.cmdPacket[15] = byte(now.Hour())
	t.cmdPacket[16] = byte(now.Minute())
	t.cmdPacket[17] = byte(now.Second())
	t.cmdPacket[18] = byte(now.UnixNano() / int64(time.Millisecond) & 0xff)
	t.cmdPacket[19] = byte(now.UnixNano() / int64(time.Millisecond) >> 8)

	binary.LittleEndian.PutUint16(t.cmdPacket[20:], CalculateCRC16(t.cmdPacket[:20]))

	_, err = t.conn.Write(t.cmdPacket[:22])
	return
}

func (t *Tello) createPacketHeader(cmd int16, pktType byte, len int16) (err error) {
	l := len + 11

	t.cmdPacket[0] = byte(messageStart)
	binary.LittleEndian.PutUint16(t.cmdPacket[1:], uint16(l<<3))
	t.cmdPacket[3] = CalculateCRC8(t.cmdPacket[0:3])
	t.cmdPacket[4] = pktType
	binary.LittleEndian.PutUint16(t.cmdPacket[5:], uint16(cmd))

	return nil
}

func (t *Tello) connectionString() string {
	x, _ := strconv.Atoi(t.videoPort)
	msg := []byte("conn_req:xx")
	binary.LittleEndian.PutUint16(msg[9:], uint16(x))
	return string(msg)
}

func validatePitch(val int) int {
	if val > 100 {
		return 100
	} else if val < 0 {
		return 0
	}

	return val
}
