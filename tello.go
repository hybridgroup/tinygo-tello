// package tello is a client for the Tello drone that works with TinyGo.
package tello

import (
	"encoding/binary"
	"net"
	"strconv"
	"sync"
	"time"
)

// Tello represents a client to the DJI Tello drone.
type Tello struct {
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

func New(port string) *Tello {
	n := &Tello{
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

	respAddr, err := net.ResolveUDPAddr("udp", ":"+t.respPort)
	if err != nil {
		return err
	}

	t.conn, err = net.DialUDP("udp", respAddr, reqAddr)
	if err != nil {
		return err
	}

	// send connection request using video port
	if _, err := t.conn.Write([]byte(t.connectionString())); err != nil {
		return err
	}

	go func() {
		for {
			err := t.SendStickCommand()
			if err != nil {
				println("stick command error:", err)
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	return nil
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
	t.ly = float32(val) / 100.0
	return nil
}

// Down tells the drone to descend. Pass in an int from 0-100.
func (t *Tello) Down(val int) error {
	t.ly = float32(val) / 100.0 * -1.0
	return nil
}

// Forward tells the drone to go forward. Pass in an int from 0-100.
func (t *Tello) Forward(val int) error {
	t.ry = float32(val) / 100.0
	return nil
}

// Backward tells drone to go in reverse. Pass in an int from 0-100.
func (t *Tello) Backward(val int) error {
	t.ry = float32(val) / 100.0 * -1.0
	return nil
}

// Right tells drone to go right. Pass in an int from 0-100.
func (t *Tello) Right(val int) error {
	t.rx = float32(val) / 100.0
	return nil
}

// Left tells drone to go left. Pass in an int from 0-100.
func (t *Tello) Left(val int) error {
	t.rx = float32(val) / 100.0 * -1.0
	return nil
}

// Clockwise tells drone to rotate in a clockwise direction. Pass in an int from 0-100.
func (t *Tello) Clockwise(val int) error {
	t.lx = float32(val) / 100.0
	return nil
}

// CounterClockwise tells drone to rotate in a counter-clockwise direction.
// Pass in an int from 0-100.
func (t *Tello) CounterClockwise(val int) error {
	t.lx = float32(val) / 100.0 * -1.0
	return nil
}

// Throw & Go support
func (t *Tello) ThrowTakeOff() error {
	t.cmdMutex.Lock()
	defer t.cmdMutex.Unlock()

	t.createPacketHeader(throwtakeoffCommand, 0x48, 0)
	t.seq++
	binary.LittleEndian.PutUint16(t.cmdPacket[7:], uint16(t.seq))
	binary.LittleEndian.PutUint16(t.cmdPacket[9:], CalculateCRC16(t.cmdPacket[:9]))

	_, err := t.conn.Write(t.cmdPacket[:11])

	return err
}

// PalmLand tells drone to come in for a landing on the palm of your hand.
func (t *Tello) PalmLand() error {
	t.cmdMutex.Lock()
	defer t.cmdMutex.Unlock()

	t.createPacketHeader(palmLandCommand, 0x68, 1)
	t.seq++
	binary.LittleEndian.PutUint16(t.cmdPacket[7:], uint16(t.seq))
	t.cmdPacket[9] = 0x00
	binary.LittleEndian.PutUint16(t.cmdPacket[10:], CalculateCRC16(t.cmdPacket[:10]))

	_, err := t.conn.Write(t.cmdPacket[:12])

	return err
}

// Flip tells drone to flip
func (t *Tello) Flip(direction FlipType) error {
	t.cmdMutex.Lock()
	defer t.cmdMutex.Unlock()

	t.createPacketHeader(flipCommand, 0x70, 1)
	t.seq++
	binary.LittleEndian.PutUint16(t.cmdPacket[7:], uint16(t.seq))
	t.cmdPacket[9] = byte(direction)
	binary.LittleEndian.PutUint16(t.cmdPacket[10:], CalculateCRC16(t.cmdPacket[:10]))

	_, err := t.conn.Write(t.cmdPacket[:12])

	return err
}

// StartVideo tells Tello to send start info (SPS/PPS) for video stream.
func (t *Tello) StartVideo() (err error) {
	t.cmdMutex.Lock()
	defer t.cmdMutex.Unlock()

	t.createPacketHeader(videoStartCommand, 0x60, 0)
	binary.LittleEndian.PutUint16(t.cmdPacket[7:], 0x00) // seq = 0
	binary.LittleEndian.PutUint16(t.cmdPacket[9:], CalculateCRC16(t.cmdPacket[:9]))

	_, err = t.conn.Write(t.cmdPacket[:11])

	return err
}

func (t *Tello) SendStickCommand() (err error) {
	t.cmdMutex.Lock()
	defer t.cmdMutex.Unlock()

	t.createPacketHeader(stickCommand, 0x60, 11)
	binary.LittleEndian.PutUint16(t.cmdPacket[7:], 0x00) // seq = 0

	// All axes range from 364 to 1684
	// RightX left =364 right =1684
	axis1 := int16(660.0*t.rx + 1024.0)

	// RightY down =364 up =1684
	axis2 := int16(660.0*t.ry + 1024.0)

	// LeftY down =364 up =1684
	axis3 := int16(660.0*t.ly + 1024.0)

	// LeftX left =364 right =1684
	axis4 := int16(660.0*t.lx + 1024.0)

	// speed control
	axis5 := int16(t.throttle)

	packedAxis := int64(axis1)&0x7FF | int64(axis2&0x7FF)<<11 | int64(axis3&0x7FF)<<22 | int64(axis4&0x7FF)<<33 | int64(axis5)<<44
	t.cmdPacket[9] = byte(0xFF & packedAxis)
	t.cmdPacket[10] = byte(packedAxis >> 8)
	t.cmdPacket[11] = byte(packedAxis >> 16)
	t.cmdPacket[12] = byte(packedAxis >> 24)
	t.cmdPacket[13] = byte(packedAxis >> 32)
	t.cmdPacket[14] = byte(packedAxis >> 40)

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
