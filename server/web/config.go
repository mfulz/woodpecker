// Copyright 2021 Woodpecker Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/woodpecker-ci/woodpecker/server"
	"github.com/woodpecker-ci/woodpecker/server/router/middleware/session"
	"github.com/woodpecker-ci/woodpecker/shared/token"
	"github.com/woodpecker-ci/woodpecker/version"
)

func WebConfig(c *gin.Context) {
	var err error

	user := session.User(c)

	var csrf string
	if user != nil {
		csrf, _ = token.New(
			token.CsrfToken,
			user.Login,
		).Sign(user.Hash)
	}

	var syncing bool
	if user != nil {
		syncing = time.Unix(user.Synced, 0).Add(time.Hour * 72).Before(time.Now())
	}

	configData := map[string]interface{}{
		"user":    user,
		"csrf":    csrf,
		"syncing": syncing,
		"docs":    server.Config.Server.Docs,
		"version": version.String(),
	}

	// default func map with json parser.
	var funcMap = template.FuncMap{
		"json": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}

	c.Header("Content-Type", "text/javascript; charset=utf-8")
	tmpl := template.Must(template.New("").Funcs(funcMap).Parse(configTemplate))

	err = tmpl.Execute(c.Writer, configData)
	if err != nil {
		fmt.Println(err)
		c.AbortWithError(http.StatusInternalServerError, nil)
	}
}

const configTemplate = `
window.WOODPECKER_USER = {{ json .user }};
window.WOODPECKER_SYNC = {{ .syncing }};
window.WOODPECKER_CSRF = "{{ .csrf }}";
window.WOODPECKER_VERSION = "{{ .version }}";
window.WOODPECKER_DOCS = "{{ .docs }}";
`
