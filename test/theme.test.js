const { describe, it } = require('node:test');
const assert = require('node:assert/strict');

describe('Theme', () => {
  it('parseArgs handles --theme flag', () => {
    const { parseArgs } = require('../src/utils');

    const args1 = parseArgs(['--theme', 'light']);
    assert.equal(args1.theme, 'light');

    const args2 = parseArgs(['--theme', 'dark']);
    assert.equal(args2.theme, 'dark');

    const args3 = parseArgs(['--theme', 'auto']);
    assert.equal(args3.theme, 'auto');
  });

  it('parseArgs rejects invalid theme values', () => {
    const { parseArgs } = require('../src/utils');
    const args = parseArgs(['--theme', 'invalid']);
    assert.equal(args.theme, undefined);
  });

  it('base template includes theme detection script', () => {
    const { createBaseTemplate, generateBreadcrumb } = require('../src/templates/base');
    const { escapeHtml } = require('../src/utils');

    const html = createBaseTemplate({
      title: 'Test',
      breadcrumb: generateBreadcrumb('/', escapeHtml),
      content: '<p>Hello</p>'
    });

    assert.ok(html.includes('lookit-theme'), 'Should include theme localStorage key');
    assert.ok(html.includes('data-theme'), 'Should reference data-theme attribute');
  });

  it('styles include light theme variables', () => {
    const { baseStyles } = require('../src/styles');
    assert.ok(baseStyles.includes('prefers-color-scheme: light'), 'Should have light theme media query');
    assert.ok(baseStyles.includes('data-theme'), 'Should have data-theme selectors');
  });
});
