/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package dock

import (
	"encoding/base64"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"pkg.deepin.io/dde/daemon/appinfo"
	"pkg.deepin.io/lib/xdg/basedir"
	"strings"
	"time"
)

var xdgAutostartDirs []string

func init() {
	configDirs := make([]string, 0, 3)
	configDirs = append(configDirs, basedir.GetUserConfigDir())
	sysConfigDirs := basedir.GetSystemConfigDirs()
	configDirs = append(configDirs, sysConfigDirs...)

	for idx, configDir := range configDirs {
		configDirs[idx] = filepath.Join(configDir, "autostart")
	}
	xdgAutostartDirs = configDirs
}

func isInAutostartDir(file string) bool {
	dir := filepath.Dir(file)
	for _, adir := range xdgAutostartDirs {
		if adir == dir {
			return true
		}
	}
	return false
}

func dataUriToFile(dataUri, path string) (string, error) {
	// dataUri starts with string "data:image/png;base64,"
	commaIndex := strings.Index(dataUri, ",")
	img, err := base64.StdEncoding.DecodeString(dataUri[commaIndex+1:])
	if err != nil {
		return path, err
	}

	return path, ioutil.WriteFile(path, img, 0644)
}

func strSliceEqual(sa, sb []string) bool {
	if len(sa) != len(sb) {
		return false
	}
	for i, va := range sa {
		vb := sb[i]
		if va != vb {
			return false
		}
	}
	return true
}

func uniqStrSlice(slice []string) []string {
	newSlice := make([]string, 0)
	for _, e := range slice {
		if !strSliceContains(newSlice, e) {
			newSlice = append(newSlice, e)
		}
	}
	return newSlice
}

func strSliceContains(slice []string, v string) bool {
	for _, e := range slice {
		if e == v {
			return true
		}
	}
	return false
}

func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

func getCurrentTimestamp() uint32 {
	return uint32(time.Now().Unix())
}

func recordFrequency(appId string) {
	f, err := appinfo.GetFrequencyRecordFile()
	if err == nil {
		appinfo.SetFrequency(appId, appinfo.GetFrequency(appId, f)+1, f) // FIXME: DesktopID???
		f.Free()
	}
}
