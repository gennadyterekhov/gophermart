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

	"github.com/gennadyterekhov/gophermart/internal/luhn"

	"github.com/gennadyterekhov/gophermart/internal/domain/models/order"

	"github.com/stretchr/testify/require"

	"github.com/gennadyterekhov/gophermart/internal/logger"

	"github.com/gennadyterekhov/gophermart/internal/tests/helpers"

	"github.com/gennadyterekhov/gophermart/internal/domain/responses"
	"github.com/gennadyterekhov/gophermart/internal/repositories"
	"github.com/stretchr/testify/assert"

	"github.com/go-resty/resty/v2"

	"github.com/gennadyterekhov/gophermart/internal/config"

	"github.com/gennadyterekhov/gophermart/internal/fork"

	"github.com/gennadyterekhov/gophermart/internal/tests"
	"github.com/stretchr/testify/suite"
)

type AccrualClientSuite struct {
	suite.Suite

	serverAddress string
	serverPort    string
	serverProcess *fork.BackgroundProcess
	serverArgs    []string
	envs          []string
}

var testSuite *AccrualClientSuite

func (suite *AccrualClientSuite) SetupSuite() {
}

func (suite *AccrualClientSuite) serverUp(ctx context.Context, envs, args []string, port string) {
	serverBinaryPath := "/Users/gena/code/yandex/practicum/golang_advanced/gophermart/cmd/accrual/accrual_darwin_arm64"
	suite.serverProcess = fork.NewBackgroundProcess(
		context.Background(),
		serverBinaryPath,
		fork.WithEnv(envs...),
		fork.WithArgs(args...),
	)

	err := suite.serverProcess.Start(ctx)
	if err != nil {
		logger.ZapSugarLogger.Debugln(err.Error())
		suite.T().Errorf("Невозможно запустить процесс командой %q: %s. Переменные окружения: %+v, флаги командной строки: %+v", suite.serverProcess, err, envs, args)
		return
	}

	err = suite.serverProcess.WaitPort(ctx, "tcp", port)
	if err != nil {
		logger.ZapSugarLogger.Debugln(err.Error()) // context deadline exceeded
		suite.T().Errorf("Не удалось дождаться пока порт %s станет доступен для запроса: %s", port, err)
		return
	}
}

func (suite *AccrualClientSuite) TearDownSuite() {
	suite.serverShutdown()
}

func (suite *AccrualClientSuite) serverShutdown() {
	if suite.serverProcess == nil {
		return
	}

	exitCode, err := suite.serverProcess.Stop(syscall.SIGINT, syscall.SIGKILL)
	if err != nil {
		if errors.Is(err, os.ErrProcessDone) {
			return
		}
		logger.ZapSugarLogger.Debugf("Не удалось остановить процесс с помощью сигнала ОС: %s", err)
		return
	}

	if exitCode > 0 {
		logger.ZapSugarLogger.Debugf("Процесс завершился с не нулевым статусом %d", exitCode)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	out := suite.serverProcess.Stderr(ctx)
	if len(out) > 0 {
		logger.ZapSugarLogger.Debugf("Получен STDERR лог процесса:\n\n%s", string(out))
	}
	out = suite.serverProcess.Stdout(ctx)
	if len(out) > 0 {
		logger.ZapSugarLogger.Debugf("Получен STDOUT лог процесса:\n\n%s", string(out))
	}
}

func TestMain(m *testing.M) {
	tests.BeforeAll()
	testSuite = &AccrualClientSuite{
		serverAddress: config.ServerConfig.AccrualURL,
		serverPort:    "8080",
		serverProcess: nil,
		serverArgs:    []string{""},
		envs:          []string{""},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	testSuite.serverUp(ctx, testSuite.envs, testSuite.serverArgs, testSuite.serverPort)
	code := m.Run()
	cancel()
	testSuite.TearDownSuite()
	tests.AfterAll()
	os.Exit(code)
}

func TestCanGetOrderStatus(t *testing.T) {
	testSuite.SetT(t)

	run := tests.UsingTransactions()
	t.Run("", run(func(t *testing.T) {
		userDto := helpers.RegisterForTest("a", "a")
		number := luhn.Generate(1)
		createOrder(t, userDto, number)

		httpc := resty.NewWithClient(&http.Client{
			Transport: &http.Transport{
				DisableCompression: true,
			},
		}).SetHostURL(testSuite.serverAddress)

		registerOrderInAccrual(t, httpc, number)
		time.Sleep(time.Millisecond * 20)

		req := httpc.R()
		resp, err := req.Get(fmt.Sprintf("/api/orders/%v", number))
		assert.NoError(t, err)
		logger.ZapSugarLogger.Debugln("resp.StatusCode()", resp.StatusCode())
		assert.Equal(t, http.StatusOK, resp.StatusCode())
	}))
}

func TestTooManyRequests(t *testing.T) {
	t.Skipf("cannot test")
}

func TestNoContent(t *testing.T) {
	t.Skipf("cannot test in suite because it can run after registration and will fail")

	testSuite.SetT(t)

	run := tests.UsingTransactions()
	t.Run("", run(func(t *testing.T) {
		userDto := helpers.RegisterForTest("a", "a")
		number := luhn.Generate(1)
		createOrder(t, userDto, number)

		httpc := resty.NewWithClient(&http.Client{
			Transport: &http.Transport{
				DisableCompression: true,
			},
		}).SetHostURL(testSuite.serverAddress)

		req := httpc.R()
		resp, err := req.Get(fmt.Sprintf("/api/orders/%v", number))
		assert.NoError(t, err)
		logger.ZapSugarLogger.Debugln("resp.StatusCode()", resp.StatusCode())
		assert.Equal(t, http.StatusNoContent, resp.StatusCode())
	}))
}

func TestInternalServerError(t *testing.T) {
	t.Skipf("cannot test")
}

func createOrder(
	t *testing.T,
	userDto *responses.Register,
	number string,
) *order.Order {
	var ten int64 = 10
	orderNewest, err := repositories.AddOrder(
		context.Background(),
		number,
		userDto.ID,
		order.Registered,
		&ten,
		time.Time{},
	)
	assert.NoError(t, err)

	return orderNewest
}

func registerOrderInAccrual(
	t *testing.T,
	httpc *resty.Client,
	number string,
) {
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
	require.NoError(t, err)
	logger.ZapSugarLogger.Debugln(resp.StatusCode())
	logger.ZapSugarLogger.Debugln(string(resp.Body()))
}
