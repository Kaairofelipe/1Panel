import test from 'node:test';
import assert from 'node:assert';
import { compareVersion } from './version';

test('compareVersion', async (t) => {
    await t.test('returns true for identical versions', () => {
        assert.strictEqual(compareVersion('1.0.0', '1.0.0'), true);
        assert.strictEqual(compareVersion('2.1', '2.1'), true);
        assert.strictEqual(compareVersion('0.0.1', '0.0.1'), true);
    });

    await t.test('returns true when version1 is greater than version2 (same length)', () => {
        assert.strictEqual(compareVersion('2.0.0', '1.0.0'), true);
        assert.strictEqual(compareVersion('1.1.0', '1.0.0'), true);
        assert.strictEqual(compareVersion('1.0.1', '1.0.0'), true);
    });

    await t.test('returns false when version1 is less than version2 (same length)', () => {
        assert.strictEqual(compareVersion('1.0.0', '2.0.0'), false);
        assert.strictEqual(compareVersion('1.0.0', '1.1.0'), false);
        assert.strictEqual(compareVersion('1.0.0', '1.0.1'), false);
    });

    await t.test('handles versions of different lengths', () => {
        // v1 longer, v1 > v2
        assert.strictEqual(compareVersion('1.0.1', '1.0'), true);
        // v1 shorter, v1 < v2
        assert.strictEqual(compareVersion('1.0', '1.0.1'), false);

        // Equal values but different string lengths
        assert.strictEqual(compareVersion('1.0.0', '1.0'), true);
        assert.strictEqual(compareVersion('1.0', '1.0.0'), true);

        // v1 longer, v1 < v2
        assert.strictEqual(compareVersion('1.0.0', '1.1'), false);
        // v1 shorter, v1 > v2
        assert.strictEqual(compareVersion('1.1', '1.0.0'), true);
    });

    await t.test('handles non-digit prefixes and suffixes', () => {
        assert.strictEqual(compareVersion('v1.0.0', '1.0.0'), true);
        assert.strictEqual(compareVersion('v2.0.0', '1.0.0'), true);
        assert.strictEqual(compareVersion('1.0.0', 'v2.0.0'), false);

        // Suffixes
        assert.strictEqual(compareVersion('1.0.0-rc2', '1.0.0-rc1'), true);
        assert.strictEqual(compareVersion('1.0.0-rc1', '1.0.0-rc2'), false);
        assert.strictEqual(compareVersion('1.0.0-alpha', '1.0.0-alpha'), true);
    });

    await t.test('handles complex string versions', () => {
        assert.strictEqual(compareVersion('1.2.3.4.5', '1.2.3.4'), true);
        assert.strictEqual(compareVersion('1.0-beta.1', '1.0-alpha.2'), false); // alpha vs beta numbers are stripped, so '1.0.1' vs '1.0.2'
    });
});
