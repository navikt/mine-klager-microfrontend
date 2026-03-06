import { describe, expect, test } from 'bun:test';
import { findComponentCss } from './find-component-css';

describe('findComponentCss', () => {
  test('returns linkcard CSS for HTML with aksel-link-card class', () => {
    const html = '<div class="aksel-link-card">content</div>';
    const result = findComponentCss(html);

    expect(result.files).toContain('linkcard.min.css');
    expect(result.css).toContain('.aksel-link-card');
  });

  test('returns typography CSS for HTML with aksel-heading class', () => {
    const html = '<h2 class="aksel-heading aksel-heading--small">Title</h2>';
    const result = findComponentCss(html);

    expect(result.files).toContain('typography.min.css');
    expect(result.css).toContain('.aksel-heading');
  });

  test('returns linkanchor CSS for HTML with aksel-link-anchor class', () => {
    const html = '<a class="aksel-link-anchor" href="#">Link</a>';
    const result = findComponentCss(html);

    expect(result.files).toContain('linkanchor.min.css');
    expect(result.css).toContain('.aksel-link-anchor');
  });

  test('returns multiple component CSS files when HTML uses classes from several components', () => {
    const html =
      '<div class="aksel-link-card"><h2 class="aksel-heading aksel-heading--small"><a class="aksel-link-anchor" href="#">Link</a></h2></div>';
    const result = findComponentCss(html);

    expect(result.files).toContain('linkcard.min.css');
    expect(result.files).toContain('typography.min.css');
    expect(result.files).toContain('linkanchor.min.css');
  });

  test('does not return unrelated component CSS files', () => {
    const html = '<div class="aksel-link-card">content</div>';
    const result = findComponentCss(html);

    expect(result.files).not.toContain('button.min.css');
    expect(result.files).not.toContain('accordion.min.css');
    expect(result.files).not.toContain('table.min.css');
    expect(result.files).not.toContain('alert.min.css');
  });

  test('returns empty results for HTML with no class attributes', () => {
    const html = '<div><span>no classes here</span></div>';
    const result = findComponentCss(html);

    expect(result.files).toHaveLength(0);
    expect(result.css).toBe('');
  });

  test('returns empty results for HTML with classes that match no component', () => {
    const html = '<div class="my-custom-class another-class">content</div>';
    const result = findComponentCss(html);

    expect(result.files).toHaveLength(0);
    expect(result.css).toBe('');
  });

  test('handles multiple HTML strings passed as separate arguments', () => {
    const html1 = '<div class="aksel-link-card">card</div>';
    const html2 = '<button class="aksel-button">click</button>';
    const result = findComponentCss(html1, html2);

    expect(result.files).toContain('linkcard.min.css');
    expect(result.files).toContain('button.min.css');
  });

  test('deduplicates CSS files when same classes appear in multiple HTML strings', () => {
    const html1 = '<div class="aksel-link-card">card 1</div>';
    const html2 = '<div class="aksel-link-card">card 2</div>';
    const result = findComponentCss(html1, html2);

    const linkcardCount = result.files.filter((f) => f === 'linkcard.min.css').length;
    expect(linkcardCount).toBe(1);
  });

  test('handles BEM modifier classes that reference their component file', () => {
    const html = '<div class="aksel-link-card aksel-link-card--small">content</div>';
    const result = findComponentCss(html);

    expect(result.files).toContain('linkcard.min.css');
    expect(result.css).toContain('.aksel-link-card--small');
  });

  test('handles BEM element classes that reference their component file', () => {
    const html = '<div class="aksel-link-card__icon">icon</div>';
    const result = findComponentCss(html);

    expect(result.files).toContain('linkcard.min.css');
  });

  test('returns files in sorted order', () => {
    const html = '<div class="aksel-heading aksel-link-card aksel-link-anchor">content</div>';
    const result = findComponentCss(html);

    const sorted = [...result.files].sort();
    expect(result.files).toEqual(sorted);
  });

  test('handles empty string input', () => {
    const result = findComponentCss('');

    expect(result.files).toHaveLength(0);
    expect(result.css).toBe('');
  });

  test('matches the exact three files needed for the actual microfrontend markup', () => {
    // This is the actual rendered HTML structure from the LinkCard component
    const html =
      '<div data-color="neutral" data-align-arrow="baseline" class="aksel-link-anchor__overlay aksel-link-card aksel-link-card--medium aksel-body-long aksel-body-long--medium">' +
      '<div aria-hidden="true" class="aksel-link-card__icon"><svg></svg></div>' +
      '<h2 class="aksel-link-card__title aksel-heading aksel-heading--small">' +
      '<a href="#" class="aksel-link-anchor">Title</a></h2>' +
      '<div class="aksel-link-card__description">Description</div>' +
      '<svg class="aksel-link-anchor__arrow aksel-link-card__arrow"></svg>' +
      '</div>';
    const result = findComponentCss(html);

    expect(result.files).toContain('linkanchor.min.css');
    expect(result.files).toContain('linkcard.min.css');
    expect(result.files).toContain('typography.min.css');
    expect(result.files).toHaveLength(3);
  });
});
