import { test, describe, beforeEach } from 'node:test';
import assert from 'node:assert';
import { getDashboardCache, setDashboardCache, clearDashboardCache, clearDashboardCacheByPrefix, DASHBOARD_CACHE_KEY } from './dashboardCache.ts';

describe('dashboardCache', () => {
    let mockStorage: Record<string, string> = {};

    beforeEach(() => {
        mockStorage = {};
        global.localStorage = {
            getItem: (key: string) => mockStorage[key] || null,
            setItem: (key: string, value: string) => { mockStorage[key] = value; },
            removeItem: (key: string) => { delete mockStorage[key]; },
            clear: () => { mockStorage = {}; },
            length: 0,
            key: () => null,
        } as unknown as Storage;
    });

    test('should return null when reading invalid JSON from cache', () => {
        // Set invalid JSON
        mockStorage[DASHBOARD_CACHE_KEY] = 'invalid-json{';

        // When getDashboardCache calls readCache(), it should hit the catch block and return null
        const result = getDashboardCache('some-key');

        assert.strictEqual(result, null);
    });

    test('should return value when reading valid JSON from cache', () => {
        // Set valid JSON with an unexpired entry
        mockStorage[DASHBOARD_CACHE_KEY] = JSON.stringify({
            'some-key': {
                value: 'some-value',
                expireAt: Date.now() + 10000, // 10 seconds in the future
            }
        });

        const result = getDashboardCache('some-key');

        assert.strictEqual(result, 'some-value');
    });

    test('should return null when reading expired entry', () => {
        // Set valid JSON with an expired entry
        mockStorage[DASHBOARD_CACHE_KEY] = JSON.stringify({
            'some-key': {
                value: 'some-value',
                expireAt: Date.now() - 10000, // 10 seconds in the past
            }
        });

        const result = getDashboardCache('some-key');

        assert.strictEqual(result, null);
    });

    test('should set value correctly', () => {
        setDashboardCache('some-key', 'some-value', 10000);

        const rawCache = mockStorage[DASHBOARD_CACHE_KEY];
        assert.ok(rawCache);

        const parsedCache = JSON.parse(rawCache);
        assert.ok(parsedCache['some-key']);
        assert.strictEqual(parsedCache['some-key'].value, 'some-value');
    });

    test('should clear cache entirely when setting with invalid JSON', () => {
        // Set invalid JSON
        mockStorage[DASHBOARD_CACHE_KEY] = 'invalid-json{';

        setDashboardCache('some-key', 'some-value', 10000);

        // Instead of error, it should remove the key
        assert.strictEqual(mockStorage[DASHBOARD_CACHE_KEY], undefined);
    });

    test('should clear all cache', () => {
        mockStorage[DASHBOARD_CACHE_KEY] = JSON.stringify({
            'some-key': { value: 'some-value', expireAt: Date.now() + 10000 }
        });

        clearDashboardCache();

        assert.strictEqual(mockStorage[DASHBOARD_CACHE_KEY], undefined);
    });

    test('should clear cache by prefix', () => {
        mockStorage[DASHBOARD_CACHE_KEY] = JSON.stringify({
            'prefix1-key': { value: 'val1', expireAt: Date.now() + 10000 },
            'prefix2-key': { value: 'val2', expireAt: Date.now() + 10000 },
            'other-key': { value: 'val3', expireAt: Date.now() + 10000 },
        });

        clearDashboardCacheByPrefix(['prefix1-', 'prefix2-']);

        const rawCache = mockStorage[DASHBOARD_CACHE_KEY];
        const parsedCache = JSON.parse(rawCache);

        assert.strictEqual(parsedCache['prefix1-key'], undefined);
        assert.strictEqual(parsedCache['prefix2-key'], undefined);
        assert.strictEqual(parsedCache['other-key'].value, 'val3');
    });

    test('should clear entire cache when clearing by prefix with invalid JSON', () => {
        mockStorage[DASHBOARD_CACHE_KEY] = 'invalid-json{';

        clearDashboardCacheByPrefix(['prefix1-']);

        assert.strictEqual(mockStorage[DASHBOARD_CACHE_KEY], undefined);
    });
});
