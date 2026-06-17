'use strict';

var SwaggerClient = require('swagger-client').default || require('swagger-client');

describe('EchoService', function () {
  var client;

  beforeEach(function (done) {
    new SwaggerClient({
      url: "http://localhost:8080/openapiv2/echo_service.swagger.json",
    }).then(function (c) {
      client = c;
      done();
    }).catch(function (err) {
      done.fail(err);
    });
  });

  describe('Echo', function () {
    it('should echo the request back', function (done) {
      var expected = {
        id: "foo",
        num: "0",
        status: null,
        resourceId: '',
        nId: null
      };
      client.apis.EchoService.EchoService_Echo(
        expected,
        { responseContentType: "application/json" }
      ).then(function (resp) {
        expect(resp.obj).toEqual(expected);
      }).catch(function (err) {
        done.fail(err);
      }).then(done);
    });
  });

  describe('EchoBody', function () {
    it('should echo the request back', function (done) {
      var expected = {
        id: "foo",
        num: "0",
        status: null,
        resourceId: '',
        nId: null
      };
      client.apis.EchoService.EchoService_EchoBody(
        { body: expected },
        { responseContentType: "application/json" }
      ).then(function (resp) {
        expect(resp.obj).toEqual(expected);
      }).catch(function (err) {
        done.fail(err);
      }).then(done);
    });
  });
});
