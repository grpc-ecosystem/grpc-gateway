# Adopters

This is a list of organizations that have spoken publicly about their adoption or
production users that have added themselves (in alphabetical order):

- [Ad Hoc](http://adhocteam.us/) uses the gRPC-Gateway to serve millions of
  API requests per day.
- [Chef](https://www.chef.io/) uses the gRPC-Gateway to provide the user-facing
  API of [Chef Automate](https://automate.chef.io/). Furthermore, the generated
  OpenAPI data serves as the basis for its [API documentation](https://automate.chef.io/docs/api/).
  The code is Open Source, [see `github.com/chef/automate`](https://github.com/chef/automate).
- [Cho Tot](https://careers.chotot.com/about-us/) utilizes gRPC Gateway to seamlessly integrate HTTP and gRPC services, enabling efficient communication for both legacy and modern systems.
- [Conduit](https://github.com/ConduitIO/conduit), a data streaming tool written in Go,
  uses the gRPC-Gateway since its very beginning to provide an HTTP API in addition to its gRPC API. 
  This makes it easier to integrate with Conduit, and the generated OpenAPI data is used in the documentation.
- [PITS Global Data Recovery Services](https://www.pitsdatarecovery.net/) uses the gRPC-Gateway to generate efficient reverse-proxy servers for internal needs. 
- [Scaleway](https://www.scaleway.com/en/) uses the gRPC-Gateway since 2018 to
  serve millions of API requests per day [1].
- [SpiceDB](https://github.com/authzed/spicedb) uses the gRPC-Gateway to handle
  requests for security-critical permissions checks in environments where gRPC
  is unavailable.

If you have adopted the gRPC-Gateway and would like to be included in this list,
feel free to submit a PR.

[1]: [The odyssey of an HTTP request in Scaleway](https://www.youtube.com/watch?v=eLxD-zIUraE&feature=youtu.be&t=480).
