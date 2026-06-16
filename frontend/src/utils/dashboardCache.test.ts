import { test, describe, beforeEach, afterEach } from 'node:test';
import assert from 'node:assert';
import { setDashboardCache, DASHBOARD_CACHE_KEY } from './dashboardCache.ts';

describe('setDashboardCache', () => {
    let originalLocalStorage: any;
    let mockStorage: Record<string, string> = {};

    beforeEach(() => {
        originalLocalStorage = global.localStorage;
        mockStorage = {};
        global.localStorage = {
            getItem: (key: string) => mockStorage[key] || null,
            setItem: (key: string, value: string) => {
                mockStorage[key] = value;
            },
            removeItem: (key: string) => {
                delete mockStorage[key];
            },
            clear: () => {
                mockStorage = {};
            },
        } as any;
    });

    afterEach(() => {
        global.localStorage = originalLocalStorage;
    });

    test('should save new value to localStorage', () => {
        const now = Date.now();
        setDashboardCache('testKey', 'testValue', 1000);

        const storedRaw = mockStorage[DASHBOARD_CACHE_KEY];
        assert.ok(storedRaw);

        const stored = JSON.parse(storedRaw);
        assert.strictEqual(stored['testKey'].value, 'testValue');
        assert.ok(stored['testKey'].expireAt >= now + 1000);
    });

    test('should append to existing cache in localStorage', () => {
        const now = Date.now();
        mockStorage[DASHBOARD_CACHE_KEY] = JSON.stringify({
            existingKey: {
                value: 'existingValue',
                expireAt: now + 5000,
            }
        });

        setDashboardCache('newKey', 'newValue', 1000);

        const storedRaw = mockStorage[DASHBOARD_CACHE_KEY];
        assert.ok(storedRaw);

        const stored = JSON.parse(storedRaw);
        assert.strictEqual(stored['existingKey'].value, 'existingValue');
        assert.strictEqual(stored['newKey'].value, 'newValue');
        assert.ok(stored['newKey'].expireAt >= now + 1000);
    });

    test('should overwrite existing key in localStorage', () => {
        const now = Date.now();
        mockStorage[DASHBOARD_CACHE_KEY] = JSON.stringify({
            testKey: {
                value: 'oldValue',
                expireAt: now + 5000,
            }
        });

        setDashboardCache('testKey', 'newValue', 1000);

        const storedRaw = mockStorage[DASHBOARD_CACHE_KEY];
        assert.ok(storedRaw);

        const stored = JSON.parse(storedRaw);
        assert.strictEqual(stored['testKey'].value, 'newValue');
        assert.ok(stored['testKey'].expireAt >= now + 1000);
        assert.ok(stored['testKey'].expireAt < now + 2000);
    });

    test('should clear cache if JSON.parse fails', () => {
        mockStorage[DASHBOARD_CACHE_KEY] = 'invalid json';

        setDashboardCache('testKey', 'testValue', 1000);

        const storedRaw = mockStorage[DASHBOARD_CACHE_KEY];
        assert.strictEqual(storedRaw, undefined); // removeItem should be called
    });

    test('should clear cache if localStorage.setItem fails', () => {
        global.localStorage.setItem = () => {
            throw new Error('QuotaExceededError');
        };

        setDashboardCache('testKey', 'testValue', 1000);

        const storedRaw = mockStorage[DASHBOARD_CACHE_KEY];
        assert.strictEqual(storedRaw, undefined); // removeItem should be called
    });
});
