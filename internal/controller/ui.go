package controller

import (
	"io"
	"net/http"
	"os"

	"go.uber.org/zap"
)

const mainUI = "./../../ui/mainui.html"

type UI interface {
	GetMainUi(w http.ResponseWriter, r *http.Request)
}

type UIController struct {
	logger *zap.Logger
}

func NewUIController(logger *zap.Logger) UI {

	return &UIController{
		logger: logger,
	}
}

func (u *UIController) GetMainUi(w http.ResponseWriter, r *http.Request) {
	data, err := u.loadFile(mainUI)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (u *UIController) loadFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		u.logger.Error(err.Error())
		return nil, err
	}

	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		u.logger.Error(err.Error())
		return nil, err
	}

	data := make([]byte, fileInfo.Size())

	for {
		_, err := file.Read(data)
		if err == io.EOF {
			break
		}
	}

	return data, nil
}
