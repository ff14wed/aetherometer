/*
Original work Copyright (c) 2018 Adam Scarr
from https://github.com/99designs/gqlgen/blob/56f3f92b8ee315ef8a3c32484e6b63dd13ae574a/handler/playground.go

Modified work Copyright (c) 2019 ff14wed

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package handlers

import (
	"html/template"
	"net/http"
)

var page = template.Must(template.New("graphiql").Parse(`<!DOCTYPE html>
<html>
<head>
	<meta charset=utf-8/>
	<meta name="viewport" content="user-scalable=no, initial-scale=1.0, minimum-scale=1.0, maximum-scale=1.0, minimal-ui">
	<link rel="shortcut icon" href="https://graphcool-playground.netlify.com/favicon.png">
	<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/graphql-playground-react@{{ .version }}/build/static/css/index.css"
		integrity="{{ .cssSRI }}" crossorigin="anonymous"/>
	<link rel="shortcut icon" href="https://cdn.jsdelivr.net/npm/graphql-playground-react@{{ .version }}/build/favicon.png"
		integrity="{{ .faviconSRI }}" crossorigin="anonymous"/>
	<script src="https://cdn.jsdelivr.net/npm/graphql-playground-react@{{ .version }}/build/static/js/middleware.js"
		integrity="{{ .jsSRI }}" crossorigin="anonymous"></script>
	<title>{{.title}}</title>
</head>
<body>
<style type="text/css">
	html { font-family: "Open Sans", sans-serif; overflow: hidden; }
	body { margin: 0; background: #172a3a; }
</style>
<div id="root"/>
<script type="text/javascript">
	const getPluginParams = async () => {
		if (!window.waitForInit) {
			return null
		} else {
			return new Promise((resolve) => {
				window.initPlugin = (params) => {
					resolve(params);
				};
			});
		}
	};
	window.addEventListener('load', async function (event) {
		const params = await getPluginParams();
		const headers = {};
		if (params) {
			headers.Authorization = params.apiToken;
		}
		const root = document.getElementById('root');
		root.classList.add('playgroundIn');
		const wsProto = location.protocol == 'https:' ? 'wss:' : 'ws:'
		GraphQLPlayground.init(root, {
			endpoint: location.protocol + '//' + location.host + '{{.endpoint}}',
			subscriptionsEndpoint: wsProto + '//' + location.host + '{{.endpoint }}',
			settings: {
				'request.credentials': 'same-origin'
			},
			headers: headers,
		})
	})
</script>
</body>
</html>
`))

func Playground(title string, endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		err := page.Execute(w, map[string]string{
			"title":      title,
			"endpoint":   endpoint,
			"version":    "1.7.20",
			"cssSRI":     "sha256-cS9Vc2OBt9eUf4sykRWukeFYaInL29+myBmFDSa7F/U=",
			"faviconSRI": "sha256-GhTyE+McTU79R4+pRO6ih+4TfsTOrpPwD8ReKFzb3PM=",
			"jsSRI":      "sha256-4QG1Uza2GgGdlBL3RCBCGtGeZB6bDbsw8OltCMGeJsA=",
		})
		if err != nil {
			panic(err)
		}
	}
}
