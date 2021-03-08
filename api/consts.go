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
	CODE_LOG_PEFORM = 102
	CODE_LOG_DELETE = 103
)

const (
	EMAIL_TEXT = `<p><span style="font-size: 14px;">Colleagues,
please review the candidacy of %s for the position of %s </span></p>
<p>CV:%s</p>
<p>Experience:%s</p>
<p>English level:%s</p>
<p>Salary expectations:$%v</p>
<p><br>Will be appreciate for quick feedback.</p>
<p><br><br></p>
<p>Best regards,
UPPeople team.</p>
<p><span style="font-size: 14px;">Добрый день,
рассмотрите, пожалуйста, кандидата %[1]s на позицию  %s </span></p>
<p>CV:%s</p>
<p>Опыт:%s</p>
<p>Уровень английского:%s</p>
<p>Ожидания по заработной плате:$%v</p>
<p><br>Будем благодарны за фидбек.</p>
<p><br><br></p>
<p>С наилучшим пожеланиями,
команда UPPeople.</p>
<p>&nbsp;<a href="http://my.uppeople.co/" target="_self"><span style="color: blue;font-size: 16px;font-family: Journal, serif;"@"UPpeople" Recruiting agency</span></a><span style="font-size: 16px;"> </span></p>
`
)
