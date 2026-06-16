import test from 'node:test';
import assert from 'node:assert';
import { clearDashboardCacheByPrefix, DASHBOARD_CACHE_KEY } from './dashboardCache.ts';

test('clearDashboardCacheByPrefix handles empty cache', () => {
    const originalLocalStorage = (global as any).localStorage;
    let storage: Record<string, string> = {};
    (global as any).localStorage = {
        getItem: (key: string) => storage[key] || null,
        setItem: (key: string, value: string) => { storage[key] = value; },
        removeItem: (key: string) => { delete storage[key]; },
        clear: () => { storage = {}; }
    };

    try {
        clearDashboardCacheByPrefix(['prefix_']);
        assert.strictEqual(localStorage.getItem(DASHBOARD_CACHE_KEY), null);
    } finally {
        (global as any).localStorage = originalLocalStorage;
    }
});

test('clearDashboardCacheByPrefix handles invalid JSON', () => {
    const originalLocalStorage = (global as any).localStorage;
    let storage: Record<string, string> = {};
    (global as any).localStorage = {
        getItem: (key: string) => storage[key] || null,
        setItem: (key: string, value: string) => { storage[key] = value; },
        removeItem: (key: string) => { delete storage[key]; },
        clear: () => { storage = {}; }
    };

    try {
        localStorage.setItem(DASHBOARD_CACHE_KEY, 'invalid json');

        clearDashboardCacheByPrefix(['prefix_']);
        assert.strictEqual(localStorage.getItem(DASHBOARD_CACHE_KEY), null); // Because clearDashboardCache is called on catch
    } finally {
        (global as any).localStorage = originalLocalStorage;
    }
});

test('clearDashboardCacheByPrefix deletes only keys matching given prefixes', () => {
    const originalLocalStorage = (global as any).localStorage;
    let storage: Record<string, string> = {};
    (global as any).localStorage = {
        getItem: (key: string) => storage[key] || null,
        setItem: (key: string, value: string) => { storage[key] = value; },
        removeItem: (key: string) => { delete storage[key]; },
        clear: () => { storage = {}; }
    };

    try {
        const initialCache = {
            'app1_data': 'some data',
            'app2_data': 'some other data',
            'system_config': 'config value',
            'app1_settings': 'settings value',
            'app3_data': 'app3 data'
        };

        localStorage.setItem(DASHBOARD_CACHE_KEY, JSON.stringify(initialCache));

        clearDashboardCacheByPrefix(['app1_', 'system_']);

        const rawCache = localStorage.getItem(DASHBOARD_CACHE_KEY);
        assert.ok(rawCache);

        const parsedCache = JSON.parse(rawCache);
        assert.deepStrictEqual(parsedCache, {
            'app2_data': 'some other data',
            'app3_data': 'app3 data'
        });
    } finally {
        (global as any).localStorage = originalLocalStorage;
    }
});
