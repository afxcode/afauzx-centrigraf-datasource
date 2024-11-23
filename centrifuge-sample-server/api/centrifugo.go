package api

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/centrifugal/centrifuge"
)

type CentrifugoAPI struct {
	node         *centrifuge.Node
	allowOrigins []string
	userIDStore  *sync.Map
}

type clientInfo struct {
	UserID     string `json:"user_id"`
	AppVersion string `json:"app_version"`
}

func New(mux *http.ServeMux, allowOrigins []string) (app CentrifugoAPI, err error) {
	node, err := centrifuge.New(centrifuge.Config{
		LogHandler: loggerMiddleware,
		LogLevel:   centrifuge.LogLevelDebug,
	})
	if err != nil {
		return
	}

	app.node = node
	app.userIDStore = &sync.Map{}

	node.OnConnecting(func(ctx context.Context, event centrifuge.ConnectEvent) (centrifuge.ConnectReply, error) {
		userID := app.newUserID()
		return centrifuge.ConnectReply{Credentials: &centrifuge.Credentials{
			UserID: userID,
		}}, nil
	})

	node.OnConnect(func(client *centrifuge.Client) {
		app.userIDStore.Store(client.UserID(), clientInfo{UserID: client.UserID()})

		// TODO: check client permissions
		client.OnSubscribe(func(e centrifuge.SubscribeEvent, cb centrifuge.SubscribeCallback) {
			cb(centrifuge.SubscribeReply{}, nil)
		})

		// TODO: publish validation permission
		client.OnPublish(func(e centrifuge.PublishEvent, cb centrifuge.PublishCallback) {
			app.handleOnPublish(e)
			cb(centrifuge.PublishReply{}, nil)
		})

		// Set Disconnect handler to react on client disconnect events.
		client.OnDisconnect(func(e centrifuge.DisconnectEvent) {
			app.userIDStore.Delete(client.UserID())
		})

		data, _ := json.Marshal(clientInfo{UserID: client.UserID()})
		client.Send(data)

		client.OnMessage(func(event centrifuge.MessageEvent) {
			var appInfo struct {
				Version string `json:"version"`
			}

			json.Unmarshal(event.Data, &appInfo)
			userID := client.UserID()
			if res, ok := app.userIDStore.Load(userID); ok {
				info := res.(clientInfo)
				info.AppVersion = appInfo.Version
				app.userIDStore.Store(userID, info)
			}
		})
	})

	// Run node. This method does not block
	if err = node.Run(); err != nil {
		return
	}

	// Serve Websocket connections using WebsocketHandler.
	for _, o := range allowOrigins {
		origin, e := url.Parse(o)
		if e != nil {
			continue
		}
		app.allowOrigins = append(app.allowOrigins, origin.Host)
	}

	wsHandler := centrifuge.NewWebsocketHandler(node, centrifuge.WebsocketConfig{
		CheckOrigin: app.hostOriginCheck,
	})

	mux.Handle("/api-ws", auth(wsHandler))
	return
}

func (c *CentrifugoAPI) Shutdown(ctx context.Context) error {
	return c.node.Shutdown(ctx)
}

func (c *CentrifugoAPI) hostOriginCheck(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return false
	}
	u, err := url.Parse(origin)
	if err != nil {
		return false
	}

	if len(c.allowOrigins) <= 0 {
		return true
	}

	for _, o := range c.allowOrigins {
		if strings.EqualFold(o, u.Host) {
			return true
		}
	}
	return false
}

func (c *CentrifugoAPI) newUserID() string {
	userID := generateUserID()
	for i := 0; i < 5; i++ {
		if _, ok := c.userIDStore.Load(userID); ok {
			break
		}
		userID = generateUserID()
	}
	return userID
}

func (c *CentrifugoAPI) GetUserList() (clientInfos []clientInfo) {
	c.userIDStore.Range(func(key, value any) bool {
		clientInfos = append(clientInfos, value.(clientInfo))
		return true
	})
	return
}

func generateUserID() string {
	b := make([]byte, 3)
	rand.Read(b)
	return strings.ToUpper(hex.EncodeToString(b))
}
