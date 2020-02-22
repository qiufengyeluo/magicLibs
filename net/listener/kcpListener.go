package listener

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/yamakiller/magicLibs/net/contractor"
	"github.com/yamakiller/mgokcp/mkcp"
)

const (
	kcpHeaderLength = 24
	statIdle        = int32(1)
	statRun         = int32(2)
)

//KCPData kcp读取数据包
type KCPData struct {
	Data   []byte
	Length int
}

//SpawnKCPListener create an kcp(udp) listener
func SpawnKCPListener(l *net.UDPConn, mtu int) *KCPListener {

	return &KCPListener{_l: l,
		_buffer: make([]byte, mtu),
		_mtu:    mtu,
		_pool: sync.Pool{
			New: func() interface{} {
				return &KCPData{
					Data:   make([]byte, mtu),
					Length: 0,
				}
			}},
		_convUsed: make(map[uint32]*KCPConn),
	}
}

func output(buff []byte, user interface{}) int32 {
	conn := user.(*KCPConn)
	fmt.Println(conn._parent._l.WriteToUDP(buff, conn._addr))
	return 0
}

//KCPListener KCP Listener
type KCPListener struct {
	RecvWndSize   int32
	SendWndSize   int32
	RecvQueueSize int32
	NoDelay       int32
	Interval      int32
	Resend        int32
	Nc            int32
	RxMinRto      int32
	FastResend    int32
	Reser         contractor.Conv

	_l        *net.UDPConn
	_buffer   []byte
	_mtu      int
	_convSz   int
	_convUsed map[uint32]*KCPConn
	_convSync sync.Mutex
	_pool     sync.Pool
	_wg       sync.WaitGroup
}

//Accept kcp accept connection
func (slf *KCPListener) Accept([]interface{}) (interface{}, error) {
	n, addr, err := slf._l.ReadFromUDP(slf._buffer)
	if err != nil {
		return nil, err
	}

	//TODO:有问题的数据包
	if n < kcpHeaderLength {
		if slf.Reser == nil {
			return nil, nil
		}

		if n != 9 {
			return nil, nil
		}

		switch slf._buffer[0] {
		case 0x01:
			id, key, err := slf.Reser.Reserve(binary.BigEndian.Uint64(slf._buffer[1:]))
			if err != nil {
				return nil, nil
			}
			slf._buffer[0] = 0x01
			binary.BigEndian.PutUint32(slf._buffer[1:], id)
			binary.BigEndian.PutUint64(slf._buffer[5:], key)
			slf._l.WriteToUDP(slf._buffer[:13], addr)
			return nil, nil
		case 0x02:
			if !slf.Reser.Confirm(slf._buffer[1:]) {
				slf._buffer[0] = 0x12
				slf._l.WriteToUDP(slf._buffer[:1], addr)
				return nil, nil
			}
			slf._buffer[0] = 0x02
			slf._l.WriteToUDP(slf._buffer[:1], addr)
			return nil, nil
		default:
			return nil, nil
		}
	}

	conv := mkcp.GetConv(slf._buffer)

	if slf.Reser != nil {
		if !slf.Reser.Authorized(conv) {
			//未预约ID无法通信
			return nil, nil
		}
	}

	conn := slf.get(conv)
	//TODO: 连接不存在重新创建
	isAccept := false
	if conn == nil {
		slf._wg.Add(1)
		conn = &KCPConn{
			_id:     conv,
			_addr:   addr,
			_parent: slf,
			/*_sync: mutex.SpinLock{
				Deplay: time.Millisecond,
				Check:  8,
			}，*/
			_wg: &slf._wg,
		}

		conn._kcp = mkcp.New(conv, conn)
		conn._kcp.WithOutput(output)
		conn._kcp.WndSize(slf.RecvWndSize, slf.SendWndSize)
		conn._kcp.NoDelay(slf.NoDelay, slf.Interval, slf.Resend, slf.Nc)
		conn._kcp.SetMTU(int32(slf._mtu))
		conn._recv = make(chan *KCPData, slf.RecvQueueSize)
		conn._destored = make(chan bool, 1)
		if slf.RxMinRto > 0 {
			conn._kcp.SetRxMinRto(slf.RxMinRto)
		}

		if slf.FastResend > 0 {
			conn._kcp.SetFastResend(slf.FastResend)
		}

		conn._addr = addr
		slf.join(conn)
		isAccept = true
	}

	conn._sync.Lock()
	conn._addr = addr
	if conn._kcp != nil {
		conn._kcp.Input(slf._buffer, int32(n))
	}
	conn._sync.Unlock()

	for {
		conn._sync.Lock()
		if conn._kcp == nil {
			conn._sync.Unlock()
			break
		}

		n := int(conn._kcp.Recv(slf._buffer, int32(len(slf._buffer))))
		if n < 0 {
			conn._sync.Unlock()
			break
		}
		conn._sync.Unlock()

		data := slf._pool.Get().(*KCPData)
		copy(data.Data, slf._buffer[:n])
		data.Length = n
		select {
		case <-conn._destored:
			slf._pool.Put(data)
			goto exit
		default:
		}

		conn._recv <- data

		continue
	exit:
		break
	}

	if isAccept {
		return conn, nil
	}
	return nil, nil
}

//Update 更新连接状态
func (slf *KCPListener) Update(tss int64) int {
	current := uint32(tss & 0xffffffff)
	var convs []uint32
	slf._convSync.Lock()
	n := len(slf._convUsed)
	if n > 0 {
		convs = make([]uint32, n)
		i := 0
		for _, v := range slf._convUsed {
			convs[i] = v._id
			i++
		}
	}
	slf._convSync.Unlock()

	if convs == nil {
		return -1
	}

	for _, v := range convs {
		con := slf.get(v)
		if con == nil {
			continue
		}

		if con._closed {
			con._sync.Lock()
			mkcp.Free(con._kcp)
			con._kcp = nil
			con._sync.Unlock()
			close(con._destored)

			for {
				select {
				case v := <-con._recv:
					slf._pool.Put(v)
					continue
				default:
				}
				close(con._recv)
				break
			}

			slf.unjoin(con._id)
			con._wg.Done()
			fmt.Println("remove complate")
			continue
		}

		con._sync.Lock()
		if con._kcp.Check(current) >= con._kcp.Check(current) {
			con._kcp.Update(current)
		}
		con._sync.Unlock()
	}

	return 0
}

//Addr Returns  address
func (slf *KCPListener) Addr() net.Addr {
	return slf._l.LocalAddr()
}

//Close close listener
func (slf *KCPListener) Close() error {
	if err := slf._l.Close(); err != nil {
		return err
	}
	//设置所有连接关闭状态
	slf._convSync.Lock()
	for _, v := range slf._convUsed {
		if v._closed {
			continue
		}
		v._closed = true
	}
	slf._convSync.Unlock()
	slf._wg.Wait()
	return nil
}

func (slf *KCPListener) reallocate() {
	slf._buffer = make([]byte, slf._mtu*2)
}

//ToString ....
func (slf *KCPListener) ToString() string {
	return "kcp(udp) listener"
}

func (slf *KCPListener) join(c *KCPConn) {
	slf._convSync.Lock()
	defer slf._convSync.Unlock()
	slf._convUsed[c._id] = c
}

func (slf *KCPListener) unjoin(id uint32) {
	slf._convSync.Lock()
	defer slf._convSync.Unlock()
	if _, ok := slf._convUsed[id]; ok {
		delete(slf._convUsed, id)
	}
}

func (slf *KCPListener) get(id uint32) *KCPConn {
	slf._convSync.Lock()
	defer slf._convSync.Unlock()

	if v, ok := slf._convUsed[id]; ok {
		return v
	}

	return nil
}

//KCPConn KCP连接者
type KCPConn struct {
	_id        uint32
	_addr      *net.UDPAddr
	_kcp       *mkcp.KCP
	_parent    *KCPListener
	_deathtime *time.Time
	_recv      chan *KCPData
	_destored  chan bool
	_closed    bool
	_sync      sync.Mutex
	_wg        *sync.WaitGroup
}

//SetReadDeadline 设置读取超时
func (slf *KCPConn) SetReadDeadline(t time.Time) {
	slf._deathtime = &t
}

//Recv 接收数据
func (slf *KCPConn) Recv(buffer []byte, size int32) (int32, error) {
	if slf._deathtime != nil {
		timeout := slf._deathtime.Sub(time.Now())
		select {
		case <-slf._destored:
			return -1, errors.New("closed")
		case <-time.After(timeout):
			return -1, errors.New("reader timeout")
		case data := <-slf._recv:
			if data == nil {
				return -1, errors.New("closed")
			}

			if int32(data.Length) > size {
				return -1, errors.New("buffer overflow")
			}

			n := int32(data.Length)
			copy(buffer, data.Data[:n])
			return n, nil
		}
	} else {
		select {
		case <-slf._destored:
			return -1, errors.New("closed")
		case data := <-slf._recv:
			if data == nil {
				return -1, errors.New("closed")
			}

			if int32(data.Length) > size {
				return -1, errors.New("buffer overflow")
			}

			n := int32(data.Length)
			copy(buffer, data.Data[:n])
			return n, nil
		}
	}
}

//Send 发送数据
func (slf *KCPConn) Write(buffer []byte, size int32) (int32, error) {
	slf._sync.Lock()
	if slf._kcp == nil {
		slf._sync.Unlock()
		return -1, errors.New("closed")
	}
	n, err := slf._kcp.Send(buffer, size)
	slf._sync.Unlock()
	return n, err
}

// LocalAddr returns the local network address.
func (slf *KCPConn) LocalAddr() net.Addr {
	p := unsafe.Pointer(slf._parent._l)
	l := (*net.UDPConn)(atomic.LoadPointer(&p))
	if l == nil {
		return &net.UDPAddr{}
	}

	return l.LocalAddr()
}

// RemoteAddr returns the remote network address.
func (slf *KCPConn) RemoteAddr() net.Addr {
	p := unsafe.Pointer(slf._addr)
	addr := (*net.UDPAddr)(atomic.LoadPointer(&p))
	if addr == nil {
		return &net.UDPAddr{}
	}

	return addr
}

//Close 关闭连接
func (slf *KCPConn) Close() error {
	if slf._closed {
		return errors.New("Repeatedly closed")
	}
	slf._closed = true
	return nil
}
