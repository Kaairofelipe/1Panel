import test from 'node:test';
import assert from 'node:assert';
import { deepCopy } from './misc.ts';

test('deepCopy', async (t) => {
    await t.test('should copy a plain object with primitive values', () => {
        const original = { name: 'Alice', age: 30, active: true };
        const copy = deepCopy(original);
        assert.deepStrictEqual(copy, original);
        assert.notStrictEqual(copy, original);
    });

    await t.test('should perform a deep copy of nested objects', () => {
        const original = { user: { profile: { id: 1 } } };
        const copy = deepCopy(original);
        assert.deepStrictEqual(copy, original);
        assert.notStrictEqual(copy, original);
        assert.notStrictEqual(copy.user, original.user);
        assert.notStrictEqual(copy.user.profile, original.user.profile);
    });

    await t.test('should handle arrays', () => {
        const original = [1, 2, 3];
        const copy = deepCopy(original);
        assert.deepStrictEqual(copy, original);
        assert.notStrictEqual(copy, original);
        assert.ok(Array.isArray(copy));
    });

    await t.test('should handle objects with nested arrays', () => {
        const original = { data: [ { id: 1 }, { id: 2 } ] };
        const copy = deepCopy(original);
        assert.deepStrictEqual(copy, original);
        assert.notStrictEqual(copy, original);
        assert.notStrictEqual(copy.data, original.data);
        assert.notStrictEqual(copy.data[0], original.data[0]);
    });

    await t.test('should fallback to object for non-array/non-object like numbers', () => {
        const original = 123;
        const copy = deepCopy(original);
        assert.deepStrictEqual(copy, {});
    });

    await t.test('should fallback to object for null', () => {
        const copy = deepCopy(null);
        assert.deepStrictEqual(copy, {});
    });
});
