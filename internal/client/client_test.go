package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/gennadyterekhov/gophermart/internal/tests/suites/base"

	"github.com/gennadyterekhov/gophermart/internal/luhn"

	"github.com/gennadyterekhov/gophermart/internal/domain/models/order"

	"github.com/stretchr/testify/require"

	"github.com/gennadyterekhov/gophermart/internal/logger"

	"github.com/gennadyterekhov/gophermart/internal/domain/responses"
	"github.com/stretchr/testify/assert"

	"github.com/go-resty/resty/v2"

	"github.com/gennadyterekhov/gophermart/internal/fork"

	"github.com/stretchr/testify/suite"
)

type testSuite struct {
	base.BaseSuite

	serverAddress string
	serverPort    string
	serverProcess *fork.BackgroundProcess
	serverArgs    []string
	envs          []string
	cancelContext context.CancelFunc
}

func (suite *testSuite) SetupSuite() {
	base.InitBaseSuite(suite)

	suite.serverAddress = "http://localhost"
	suite.serverPort = "8089"
	suite.serverProcess = nil
	suite.serverArgs = []string{""}
	suite.envs = []string{""}

	ctx, cancelContext := context.WithTimeout(context.Background(), 30*time.Second)
	suite.serverUp(ctx, suite.envs, suite.serverArgs, suite.serverPort)
	suite.cancelContext = cancelContext
}

func (suite *testSuite) serverUp(ctx context.Context, envs, args []string, port string) {
	serverBinaryPath := "/Users/gena/code/yandex/practicum/golang_advanced/gophermart/cmd/accrual/accrual_darwin_arm64"
	suite.serverProcess = fork.NewBackgroundProcess(
		context.Background(),
		serverBinaryPath,
		fork.WithEnv(envs...),
		fork.WithArgs(args...),
	)

	err := suite.serverProcess.Start(ctx)
	if err != nil {
		logger.CustomLogger.Debugln(err.Error())
		suite.T().Errorf("Невозможно запустить процесс командой %q: %s. Переменные окружения: %+v, флаги командной строки: %+v", suite.serverProcess, err, envs, args)
		return
	}

	err = suite.serverProcess.WaitPort(ctx, "tcp", port)
	if err != nil {
		logger.CustomLogger.Debugln(err.Error()) // context deadline exceeded
		suite.T().Errorf("Не удалось дождаться пока порт %s станет доступен для запроса: %s", port, err)
		return
	}
}

func (suite *testSuite) TearDownSuite() {
	suite.cancelContext()
	suite.serverShutdown()
}

func (suite *testSuite) serverShutdown() {
	if suite.serverProcess == nil {
		return
	}

	exitCode, err := suite.serverProcess.Stop(syscall.SIGINT, syscall.SIGKILL)
	if err != nil {
		if errors.Is(err, os.ErrProcessDone) {
			return
		}
		logger.CustomLogger.Debugln("Не удалось остановить процесс с помощью сигнала ОС: %s", err.Error())
		return
	}

	if exitCode > 0 {
		logger.CustomLogger.Debugln("Процесс завершился с не нулевым статусом %d", exitCode)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	out := suite.serverProcess.Stderr(ctx)
	if len(out) > 0 {
		logger.CustomLogger.Debugln("Получен STDERR лог процесса:\n\n%s", string(out))
	}
	out = suite.serverProcess.Stdout(ctx)
	if len(out) > 0 {
		logger.CustomLogger.Debugln("Получен STDOUT лог процесса:\n\n%s", string(out))
	}
}

func Test(t *testing.T) {
	suite.Run(t, new(testSuite))
}

func (suite *testSuite) TestCanGetOrderStatus() {
	userDto := suite.RegisterForTest("b", "a")
	number := luhn.Generate(1)
	suite.createOrder(userDto, number)

	httpc := resty.NewWithClient(&http.Client{
		Transport: &http.Transport{
			DisableCompression: true,
		},
	}).SetBaseURL(suite.serverAddress)
	suite.registerOrderInAccrual(httpc, number)
	time.Sleep(time.Millisecond * 20)

	req := httpc.R()
	resp, err := req.Get(fmt.Sprintf("/api/orders/%v", number))
	assert.NoError(suite.T(), err)
	logger.CustomLogger.Debugln("resp.StatusCode()", resp.StatusCode())
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode())
}

func (suite *testSuite) TestTooManyRequests() {
	suite.T().Skipf("cannot test")
}

func (suite *testSuite) TestNoContent() {
	suite.T().Skipf("cannot test in suite because it can run after registration and will fail")

	userDto := suite.RegisterForTest("a", "a")
	number := luhn.Generate(1)
	suite.createOrder(userDto, number)

	httpc := resty.NewWithClient(&http.Client{
		Transport: &http.Transport{
			DisableCompression: true,
		},
	}).SetBaseURL(suite.serverAddress)

	req := httpc.R()
	resp, err := req.Get(fmt.Sprintf("/api/orders/%v", number))
	assert.NoError(suite.T(), err)
	logger.CustomLogger.Debugln("resp.StatusCode()", resp.StatusCode())
	assert.Equal(suite.T(), http.StatusNoContent, resp.StatusCode())
}

func (suite *testSuite) TestInternalServerError() {
	suite.T().Skipf("cannot test")
}

func (suite *testSuite) createOrder(userDto *responses.Register, number string) *order.Order {
	var ten int64 = 10
	orderNewest, err := suite.Repository.AddOrder(
		context.Background(),
		number,
		userDto.ID,
		order.Registered,
		&ten,
		time.Time{},
	)
	assert.NoError(suite.T(), err)

	return orderNewest
}

func (suite *testSuite) registerOrderInAccrual(httpc *resty.Client, number string) {
	o := []byte(`
			{
				"order": "` + number + `",
				"goods": [
					{
						"description": "Стиральная машинка LG",
						"price": 47399.99
					}
				]
			}
		`)

	req := httpc.R().
		SetHeader("Content-Type", "application/json").
		SetBody(o)

	resp, err := req.Post("/api/orders")
	require.NoError(suite.T(), err)
	logger.CustomLogger.Debugln(resp.StatusCode())
	logger.CustomLogger.Debugln(string(resp.Body()))
}
