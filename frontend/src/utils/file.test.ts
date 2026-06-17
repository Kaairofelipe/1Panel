import { describe, it, beforeEach, afterEach } from 'node:test';
import assert from 'node:assert';
import { buildFileSharePageUrl } from './file';

describe('buildFileSharePageUrl', () => {
    let originalWindow: any;

    beforeEach(() => {
        // Save original window if it exists
        if (typeof global !== 'undefined') {
            originalWindow = (global as any).window;
            // Mock global window object
            (global as any).window = {
                location: {
                    origin: 'http://localhost:8080',
                },
            };
        }
    });

    afterEach(() => {
        // Restore original window
        if (typeof global !== 'undefined') {
            (global as any).window = originalWindow;
        }
    });

    it('should build the correct share page URL with basic parameters', () => {
        const url = buildFileSharePageUrl('myCode123', 'local');
        assert.strictEqual(url, 'http://localhost:8080/s/myCode123?operateNode=local');
    });

    it('should properly encode special characters in the share code', () => {
        const url = buildFileSharePageUrl('my Code/test+123@!', 'localNode');
        assert.strictEqual(url, 'http://localhost:8080/s/my%20Code%2Ftest%2B123%40!?operateNode=localNode');
    });

    it('should handle empty current node parameter', () => {
        const url = buildFileSharePageUrl('validCode', '');
        assert.strictEqual(url, 'http://localhost:8080/s/validCode?operateNode=');
    });
});
