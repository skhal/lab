# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# vim: tabstop=2

PORTNAME=       lab-certval
PORTVERSION=    260929
PORTREVISION=   1
CATEGORIES=     sysutils

MAINTAINER=     sn.khalatyan@gmail.com
COMMENT=        Certificate validator with Slack notifications.
WWW=            https://github.com/skhal/lab

LICENSE=        BSD3CLAUSE
LICENSE_FILE=   ${WRKSRC}/LICENSE

BUILD_DEPENDS=  protoc:devel/protobuf

USES=           go:modules
USE_GITHUB=     yes
GH_ACCOUNT=     skhal
GH_PROJECT=     lab
GH_TAGNAME=     aaa93283
GO_TARGET=      ./infra/cmd/certval:lab-certval
GH_TUPLE=       \
                BurntSushi:toml:a339e1f7089c:burntsushi_toml/vendor/github.com/BurntSushi/toml \
                PuerkitoBio:goquery:v1.12.0:puerkitobio_goquery/vendor/github.com/PuerkitoBio/goquery \
                andybalholm:cascadia:v1.3.3:andybalholm_cascadia/vendor/github.com/andybalholm/cascadia \
                bufbuild:protocompile:v0.14.1:bufbuild_protocompile/vendor/github.com/bufbuild/protocompile \
                chzyer:readline:v1.5.1:chzyer_readline/vendor/github.com/chzyer/readline \
                dominikh:go-tools:v0.8.1:dominikh_go_tools/vendor/honnef.co/go/tools \
                go-yaml:yaml:v3.0.1:go_yaml_yaml/vendor/gopkg.in/yaml.v3 \
                golang:mod:v0.35.0:golang_mod/vendor/golang.org/x/mod \
                golang:net:v0.54.0:golang_net/vendor/golang.org/x/net \
                golang:sync:v0.20.0:golang_sync/vendor/golang.org/x/sync \
                golang:sys:v0.44.0:golang_sys/vendor/golang.org/x/sys \
                golang:text:v0.37.0:golang_text/vendor/golang.org/x/text \
                golang:tools:19499e7caabc:golang_tools/vendor/golang.org/x/tools \
                gomarkdown:markdown:7d523f7318df:gomarkdown_markdown/vendor/github.com/gomarkdown/markdown \
                google:go-cmp:v0.7.0:google_go_cmp/vendor/github.com/google/go-cmp \
                google:go-genproto:aa98bba5eb94:google_go_genproto_rpc/vendor/google.golang.org/genproto/googleapis/rpc \
                google:pprof:545e8a4df936:google_pprof/vendor/github.com/google/pprof \
                grpc:grpc-go:cmd/protoc-gen-go-grpc/v1.6.1:grpc_grpc_go_1/vendor/google.golang.org/grpc/cmd/protoc-gen-go-grpc \
                grpc:grpc-go:v1.81.1:grpc_grpc_go/vendor/google.golang.org/grpc \
                ianlancetaylor:demangle:f615e6bd150b:ianlancetaylor_demangle/vendor/github.com/ianlancetaylor/demangle \
                json-iterator:go:v1.1.12:json_iterator_go/vendor/github.com/json-iterator/go \
                kevinpollet:nego:a65cd48cee43:kevinpollet_nego/vendor/github.com/kevinpollet/nego \
                lpar:gzipped:v2.1.0:lpar_gzipped_v2/vendor/github.com/lpar/gzipped/v2 \
                mitchellh:go-wordwrap:v1.0.1:mitchellh_go_wordwrap/vendor/github.com/mitchellh/go-wordwrap \
                mmcdole:gofeed:v1.3.0:mmcdole_gofeed/vendor/github.com/mmcdole/gofeed \
                mmcdole:goxpp:v1.1.1:mmcdole_goxpp/vendor/github.com/mmcdole/goxpp \
                modern-go:concurrent:bacd9c7ef1dd:modern_go_concurrent/vendor/github.com/modern-go/concurrent \
                modern-go:reflect2:v1.0.2:modern_go_reflect2/vendor/github.com/modern-go/reflect2 \
                protocolbuffers:protobuf-go:v1.36.11:protocolbuffers_protobuf_go/vendor/google.golang.org/protobuf \
                protocolbuffers:txtpbfmt:c39628bde8b5:protocolbuffers_txtpbfmt/vendor/github.com/protocolbuffers/txtpbfmt \
                stretchr:testify:v1.11.1:stretchr_testify/vendor/github.com/stretchr/testify

post-extract:
	@${CP} ${FILESDIR}/modules.txt ${WRKSRC}/vendor

.include <bsd.port.mk>
