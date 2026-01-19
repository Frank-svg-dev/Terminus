package metadata

import (
	"context"
	"sync"

	"k8s.io/klog/v2"
)

// 定义事件类型
type EventType int

const (
	EventUpdate EventType = iota
	EventDelete
)

// UpdateEvent 定义传给 Goroutine 的消息包
type UpdateEvent struct {
	Type      EventType
	ProjectID uint32
	Info      ContainerInfo
}

// AsyncStore 包装了底层的 MemoryStore 和异步通道
type AsyncStore struct {
	// 底层数据存储 (依然需要锁，因为 Exporter 会并发读取)
	data map[uint32]ContainerInfo
	mu   sync.RWMutex

	// 异步通道：缓冲区大小决定了能抗多少突发流量
	updateCh chan UpdateEvent
}

func NewAsyncStore(bufferSize int) *AsyncStore {
	return &AsyncStore{
		data:     make(map[uint32]ContainerInfo),
		updateCh: make(chan UpdateEvent, bufferSize),
	}
}

// ==========================================
// 1. 生产者接口 (给 NRI 调用) - 极速返回
// ==========================================

// TriggerUpdate 触发更新 (非阻塞)
func (s *AsyncStore) TriggerUpdate(id uint32, info ContainerInfo) {
	select {
	case s.updateCh <- UpdateEvent{Type: EventUpdate, ProjectID: id, Info: info}:
		// 写入成功
	default:
		// 通道满了 (极少发生，除非处理逻辑卡死)，打印日志防止阻塞 NRI
		klog.ErrorS(nil, "Metadata update channel full, dropping event", "id", id)
	}
}

// TriggerDelete 触发删除 (非阻塞)
func (s *AsyncStore) TriggerDelete(id uint32) {
	select {
	case s.updateCh <- UpdateEvent{Type: EventDelete, ProjectID: id}:
	default:
		klog.ErrorS(nil, "Metadata update channel full, dropping delete", "id", id)
	}
}

// ==========================================
// 2. 消费者逻辑 (后台 Goroutine)
// ==========================================

// Run 启动消费者循环 (在 main 中 go s.Run())
func (s *AsyncStore) Run(ctx context.Context) {
	klog.Info("Async metadata store worker started")

	for {
		select {
		case <-ctx.Done():
			klog.Info("Async store worker stopped")
			return
		case event := <-s.updateCh:
			s.handleEvent(event)
		}
	}
}

// handleEvent 处理单个事件
func (s *AsyncStore) handleEvent(e UpdateEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch e.Type {
	case EventUpdate:
		s.data[e.ProjectID] = e.Info
		klog.V(4).InfoS("Async updated metadata", "id", e.ProjectID)
		// 【扩展点】在这里可以顺手把数据写到磁盘文件，实现持久化
		// s.persistToDisk()

	case EventDelete:
		delete(s.data, e.ProjectID)
		klog.V(4).InfoS("Async deleted metadata", "id", e.ProjectID)
	}
}

// ==========================================
// 3. 读取接口 (给 Exporter 调用)
// ==========================================

// Get 直接读内存 (依然很快)
func (s *AsyncStore) Get(id uint32) (ContainerInfo, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[id]
	return val, ok
}
