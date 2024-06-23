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

	"github.com/gennadyterekhov/gophermart/internal/tests/suites/with_server"

	"github.com/gennadyterekhov/gophermart/internal/httpui/handlers/controllers"

	"github.com/gennadyterekhov/gophermart/internal/httpui/handlers/router"

	"github.com/gennadyterekhov/gophermart/internal/config"
	"github.com/gennadyterekhov/gophermart/internal/fork"
	"github.com/gennadyterekhov/gophermart/internal/logger"
	"github.com/stretchr/testify/suite"

	"github.com/gennadyterekhov/gophermart/internal/tests"
	"github.com/stretchr/testify/require"
)

type testSuite struct {
	with_server.BaseSuiteWithServer
	serverAddress string
	serverPort    string
	serverProcess *fork.BackgroundProcess
	serverArgs    []string
	envs          []string
	cancelContext context.CancelFunc
}

func (suite *testSuite) SetupSuite() {
	conf := config.NewConfig()
	with_server.InitBaseSuiteWithServer(suite)
	controllersStruct := controllers.NewControllers(conf, suite.GetRepository())
	s := tests.NewTestHTTPServer(router.NewRouter(controllersStruct).Router)
	suite.Server = s.Server
	suite.serverAddress = conf.AccrualURL
	suite.serverPort = "8080"
	suite.serverProcess = nil
	suite.serverArgs = []string{""}
	suite.envs = []string{""}

	ctx, cancelContext := context.WithTimeout(context.Background(), 30*time.Second)
	suite.cancelContext = cancelContext
	suite.serverUp(ctx, suite.envs, suite.serverArgs, suite.serverPort)
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

func (suite *testSuite) Test202IfUploadedFirstTime() {
	run := suite.UsingTransactions()

	suite.T().Run("", run(func(t *testing.T) {
		var _ error
		regDto := suite.RegisterForTest("a", "a")

		responseStatusCode := suite.SendPost(
			"/api/user/orders",
			"text/plain",
			regDto.Token,
			bytes.NewBuffer([]byte("12345678903")),
		)

		require.Equal(t, http.StatusAccepted, responseStatusCode)
	}))
}
