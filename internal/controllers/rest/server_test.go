package rest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Ragnar-BY/event-collector/internal/domain"
	mock_rest "github.com/Ragnar-BY/event-collector/internal/mocks/usecase"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func getTestGinContext(w *httptest.ResponseRecorder) *gin.Context {
	gin.SetMode(gin.TestMode)

	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
	}

	return ctx
}

func MockJsonPost(c *gin.Context, content interface{}, clientIP string) {
	c.Request.Method = http.MethodPost
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Header.Set("X-Forwarded-For", clientIP)

	c.Request.RemoteAddr = clientIP + ":42123"

	c.Request.Body = io.NopCloser(bytes.NewBuffer(content.([]byte)))
}

func TestSaveEvents(t *testing.T) {

	clientIP := "11.22.33.44"
	serverTime := time.Now()
	logger, _ := zap.NewDevelopment()
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	eventsUsecase := mock_rest.NewMockEventUsecase(ctl)

	w := httptest.NewRecorder()
	ctx := getTestGinContext(w)

	events := []domain.Event{}
	eventsJson := []byte(`[
	{"ClientTime":"2020-12-01T23:59:00Z","DeviceID":"0287d9aa-4adf-4b37-a60f-3e9e645c821e","DeviceOS":"iOS 13.5.1","Session":"ybuRi8mAUypxjbxQ","Sequence":1,"Event":"app_start","ParamInt":0,"ParamStr":"some text","ClientIP":"","ServerTime":"0001-01-01T00:00:00Z"},
	{"ClientTime":"2020-12-01T23:59:00Z","DeviceID":"0287d9aa-4adf-4b37-a60f-3e9e645c821e","DeviceOS":"iOS 13.5.1","Session":"ybuRi8mAUypxjbxQ","Sequence":2,"Event":"app_start","ParamInt":0,"ParamStr":"some text","ClientIP":"","ServerTime":"0001-01-01T00:00:00Z"}
	]`)
	err := json.Unmarshal(eventsJson, &events)
	if err != nil {
		t.Error(err)
		return
	}

	gomock.InOrder(
		eventsUsecase.EXPECT().SaveEvents(ctx.Request.Context(), events, clientIP, serverTime).Return(nil),
	)

	content := []byte(`{	"client_time":"2020-12-01 23:59:00",	"device_id":"0287D9AA-4ADF-4B37-A60F-3E9E645C821E", "device_os":"iOS 13.5.1", "session":"ybuRi8mAUypxjbxQ", "sequence":1, 			"event":"app_start", 			"param_int":0, "param_str":"some text"}
		{	"client_time":"2020-12-01 23:59:00",	"device_id":"0287D9AA-4ADF-4B37-A60F-3E9E645C821E", "device_os":"iOS 13.5.1", "session":"ybuRi8mAUypxjbxQ", "sequence":2, 			"event":"app_start", 			"param_int":0, "param_str":"some text"}`)
	MockJsonPost(ctx, content, clientIP)

	srv := NewServer("", logger, eventsUsecase)
	srv.Now = func() time.Time {
		return serverTime
	}
	srv.SendEvents(ctx)
	time.Sleep(10 * time.Millisecond)

	assert.EqualValues(t, http.StatusOK, w.Code)

}
