---
layout: default
title: gRPC-Gateway
nav_order: 1
description: "Documentation site for the gRPC-Gateway"
permalink: /
---

# gRPC-Gateway
{: .fs-9 }

gRPC-Gateway is a plugin of [protoc](https://github.com/protocolbuffers/protobuf).
It reads a [gRPC](https://grpc.io) service definition,
and generates a reverse-proxy server which translates a RESTful JSON API into gRPC.
This server is generated according to [custom options](https://cloud.google.com/service-infrastructure/docs/service-management/reference/rpc/google.api#http) in your gRPC definition.
{: .fs-6 .fw-300 }

[Get started now](#getting-started){: .btn .btn-primary .fs-5 .mb-4 .mb-md-0 .mr-2 } [View it on GitHub](https://github.com/grpc-ecosystem/grpc-gateway){: .btn .fs-5 .mb-4 .mb-md-0 }

---

## About

[![circleci](https://img.shields.io/circleci/build/github/grpc-ecosystem/grpc-gateway?color=379c9c&logo=circleci&logoColor=ffffff&style=flat-square)](https://circleci.com/gh/grpc-ecosystem/grpc-gateway)
[![codecov](https://img.shields.io/codecov/c/github/grpc-ecosystem/grpc-gateway?color=379c9c&logo=codecov&logoColor=ffffff&style=flat-square)](https://codecov.io/gh/grpc-ecosystem/grpc-gateway)
[![forks](https://img.shields.io/github/forks/grpc-ecosystem/grpc-gateway?color=379c9c&style=flat-square)](https://github.com/grpc-ecosystem/grpc-gateway/network/members)
[![issues](https://img.shields.io/github/issues/grpc-ecosystem/grpc-gateway?color=379c9c&style=flat-square)](https://github.com/grpc-ecosystem/grpc-gateway/issues)
[![license](https://img.shields.io/github/license/grpc-ecosystem/grpc-gateway?color=379c9c&style=flat-square)](https://github.com/grpc-ecosystem/grpc-gateway/blob/master/LICENSE.txt)
[![stars](https://img.shields.io/github/stars/grpc-ecosystem/grpc-gateway?color=379c9c&style=flat-square)](https://github.com/grpc-ecosystem/grpc-gateway/stargazers)

grpc-gateway is a plugin of [protoc](https://github.com/protocolbuffers/protobuf).
It reads a [gRPC](https://grpc.io) service definition,
and generates a reverse-proxy server which translates a RESTful JSON API into gRPC.
This server is generated according to [custom options](https://cloud.google.com/service-infrastructure/docs/service-management/reference/rpc/google.api#http) in your gRPC definition.

It helps you to provide your APIs in both gRPC and RESTful style at the same time.

![architecture introduction diagram](https://docs.google.com/drawings/d/12hp4CPqrNPFhattL_cIoJptFvlAqm5wLQ0ggqI5mkCg/pub?w=749&h=370)

To learn more about us check out our documentation.

## Contribution

See [CONTRIBUTING.md](http://github.com/grpc-ecosystem/grpc-gateway/blob/master/CONTRIBUTING.md).

## License

grpc-gateway is licensed under the BSD 3-Clause License.
See [LICENSE.txt](https://github.com/grpc-ecosystem/grpc-gateway/blob/master/LICENSE.txt) for more details.

#### Thank you to the contributors of gRPC-Gateway!

<ul class="list-style-none">
{% for contributor in site.github.contributors %}
  <li class="d-inline-block mr-1">
     <a href="{{ contributor.html_url }}"><img src="{{ contributor.avatar_url }}" width="32" height="32" alt="{{ contributor.login }}"/></a>
  </li>
{% endfor %}
</ul>
