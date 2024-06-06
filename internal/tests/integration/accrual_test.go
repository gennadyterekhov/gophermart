package integration

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/gennadyterekhov/gophermart/internal/config"
	"github.com/gennadyterekhov/gophermart/internal/fork"
	"github.com/gennadyterekhov/gophermart/internal/httpui/handlers"
	"github.com/gennadyterekhov/gophermart/internal/logger"
	"github.com/stretchr/testify/suite"

	"github.com/gennadyterekhov/gophermart/internal/tests"
	"github.com/gennadyterekhov/gophermart/internal/tests/helpers"
	"github.com/stretchr/testify/require"
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

func Test202IfUploadedFirstTime(t *testing.T) {
	testSuite.SetT(t)
	run := tests.UsingTransactions()
	tests.InitTestServer(handlers.GetRouter())

	t.Run("", run(func(t *testing.T) {
		var _ error
		regDto := helpers.RegisterForTest("a", "a")

		responseStatusCode := tests.SendPost(
			t,
			tests.TestServer,
			"/api/user/orders",
			"text/plain",
			regDto.Token,
			bytes.NewBuffer([]byte("12345678903")),
		)

		require.Equal(t, http.StatusAccepted, responseStatusCode)
	}))
}
