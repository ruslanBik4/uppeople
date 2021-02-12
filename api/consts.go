// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

// names of system environment variables
const (
	CFG_PATH    = "configPath"
	WEB_PATH    = "webPath"
	SYSTEM_PATH = "systemPath"
)

const (
	CODE_LOG_UPDATE = 100
	CODE_LOG_INSERT = 101
)

const (
	EMAIL_TEXT = `<p><span style="font-size: 14px;">Please, review %s %s  CV</span></p>
<p>%s</p>
<p><br>Will be appreciate for quick feedback.</p>
<p><br><br></p>
<p>@"UPpeople" Recruiting agency</p>
<p>&nbsp;<a href="http://www.rock-it.com.ua/" target="_self"><span style="color: blue;font-size: 16px;font-family: Journal, serif;">http://www.rock-it.com.ua/</span></a><span style="font-size: 16px;"> </span></p>`
)
