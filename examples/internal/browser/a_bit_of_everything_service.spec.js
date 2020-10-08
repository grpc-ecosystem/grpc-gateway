'use strict';

var SwaggerClient = require('swagger-client');

describe('ABitOfEverythingService', function() {
  var client;

  beforeEach(function(done) {
    new SwaggerClient({
      url: "http://localhost:8080/openapiv2/a_bit_of_everything.swagger.json",
      usePromise: true,
    }).then(function(c) {
      client = c;
    }).catch(function(err) {
      done.fail(err);
    }).then(done);
  });

  describe('Create', function() {
    var created;
    var expected = {
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
      enum_value: "ONE",
      path_enum_value: "DEF",
      nested_path_enum_value: "JKL",
      enum_value_annotation: "ONE",
    };

    beforeEach(function(done) {
      client.ABitOfEverythingService.Create(expected).then(function(resp) {
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
      enum_value: "ONE",
      path_enum_value: "DEF",
      nested_path_enum_value: "JKL",

      nested: [
       { name: "bar", amount: 10 },
       { name: "baz", amount: 20 },
      ],
      repeated_string_value: ["a", "b", "c"],
      oneof_string: "x",
      map_value: { a: "ONE", b: 2 },
      mapped_string_value: { a: "x", b: "y" },
      mapped_nested_value: {
        a: { name: "x", amount: 1 },
        b: { name: "y", amount: 2 },
      },
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

  describe('GetRepeatedQuery', function() {
    var repeated;
    var expected = {
      path_repeated_float_value: [1.5, -1.5],
      path_repeated_double_value: [2.5, -2.5],
      path_repeated_int64_value: ["4294967296", "-4294967296"],
      path_repeated_uint64_value: ["0", "9223372036854775807"],
      path_repeated_int32_value: [2147483647, -2147483648],
      path_repeated_fixed64_value: ["0", "9223372036854775807"],
      path_repeated_fixed32_value: [0, 4294967295],
      path_repeated_bool_value: [true, false],
      path_repeated_string_value: ["foo", "bar"],
      path_repeated_bytes_value: ["AA==", "_w=="],
      path_repeated_uint32_value: [4294967295, 0],
      path_repeated_enum_value: ["ONE", "ONE"],
      path_repeated_sfixed32_value: [-2147483648, 2147483647],
      path_repeated_sfixed64_value: ["-4294967296", "4294967296"],
      path_repeated_sint32_value: [2147483646, -2147483647],
      path_repeated_sint64_value: ["4611686018427387903", "-4611686018427387904"]
    };

    beforeEach(function(done) {
      client.ABitOfEverythingService.GetRepeatedQuery(expected).then(function(resp) {
        repeated = resp.obj;
      }).catch(function(err) {
        done.fail(err);
      }).then(done);
    });

    it('should echo the request back', function() {
      // API will echo a non URL safe encoding
      expected.path_repeated_bytes_value = ["AA==", "/w=="];
      expect(repeated).toEqual(expected);
    });
  });
});

