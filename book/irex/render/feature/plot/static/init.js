// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
import * as plot from "/static/plot.js";
const data = new Map([
{{- range $k, $v := .Quotes}}[{{$k}},{d:{{$v.UnixTime}}, c:{{$v.Cents}}}],{{end -}}
]);
plot.Init("plot-svg", data);
