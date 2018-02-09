'use strict';

var SwaggerClient = require('swagger-client');

describe('ABitOfEverythingService', function() {
  var client;

  beforeEach(function(done) {
    new SwaggerClient({
      url: "http://localhost:8080/swagger/a_bit_of_everything.swagger.json",
      usePromise: true,
    }).then(function(c) {
      client = c;
    }).catch(function(err) {
      done.fail(err);
    }).then(done);
  });

  describe('Create', function() {
    var created;
    var params = {
      float_value: 1.5,
      double_value: 2.5,
      int64_value: "4294967296",
      uint64_value: "9223372036854775807",
      int32_value: -2147483648,
      fixed64_value: "9223372036854775807",
      fixed32_value: 4294967295,
      bool_value: true,
      string_value: "strprefix/foo",
      uint32_value: 4294967295,
      sfixed32_value: 2147483647,
      sfixed64_value: "-4611686018427387904",
      sint32_value: 2147483647,
      sint64_value: "4611686018427387903",
      nonConventionalNameValue: "camelCase",
    };
    var expected = {
      floatValue: 1.5,
      doubleValue: 2.5,
      int64Value: "4294967296",
      uint64Value: "9223372036854775807",
      int32Value: -2147483648,
      fixed64Value: "9223372036854775807",
      fixed32Value: 4294967295,
      boolValue: true,
      stringValue: "strprefix/foo",
      uint32Value: 4294967295,
      sfixed32Value: 2147483647,
      sfixed64Value: "-4611686018427387904",
      sint32Value: 2147483647,
      sint64Value: "4611686018427387903",
      nonConventionalNameValue: "camelCase",

      bytesValue: null,
      singleNested: null,
      nested: [],
      enumValue: "ZERO",
      repeatedStringValue: [],
      mapValue: Object({}),
      mappedStringValue: Object({}),
      mappedNestedValue: Object({}),
      timestampValue: null,
      repeatedEnumValue: [],
    };

    beforeEach(function(done) {
      client.ABitOfEverythingService.Create(params).then(function(resp) {
        created = resp.obj;
      }).catch(function(err) {
        done.fail(err);
      }).then(done);
    });

    it('should assign id', function() {
      expect(created.uuid).not.toBe("");
    });

    it('should echo the request back', function() {
      delete created.uuid;
      expect(created).toEqual(expected);
    });
  });

  describe('CreateBody', function() {
    var created;
    var expected = {
      floatValue: 1.5,
      doubleValue: 2.5,
      int64Value: "4294967296",
      uint64Value: "9223372036854775807",
      int32Value: -2147483648,
      fixed64Value: "9223372036854775807",
      fixed32Value: 4294967295,
      boolValue: true,
      stringValue: "strprefix/foo",
      uint32Value: 4294967295,
      sfixed32Value: 2147483647,
      sfixed64Value: "-4611686018427387904",
      sint32Value: 2147483647,
      sint64Value: "4611686018427387903",
      nonConventionalNameValue: "camelCase",

      nested: [
       { name: "bar", amount: 10, ok: 'FALSE' },
       { name: "baz", amount: 20, ok: 'FALSE' },
      ],
      repeatedStringValue: ["a", "b", "c"],
      oneofString: "x",
      // TODO(yugui) Support enum by name
      mapValue: { a: 1, b: 2 },
      mappedStringValue: { a: "x", b: "y" },
      mappedNestedValue: {
        a: { name: "x", amount: 1, ok: 'FALSE' },
        b: { name: "y", amount: 2, ok: 'FALSE' },
      },

      bytesValue: null,
      singleNested: null,
      enumValue: "ZERO",
      timestampValue: null,
      repeatedEnumValue: [],
    };

    beforeEach(function(done) {
      client.ABitOfEverythingService.CreateBody({
        body: expected,
      }).then(function(resp) {
        created = resp.obj;
      }).catch(function(err) {
        done.fail(err);
      }).then(done);
    });

    it('should assign id', function() {
      expect(created.uuid).not.toBe("");
    });

    it('should echo the request back', function() {
      delete created.uuid;
      expect(created).toEqual(expected);
    });
  });

  describe('lookup', function() {
    var created;
    var expected = {
      bool_value: true,
      string_value: "strprefix/foo",
    };

    beforeEach(function(done) {
      client.ABitOfEverythingService.CreateBody({
        body: expected,
      }).then(function(resp) {
        created = resp.obj;
      }).catch(function(err) {
        fail(err);
      }).finally(done);
    });

    it('should look up an object by uuid', function(done) {
      client.ABitOfEverythingService.Lookup({
        uuid: created.uuid
      }).then(function(resp) {
        expect(resp.obj).toEqual(created);
      }).catch(function(err) {
        fail(err.errObj);
      }).finally(done);
    });

    it('should fail if no such object', function(done) {
      client.ABitOfEverythingService.Lookup({
        uuid: 'not_exist',
      }).then(function(resp) {
        fail('expected failure but succeeded');
      }).catch(function(err) {
        expect(err.status).toBe(404);
      }).finally(done);
    });
  });

  describe('Delete', function() {
    var created;
    var expected = {
      bool_value: true,
      string_value: "strprefix/foo",
    };

    beforeEach(function(done) {
      client.ABitOfEverythingService.CreateBody({
        body: expected,
      }).then(function(resp) {
        created = resp.obj;
      }).catch(function(err) {
        fail(err);
      }).finally(done);
    });

    it('should delete an object by id', function(done) {
      client.ABitOfEverythingService.Delete({
        uuid: created.uuid
      }).then(function(resp) {
        expect(resp.obj).toEqual({});
      }).catch(function(err) {
        fail(err.errObj);
      }).then(function() {
        return client.ABitOfEverythingService.Lookup({
          uuid: created.uuid
        });
      }).then(function(resp) {
        fail('expected failure but succeeded');
      }). catch(function(err) {
        expect(err.status).toBe(404);
      }).finally(done);
    });
  });
});

