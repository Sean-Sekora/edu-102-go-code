package main

import (
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	translation "temporal102/exercises/durable-execution/practice"
)

func main() {
	c, err := client.Dial(client.Options{
		Logger: translation.NewZapAdapter(
			NewZapLogger()),
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, translation.TaskQueueName, worker.Options{})

	w.RegisterWorkflow(translation.SayHelloGoodbye)
	w.RegisterActivity(translation.TranslateTerm)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}

func NewZapLogger() *zap.Logger {
	encodeConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:      false,
		Sampling:         nil, // consider exposing this to config for our external customer
		Encoding:         "console",
		EncoderConfig:    encodeConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	logger, err := config.Build()

	// Or simple
	// logger, err := zap.NewDevelopment()
	// can be used instead of the code above.

	if err != nil {
		log.Fatalln("Unable to create zap logger")
	}
	return logger
}
