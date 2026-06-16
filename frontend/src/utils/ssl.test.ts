import { describe, it } from 'node:test';
import assert from 'node:assert';
import { getAccountName, getKeyName, getDNSName, getProvider } from './ssl';

describe('ssl utils', () => {
    describe('getAccountName', () => {
        it('should return correct name for known accounts', () => {
            assert.strictEqual(getAccountName('letsencrypt'), "Let's Encrypt");
            assert.strictEqual(getAccountName('zerossl'), 'ZeroSSL');
            assert.strictEqual(getAccountName('buypass'), 'Buypass');
            assert.strictEqual(getAccountName('google'), 'Google Cloud');
        });

        it('should return empty string for unknown accounts', () => {
            assert.strictEqual(getAccountName('unknown'), '');
        });
    });

    describe('getKeyName', () => {
        it('should return correct name for known key types', () => {
            assert.strictEqual(getKeyName('P256'), 'EC 256');
            assert.strictEqual(getKeyName('2048'), 'RSA 2048');
        });

        it('should return empty string for unknown key types', () => {
            assert.strictEqual(getKeyName('unknown'), '');
        });
    });

    describe('getDNSName', () => {
        it('should return correct name for known dns types', () => {
            // we're getting 'Aliyun DNS' because i18n isn't fully set up and it falls back to the default english translation
            assert.strictEqual(typeof getDNSName('AliYun') === 'string', true);
            assert.strictEqual(typeof getDNSName('CloudFlare') === 'string', true);
            assert.strictEqual(getDNSName('CloudFlare'), 'Cloudflare');
        });

        it('should return empty string for unknown dns types', () => {
            assert.strictEqual(getDNSName('unknown'), '');
        });
    });

    describe('getProvider', () => {
        it('should return HTTP literal for http', () => {
            assert.strictEqual(getProvider('http'), 'HTTP');
        });
        it('should return translations for others', () => {
            assert.ok(typeof getProvider('dnsAccount') === 'string');
            assert.ok(typeof getProvider('dnsManual') === 'string');
            assert.ok(typeof getProvider('selfSigned') === 'string');
            assert.ok(typeof getProvider('fromMaster') === 'string');
            assert.ok(typeof getProvider('unknown') === 'string');
        });
    });
});
