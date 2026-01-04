package project

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
	"wails_study/project/interceptor"
	"wails_study/project/routers"
	"wails_study/project/tcp/tcpServer"

	"github.com/gin-gonic/gin"
)

type ServiceManager struct {
	services []Service
	wg       sync.WaitGroup
}

type Service interface {
	Start() error
	Stop(ctx context.Context) error
	Name() string
}

type HttpService struct {
	addr   string
	engine *gin.Engine
	server *http.Server
}

func NewHttpService(addr string) *HttpService {
	engine := gin.Default()
	engine.Use(interceptor.Log, routers.TimeCost)
	routers.OpenRouters(engine)
	routers.TokenRouters(engine)

	server := &http.Server{
		Addr:    addr,
		Handler: engine,
	}
	return &HttpService{
		addr:   addr,
		engine: engine,
		server: server,
	}
}

func (h *HttpService) Start() error {
	go func() {
		if err := h.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			println("HTTP服务启动/运行失败:", err.Error())
		}
	}()
	return nil
}

func (h *HttpService) Stop(ctx context.Context) error {
	if h.server != nil {
		return h.server.Shutdown(ctx)
	}
	return nil
}

func (h *HttpService) Name() string {
	return "HttpService"
}

type TCPService struct {
	addr     string
	listener net.Listener
	workPool *tcpServer.WorkPool
	stopChan chan struct{}
}

func NewTcpService(addr string) *TCPService {
	return &TCPService{
		addr:     addr,
		stopChan: make(chan struct{}, 1),
	}
}

func (t *TCPService) Start() error {
	listener, err := net.Listen("tcp", t.addr)
	if err != nil {
		if strings.Contains(err.Error(), "address already in use") {
			return fmt.Errorf("TCP服务启动失败: 端口 %s 已被占用，请检查是否有其他程序正在使用此端口", t.addr)
		}
		return fmt.Errorf("TCP服务启动失败: %v", err)
	}
	t.listener = listener
	t.workPool = tcpServer.NewWorkPool(10)

	go func() {
		for {
			select {
			case <-t.stopChan:
				return
			default:
			}

			connection, err := t.listener.Accept()

			if err != nil {
				select {
				case <-t.stopChan:
					return
				default:
					fmt.Printf("接受连接错误: %v\n", err)
					continue
				}
			}
			t.workPool.AddTask(connection)
		}
	}()

	return nil
}

func (t *TCPService) Stop(ctx context.Context) error {
	close(t.stopChan)
	if t.listener != nil {
		err := t.listener.Close()
		if err != nil {
			return err
		}
	}

	if t.workPool != nil {
		t.workPool.Close()
	}
	return nil
}

func (t *TCPService) Name() string {
	return "TCPService"
}

func (sm *ServiceManager) AddService(service Service) {
	sm.services = append(sm.services, service)
}

func (sm *ServiceManager) StartAll() error {
	for _, service := range sm.services {
		if err := service.Start(); err != nil {
			return err
		}
	}
	return nil
}

func (sm *ServiceManager) StopAll() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, service := range sm.services {
		if err := service.Stop(ctx); err != nil {
			println("停止服务错误", service.Name(), ":", err.Error())
		} else {
			println("服务已停止", service.Name())
		}
	}
}
