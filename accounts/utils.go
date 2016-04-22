/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package accounts

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"pkg.deepin.io/lib/graphic"
	"sort"
	"strings"
)

const (
	polkitManagerUser    = "com.deepin.daemon.accounts.user-administration"
	polkitChangeOwnData  = "com.deepin.daemon.accounts.change-own-user-data"
	polkitSetLoginOption = "com.deepin.daemon.accounts.set-login-option"
)

type ErrCodeType int32

const (
	// 未知错误
	ErrCodeUnkown ErrCodeType = iota
	// 权限认证失败
	ErrCodeAuthFailed
	// 执行命令失败
	ErrCodeExecFailed
	// 传入的参数不合法
	ErrCodeParamInvalid
)

func (code ErrCodeType) String() string {
	switch code {
	case ErrCodeUnkown:
		return "Unkown error"
	case ErrCodeAuthFailed:
		return "Policykit authentication failed"
	case ErrCodeExecFailed:
		return "Exec command failed"
	case ErrCodeParamInvalid:
		return "Invalid parameters"
	}

	return "Unkown error"
}

func clearUserDatas(name string) {
	icons := getUserCustomIcons(name)
	config := path.Join(userConfigDir, name)
	os.Remove(config)

	icons = append(icons, config)
	for _, v := range icons {
		os.Remove(v)
	}
}

func getUserStandardIcons() []string {
	imgs, err := graphic.GetImagesInDir(userIconsDir)
	if err != nil {
		return nil
	}

	var icons []string
	for _, img := range imgs {
		if strings.Contains(img, "guest") ||
			strings.Contains(img, defaultUserIcon) {
			continue
		}

		icons = append(icons, img)
	}

	return icons
}

func getUserCustomIcons(name string) []string {
	return getUserIconsFromDir(userCustomIconsDir, name+"-")
}

func getUserIconsFromDir(dir, condition string) []string {
	imgs, err := graphic.GetImagesInDir(dir)
	if err != nil {
		return nil
	}

	var icons []string
	for _, img := range imgs {
		if !strings.Contains(img, condition) {
			continue
		}

		icons = append(icons, img)
	}

	return icons
}

func isStrInArray(str string, array []string) bool {
	for _, v := range array {
		if v == str {
			return true
		}
	}

	return false
}

func isStrvEqual(l1, l2 []string) bool {
	if len(l1) != len(l2) {
		return false
	}

	sort.Strings(l1)
	sort.Strings(l2)
	for i, v := range l1 {
		if v != l2[i] {
			return false
		}
	}
	return true
}

func polkitAuthManagerUser(pid uint32) error {
	return polkitAuthentication(polkitManagerUser, pid)
}

func polkitAuthChangeOwnData(pid uint32) error {
	return polkitAuthentication(polkitChangeOwnData, pid)
}

func polkitAuthentication(action string, pid uint32) error {
	success, err := polkitAuthWithPid(action, pid)
	if err != nil {
		return err
	}

	if !success {
		return fmt.Errorf(ErrCodeAuthFailed.String())
	}

	return nil
}

const (
	pidFileStatus = "/proc/%v/status"
)

func getUidByPid(pid uint32) (string, error) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("Recover error in getUidByPid:", err)
		}
	}()

	var file = fmt.Sprintf(pidFileStatus, pid)
	datas, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	var lines = strings.Split(string(datas), "\n")
	for _, line := range lines {
		if !strings.Contains(line, "Uid:") {
			continue
		}

		strv := strings.Split(line, "\t")
		return strv[1], nil
	}

	return "", fmt.Errorf("Invalid file: %s", file)
}