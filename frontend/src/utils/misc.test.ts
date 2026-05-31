import { test } from 'node:test';
import assert from 'node:assert';
import { deepCopy } from './misc.ts';

test('deepCopy - Objects', () => {
    const obj = { a: 1, b: 'string', c: true };
    const copy = deepCopy<typeof obj>(obj);
    assert.deepStrictEqual(copy, obj);
    assert.notStrictEqual(copy, obj);
});

test('deepCopy - Arrays', () => {
    const arr = [1, 2, 3];
    const copy = deepCopy<number[]>(arr);
    assert.deepStrictEqual(copy, arr);
    assert.notStrictEqual(copy, arr);
    assert.ok(Array.isArray(copy));
});

test('deepCopy - Nested Structures', () => {
    const obj = {
        a: [1, 2, { x: 1 }],
        b: { y: [3, 4] }
    };
    const copy = deepCopy<typeof obj>(obj);
    assert.deepStrictEqual(copy, obj);
    assert.notStrictEqual(copy, obj);
    assert.notStrictEqual(copy.a, obj.a);
    assert.notStrictEqual(copy.a[2], obj.a[2]);
    assert.notStrictEqual(copy.b, obj.b);
    assert.notStrictEqual((copy.b as any).y, (obj.b as any).y);
});

test('deepCopy - null and undefined (Edge Cases)', () => {
    // Current implementation returns {} for null and undefined
    // This test will fail if we expect it to return null/undefined
    assert.strictEqual(deepCopy(null), null, 'Should handle null');
    assert.strictEqual(deepCopy(undefined), undefined, 'Should handle undefined');
});

test('deepCopy - Primitives', () => {
    assert.strictEqual(deepCopy(123), 123);
    assert.strictEqual(deepCopy('test'), 'test');
    assert.strictEqual(deepCopy(true), true);
});

test('deepCopy - Object containing null', () => {
    const obj = { a: null };
    const copy = deepCopy<any>(obj);
    assert.strictEqual(copy.a, null, 'null inside object should remain null');
});
