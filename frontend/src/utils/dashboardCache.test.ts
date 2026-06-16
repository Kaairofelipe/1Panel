import { test, describe, beforeEach, mock, afterEach } from 'node:test';
import assert from 'node:assert';
import { getDashboardCache, DASHBOARD_CACHE_KEY } from './dashboardCache.ts';

describe('getDashboardCache', () => {
    let mockGetItem: any;

    beforeEach(() => {
        mockGetItem = mock.fn();
        global.localStorage = {
            getItem: mockGetItem,
            setItem: mock.fn(),
            removeItem: mock.fn(),
        } as any;
    });

    afterEach(() => {
        mock.restoreAll();
    });

    test('should return null when localStorage returns null (empty cache)', () => {
        mockGetItem.mock.mockImplementation(() => null);
        const result = getDashboardCache('testKey');
        assert.strictEqual(result, null);
    });

    test('should return null when JSON.parse throws an error (invalid JSON)', () => {
        mockGetItem.mock.mockImplementation(() => '{invalid_json}');
        const result = getDashboardCache('testKey');
        assert.strictEqual(result, null);
    });

    test('should return null if the key does not exist in the parsed cache', () => {
        const cacheData = JSON.stringify({ otherKey: { value: 'data', expireAt: Date.now() + 10000 } });
        mockGetItem.mock.mockImplementation(() => cacheData);
        const result = getDashboardCache('testKey');
        assert.strictEqual(result, null);
    });

    test('should return the value if the key exists and expireAt is in the future', () => {
        const futureTime = Date.now() + 10000;
        const cacheData = JSON.stringify({ testKey: { value: 'testData', expireAt: futureTime } });
        mockGetItem.mock.mockImplementation(() => cacheData);

        const result = getDashboardCache('testKey');
        assert.strictEqual(result, 'testData');
    });

    test('should return null if the key exists but expireAt is in the past', () => {
        const pastTime = Date.now() - 10000;
        const cacheData = JSON.stringify({ testKey: { value: 'testData', expireAt: pastTime } });
        mockGetItem.mock.mockImplementation(() => cacheData);

        const result = getDashboardCache('testKey');
        assert.strictEqual(result, null);
    });

    test('should return null if the key exists but expireAt is exactly now', () => {
        const now = Date.now();
        const cacheData = JSON.stringify({ testKey: { value: 'testData', expireAt: now } });
        mockGetItem.mock.mockImplementation(() => cacheData);

        // Mock Date.now() so it consistently returns the exact same 'now'
        const originalDateNow = Date.now;
        Date.now = () => now;

        const result = getDashboardCache('testKey');

        // Restore original Date.now
        Date.now = originalDateNow;

        assert.strictEqual(result, null);
    });
});
