import { test, describe } from 'node:test';
import * as assert from 'node:assert';
import { isJson } from './misc.ts';

describe('isJson', () => {
    test('should return true for a valid JSON object string', () => {
        assert.strictEqual(isJson('{"key": "value"}'), true);
    });

    test('should return true for a valid JSON array string', () => {
        assert.strictEqual(isJson('[1, 2, 3]'), true);
    });

    test('should return false for an invalid JSON string', () => {
        assert.strictEqual(isJson('{invalid}'), false);
    });

    test('should return false for valid JSON that is not an object (string)', () => {
        assert.strictEqual(isJson('"string"'), false);
    });

    test('should return false for valid JSON that is not an object (number)', () => {
        assert.strictEqual(isJson('123'), false);
    });

    test('should return false for valid JSON that is not an object (boolean)', () => {
        assert.strictEqual(isJson('true'), false);
        assert.strictEqual(isJson('false'), false);
    });

    test('should return false for "null" string since we do not want it to be valid json object', () => {
        assert.strictEqual(isJson('null'), false);
    });

    test('should return false for empty string', () => {
        assert.strictEqual(isJson(''), false);
    });

    test('should return false for spaces', () => {
        assert.strictEqual(isJson('   '), false);
    });

    test('should return false for undefined', () => {
        assert.strictEqual(isJson(undefined as any), false);
    });

    test('should return false for null as input', () => {
        assert.strictEqual(isJson(null as any), false);
    });
});
