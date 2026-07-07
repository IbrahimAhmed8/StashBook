/**
 * Converts text to a URL-safe slug.
 * Lowercase, alphanumeric segments separated by hyphens.
 */
export function slugify(text) {
  if (text == null || typeof text !== 'string') return '';

  return text
    .toLowerCase()
    .trim()
    .normalize('NFD')
    .replace(/[\u0300-\u036f]/g, '')
    .replace(/[^a-z0-9\s-]/g, '')
    .replace(/[\s_]+/g, '-')
    .replace(/-+/g, '-')
    .replace(/^-+|-+$/g, '');
}
