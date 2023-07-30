package http

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/sqjian/go-kit/log"
	"io"
	"net/http"
	"net/url"
	"sync"
)

var (
	defaultUpgrader = &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	defaultDialer = websocket.DefaultDialer
)

type IncomeInterceptor func(int, []byte) (int, []byte)
type OutcomeInterceptor func(int, []byte) (int, []byte)

func newDefaultWebsocketProxyConfig() *WebsocketProxy {
	return &WebsocketProxy{
		Backend:            nil,
		Upgrader:           defaultUpgrader,
		Dialer:             defaultDialer,
		Logger:             func() log.Log { inst, _ := log.NewLogger(log.WithLevel("dummy")); return inst }(),
		IncomeInterceptor:  func(msgType int, msg []byte) (int, []byte) { return msgType, msg },
		OutcomeInterceptor: func(msgType int, msg []byte) (int, []byte) { return msgType, msg },
	}
}

type WebsocketProxy struct {
	Backend            *url.URL
	Upgrader           *websocket.Upgrader
	Dialer             *websocket.Dialer
	Logger             log.Log
	IncomeInterceptor  IncomeInterceptor
	OutcomeInterceptor OutcomeInterceptor
}

func (wp *WebsocketProxy) Init(addr string, opts ...WebsocketProxyOptionFunc) (*WebsocketProxy, error) {
	target, targetErr := url.Parse(addr)
	if targetErr != nil {
		return nil, targetErr
	}

	wpInst := newDefaultWebsocketProxyConfig()
	wpInst.Backend = target

	for _, opt := range opts {
		opt(wpInst)
	}

	return wpInst, nil
}

func (wp *WebsocketProxy) WebsocketProxyHandle(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	connBackend, connBackendResp, connBackendRespErr := wp.Dialer.Dial(wp.Backend.String(), nil)
	if connBackendRespErr != nil {
		wp.Logger.Infof("websocket proxy: couldn't dial to remote backend url:%v %s", wp.Backend.String(), connBackendRespErr)
		if connBackendResp != nil {
			if copyResponseErr := copyResponse(w, connBackendResp); copyResponseErr != nil {
				wp.Logger.Errorf("websocket proxy: couldn't write response after failed remote backend handshake: %s", copyResponseErr)
				w.Write([]byte(fmt.Errorf("connBackendRespErr:%w,copyResponseErr:%w", connBackendRespErr, copyResponseErr).Error()))
			} else {
				wp.Logger.Errorf("websocket proxy: write response successfully after failed remote backend handshake: %s", copyResponseErr)
				w.Write([]byte(fmt.Errorf("connBackendRespErr:%w", connBackendRespErr).Error()))
			}
		} else {
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			w.Write([]byte(fmt.Errorf("connBackendRespErr:%w", connBackendRespErr).Error()))
		}
		return
	} else {
		wp.Logger.Infof("websocket proxy: dial to remote backend url:%v successfully", wp.Backend.String())
	}
	defer connBackend.Close()

	connFrontend, connFrontendUpgradeErr := wp.Upgrader.Upgrade(w, r, nil)
	if connFrontendUpgradeErr != nil {
		wp.Logger.Errorf("websocket proxy: couldn't upgrade %s", connFrontendUpgradeErr)
		return
	} else {
		wp.Logger.Errorf("websocket proxy: upgrade successfully")
	}
	defer connFrontend.Close()

	processInComeMsg := func(wg *sync.WaitGroup, connFrontend, connBackend *websocket.Conn) {
		defer wg.Done()

		for {
			msgType, msg, readMessageErr := connFrontend.ReadMessage()
			if readMessageErr != nil {
				wp.Logger.Errorf("readMessage failed,writeMessageErr:%v", readMessageErr.Error())
				m := websocket.FormatCloseMessage(websocket.CloseNormalClosure, fmt.Sprintf("%v", readMessageErr))
				var closeErr *websocket.CloseError
				if errors.As(readMessageErr, &closeErr) {
					if closeErr.Code != websocket.CloseNoStatusReceived {
						m = websocket.FormatCloseMessage(closeErr.Code, closeErr.Text)
					}
				}
				connFrontend.WriteMessage(websocket.CloseMessage, m)
				break
			}

			writeMessageErr := connBackend.WriteMessage(wp.IncomeInterceptor(msgType, msg))
			if writeMessageErr != nil {
				wp.Logger.Errorf("writeMessage failed,writeMessageErr:%v", writeMessageErr.Error())
				break
			} else {
				wp.Logger.Errorf("writeMessage successfully")
			}
		}
	}
	processOutComeMsg := func(wg *sync.WaitGroup, connFrontend, connBackend *websocket.Conn) {
		wg.Done()

		for {
			msgType, msg, readMessageErr := connBackend.ReadMessage()
			if readMessageErr != nil {
				wp.Logger.Errorf("readMessage failed,writeMessageErr:%v", readMessageErr.Error())
				m := websocket.FormatCloseMessage(websocket.CloseNormalClosure, fmt.Sprintf("%v", readMessageErr))
				var closeErr *websocket.CloseError
				if errors.As(readMessageErr, &closeErr) {
					if closeErr.Code != websocket.CloseNoStatusReceived {
						m = websocket.FormatCloseMessage(closeErr.Code, closeErr.Text)
					}
				}
				connFrontend.WriteMessage(websocket.CloseMessage, m)
				break
			} else {
				wp.Logger.Infof("readMessage successfully")
			}

			writeMessageErr := connFrontend.WriteMessage(wp.OutcomeInterceptor(msgType, msg))
			if writeMessageErr != nil {
				wp.Logger.Errorf("writeMessage failed,writeMessageErr:%v", writeMessageErr.Error())
				break
			} else {
				wp.Logger.Infof("writeMessage successfully")
			}
		}
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go processInComeMsg(&wg, connFrontend, connBackend)
	go processOutComeMsg(&wg, connFrontend, connBackend)
	wg.Wait()
}

func copyResponse(r http.ResponseWriter, resp *http.Response) error {
	copyHeader := func(dst, src http.Header) {
		for k, vv := range src {
			for _, v := range vv {
				dst.Add(k, v)
			}
		}
	}

	copyHeader(r.Header(), resp.Header)
	r.WriteHeader(resp.StatusCode)
	defer resp.Body.Close()

	_, err := io.Copy(r, resp.Body)
	return err
}
