import test from 'node:test';
import assert from 'node:assert';
import { clearDashboardCache, DASHBOARD_CACHE_KEY } from './dashboardCache.ts';

test('clearDashboardCache', async (t) => {
  const originalLocalStorage = (globalThis as any).localStorage;

  t.beforeEach(() => {
    // Mock localStorage
    const mockLocalStorage = (() => {
      let store: Record<string, string> = {};
      return {
        getItem: (key: string) => store[key] || null,
        setItem: (key: string, value: string) => {
          store[key] = value.toString();
        },
        removeItem: (key: string) => {
          delete store[key];
        },
        clear: () => {
          store = {};
        }
      };
    })();

    (globalThis as any).localStorage = mockLocalStorage;
  });

  t.afterEach(() => {
    // Restore original localStorage
    (globalThis as any).localStorage = originalLocalStorage;
  });

  await t.test('should remove DASHBOARD_CACHE_KEY from localStorage', () => {
    // Arrange
    localStorage.setItem(DASHBOARD_CACHE_KEY, JSON.stringify({ key1: { value: 1, expireAt: Date.now() + 1000 }}));
    assert.strictEqual(localStorage.getItem(DASHBOARD_CACHE_KEY) !== null, true);

    // Act
    clearDashboardCache();

    // Assert
    assert.strictEqual(localStorage.getItem(DASHBOARD_CACHE_KEY), null);
  });
});
