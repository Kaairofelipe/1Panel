import test, { describe, it, mock, beforeEach, afterEach } from 'node:test';
import assert from 'node:assert';
import { toLink, preloadImage } from './browser.ts';

describe('browser utils', () => {
    describe('toLink', () => {
        let originalWindow: any;
        let openMock: any;

        beforeEach(() => {
            originalWindow = global.window;
            openMock = mock.fn();
            global.window = { open: openMock } as any;
        });

        afterEach(() => {
            global.window = originalWindow;
        });

        it('should open normal http links without modification', () => {
            toLink('http://example.com/path');
            assert.strictEqual(openMock.mock.calls.length, 1);
            assert.strictEqual(openMock.mock.calls[0].arguments[0], 'http://example.com/path');
            assert.strictEqual(openMock.mock.calls[0].arguments[1], '_blank');
        });

        it('should open normal https links without modification', () => {
            toLink('https://example.com');
            assert.strictEqual(openMock.mock.calls.length, 1);
            assert.strictEqual(openMock.mock.calls[0].arguments[0], 'https://example.com');
            assert.strictEqual(openMock.mock.calls[0].arguments[1], '_blank');
        });

        it('should open IPv4 links without modification', () => {
            toLink('http://127.0.0.1:8080/api');
            assert.strictEqual(openMock.mock.calls[0].arguments[0], 'http://127.0.0.1:8080/api');
        });

        it('should format unbracketed IPv6 addresses correctly (HTTP)', () => {
            toLink('http://::1:8080/path');
            assert.strictEqual(openMock.mock.calls[0].arguments[0], 'http://[::1]:8080/path');
        });

        it('should format unbracketed IPv6 addresses correctly (HTTPS)', () => {
            toLink('https://fe80::1ff:fe23:4567:890a:443');
            assert.strictEqual(openMock.mock.calls[0].arguments[0], 'https://[fe80::1ff:fe23:4567:890a]:443');
        });

        it('should handle unbracketed IPv6 addresses with paths correctly', () => {
            toLink('http://2001:db8::1:80/some/path?query=1');
            assert.strictEqual(openMock.mock.calls[0].arguments[0], 'http://[2001:db8::1]:80/some/path?query=1');
        });

        it('should not modify already bracketed IPv6 addresses', () => {
            toLink('http://[::1]:8080/path');
            assert.strictEqual(openMock.mock.calls[0].arguments[0], 'http://[::1]:8080/path');
        });

        it('should catch and ignore exceptions thrown by window.open', () => {
            const errorMock = mock.fn(() => {
                throw new Error('Access denied');
            });
            global.window = { open: errorMock } as any;

            // Should not throw
            assert.doesNotThrow(() => {
                toLink('http://example.com');
            });
            assert.strictEqual(errorMock.mock.calls.length, 1);
        });

        it('should safely handle missing window object', () => {
            // Remove window object entirely
            global.window = undefined as any;

            // Should not throw (handled by try/catch)
            assert.doesNotThrow(() => {
                toLink('http://example.com');
            });
        });
    });
});
