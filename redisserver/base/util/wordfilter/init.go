// Package wordfilter .

package wordfilter

import (
	"github.com/sirupsen/logrus"
)

// Default Wordfilter
var Default *Filter

// Init 初始化 Default Wordfilter
func Init(path string, ignoreNoise bool, ignoreWw bool) {
	if Default == nil {
		Default = New()
		if err := Default.LoadWordDict(path, ignoreNoise, ignoreWw); err != nil {
			logrus.Fatalf("worldfilter init error: %v", err)
		}
	}
}
