import { test, describe } from 'node:test';
import * as assert from 'node:assert/strict';
import { deepCopy } from './misc.ts';

describe('deepCopy', () => {
    test('deep copies an object', () => {
        const obj = { a: 1, b: { c: 2 } };
        const copy = deepCopy(obj);
        assert.deepEqual(copy, obj);
        assert.notEqual(copy, obj);
        assert.notEqual(copy.b, obj.b);
    });

    test('deep copies an array', () => {
        const arr = [1, [2, 3]];
        const copy = deepCopy(arr);
        assert.deepEqual(copy, arr);
        assert.notEqual(copy, arr);
        assert.notEqual(copy[1], arr[1]);
        assert.ok(Array.isArray(copy));
    });

    test('handles object where .push throws an error', () => {
        const obj = {
            a: 1
        };
        // The error will be thrown when checking push. Let's make sure the iteration also skips throwing or works with it.
        Object.defineProperty(obj, 'push', {
            get() {
                throw new Error("Accessing push is not allowed");
            },
            enumerable: false // Keep it non-enumerable so `for (let attr in obj)` doesn't trigger the getter again.
        });

        const copy = deepCopy(obj);
        assert.equal(copy.a, 1);
        assert.ok(!Array.isArray(copy));
    });
});
