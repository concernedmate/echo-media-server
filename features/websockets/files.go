package websockets

import (
	"bytes"
	"media-server/configs"
	"media-server/features/models"
	"os"
	"path"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type FileFormat struct {
	FileID    string
	Filename  []byte
	Directory []byte
	Data      *os.File
}

var upgrader = websocket.Upgrader{}
var EOF = append([]byte("TRANSMISSION_END"), []byte{0x00}...)
var SPLIT_END = append([]byte("SPLIT_END"), []byte{0x00}...)

func UploadMultipleFiles(c echo.Context) (err error) {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	username, ok := c.Get("username").(string)
	if !ok {
		return ws.WriteMessage(websocket.TextMessage, []byte("[Error] invalid user token"))
	}

	file := FileFormat{
		FileID:    uuid.NewString(),
		Filename:  []byte{},
		Directory: []byte{},
		Data:      nil,
	}
	transfer_mode := 0
	for {
		var msg []byte

		_, msg, err = ws.ReadMessage()
		if err != nil {
			return ws.WriteMessage(websocket.TextMessage, []byte("[Error] "+err.Error()))
		}

		if bytes.Equal(EOF, msg) {
			break
		}

		switch transfer_mode {
		case 0:
			{
				file.Filename = append(file.Filename, msg...)
				transfer_mode++
			}
		case 1:
			{
				file.Directory = append(file.Directory, msg...)
				transfer_mode++

				file.Data, err = os.Create(path.Join(configs.UPLOAD_BASEDIR(), file.FileID))
				if err != nil {
					return ws.WriteMessage(websocket.TextMessage, []byte("[Error] "+err.Error()))
				}
				defer func() {
					file.Data.Close()
					if err != nil {
						_ = os.Remove(path.Join(configs.UPLOAD_BASEDIR(), file.FileID))
					}
				}()
			}
		case 2:
			{
				if bytes.Equal(SPLIT_END, msg) {
					// TODO add checking before uploading
					err = models.SaveFileWebsocket(file.FileID, string(file.Filename), string(file.Directory), username)
					if err != nil {
						_ = os.Remove(path.Join(configs.UPLOAD_BASEDIR(), file.FileID))
					}

					file = FileFormat{
						FileID:    uuid.NewString(),
						Filename:  []byte{},
						Directory: []byte{},
						Data:      nil,
					}

					transfer_mode = 0
					continue
				}

				_, err = file.Data.Write(msg)
				if err != nil {
					return ws.WriteMessage(websocket.TextMessage, []byte("[Error] "+err.Error()))
				}
			}
		}
	}

	return ws.WriteMessage(websocket.TextMessage, []byte("[Success]"))
}
