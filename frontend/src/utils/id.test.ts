import { test, describe } from 'node:test';
import assert from 'node:assert/strict';
import { getRandomStr } from './id.ts';

describe('getRandomStr', () => {
    test('returns a string of the specified length', () => {
        const lengths = [0, 1, 5, 10, 100];
        for (const length of lengths) {
            const result = getRandomStr(length);
            assert.equal(typeof result, 'string');
            assert.equal(result.length, length);
        }
    });

    test('returns a string containing only characters from the allowed set', () => {
        const allowedChars = 'ABCDEFGHJKMNPQRSTWXYZabcdefhijkmnprstwxyz2345678';
        const length = 100;
        const result = getRandomStr(length);

        for (const char of result) {
            assert.ok(allowedChars.includes(char), `Character "${char}" is not in the allowed set`);
        }
    });

    test('handles negative lengths gracefully (returns empty string)', () => {
        const length = -5;
        const result = getRandomStr(length);
        assert.equal(result, '');
    });
});
