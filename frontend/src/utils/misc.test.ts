import { test, describe } from 'node:test';
import assert from 'node:assert';
import { deepCopy } from './misc.ts';

describe('deepCopy', () => {
    test('should copy primitive values', () => {
        assert.strictEqual(deepCopy(5), 5);
        assert.strictEqual(deepCopy('string'), 'string');
        assert.strictEqual(deepCopy(true), true);
        assert.strictEqual(deepCopy(null), null);
    });

    test('should copy arrays', () => {
        const arr = [1, 2, { a: 3 }];
        const copied = deepCopy(arr);
        assert.notStrictEqual(copied, arr);
        assert.deepStrictEqual(copied, arr);
        assert.notStrictEqual(copied[2], arr[2]);
    });

    test('should copy objects', () => {
        const obj = { a: 1, b: { c: 2 } };
        const copied = deepCopy(obj);
        assert.notStrictEqual(copied, obj);
        assert.deepStrictEqual(copied, obj);
        assert.notStrictEqual(copied.b, obj.b);
    });

    test('should maintain type', () => {
        type MyType = { id: number };
        const input: MyType = { id: 1 };
        const copied: MyType = deepCopy(input);
        assert.deepStrictEqual(copied, input);
    });
});
