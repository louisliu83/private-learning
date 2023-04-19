package log

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"pa.cn/fedlearn/psi/api"
)

func Debugln(ctx context.Context, args ...interface{}) {
	logrus.WithContext(ctx).
		WithField(api.Trace_ID, fmt.Sprintf("%s", ctx.Value(api.Trace_ID))).
		Debugln(args...)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	logrus.WithContext(ctx).
		WithField(api.Trace_ID, fmt.Sprintf("%s", ctx.Value(api.Trace_ID))).
		Debugf(format, args...)
}

func Infoln(ctx context.Context, args ...interface{}) {
	logrus.WithContext(ctx).
		WithField(api.Trace_ID, fmt.Sprintf("%s", ctx.Value(api.Trace_ID))).
		Infoln(args...)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	logrus.WithContext(ctx).
		WithField(api.Trace_ID, fmt.Sprintf("%s", ctx.Value(api.Trace_ID))).
		Infof(format, args...)
}

func Warningln(ctx context.Context, args ...interface{}) {
	logrus.WithContext(ctx).
		WithField(api.Trace_ID, fmt.Sprintf("%s", ctx.Value(api.Trace_ID))).
		Warningln(args...)
}

func Warningf(ctx context.Context, format string, args ...interface{}) {
	logrus.WithContext(ctx).
		WithField(api.Trace_ID, fmt.Sprintf("%s", ctx.Value(api.Trace_ID))).
		Warningf(format, args...)
}

func Warnln(ctx context.Context, args ...interface{}) {
	logrus.WithContext(ctx).
		WithField(api.Trace_ID, fmt.Sprintf("%s", ctx.Value(api.Trace_ID))).
		Warnln(args...)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	logrus.WithContext(ctx).
		WithField(api.Trace_ID, fmt.Sprintf("%s", ctx.Value(api.Trace_ID))).
		Warnf(format, args...)
}

func Errorln(ctx context.Context, args ...interface{}) {
	logrus.WithContext(ctx).
		WithField(api.Trace_ID, fmt.Sprintf("%s", ctx.Value(api.Trace_ID))).
		Errorln(args...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	logrus.WithContext(ctx).
		WithField(api.Trace_ID, fmt.Sprintf("%s", ctx.Value(api.Trace_ID))).
		Errorf(format, args...)
}

func Fatalln(ctx context.Context, args ...interface{}) {
	logrus.WithContext(ctx).
		WithField(api.Trace_ID, fmt.Sprintf("%s", ctx.Value(api.Trace_ID))).
		Fatalln(args...)
}

func Fatalf(ctx context.Context, format string, args ...interface{}) {
	logrus.WithContext(ctx).
		WithField(api.Trace_ID, fmt.Sprintf("%s", ctx.Value(api.Trace_ID))).
		Fatalf(format, args...)
}
