import { describe, expect, it } from 'vitest';
import {
  filterActiveFreeGames,
  filterExpiredGames,
  filterIconicDeals,
  filterMainFeedGames,
  unescapeHtmlEntities,
  validateGamePartitions,
} from './gameFilters';

const activeFree = {
  id: 'steam_live_123',
  title: 'Free Game',
  is_expired: false,
  description: 'Free game on Steam.',
  original_price: 19.99,
};

const expiredFree = {
  id: 'epic_456',
  title: 'Expired Game',
  is_expired: true,
  description: 'Free game on the Epic Games Store.',
  original_price: 9.99,
};

const legendaryDeal = {
  id: 'legendary_steam_789',
  title: 'Iconic Steam Deal',
  is_expired: false,
  description: '[ 75% OFF ] A legendary title on sale.',
  original_price: 14.99,
};

const partialDiscount = {
  id: 'cs_partial_001',
  title: 'Partial Discount',
  is_expired: false,
  description: '[ 50% OFF ] Still costs money.',
  original_price: 24.99,
};

const fixtureGames = [activeFree, expiredFree, legendaryDeal, partialDiscount];

describe('unescapeHtmlEntities', () => {
  it('decodes common HTML entities', () => {
    expect(unescapeHtmlEntities('Tom &amp; Jerry')).toBe('Tom & Jerry');
    expect(unescapeHtmlEntities('&quot;StashBook&quot;')).toBe('"StashBook"');
    expect(unescapeHtmlEntities('It&#39;s free')).toBe("It's free");
    expect(unescapeHtmlEntities('&lt;tag&gt;')).toBe('<tag>');
  });

  it('returns empty values unchanged', () => {
    expect(unescapeHtmlEntities('')).toBe('');
    expect(unescapeHtmlEntities(null)).toBe('');
    expect(unescapeHtmlEntities(undefined)).toBe('');
  });
});

describe('game feed filters', () => {
  it('keeps only 100% free non-iconic games in the main feed', () => {
    expect(filterMainFeedGames(fixtureGames)).toEqual([activeFree, expiredFree]);
  });

  it('routes active offers to the active feed only', () => {
    expect(filterActiveFreeGames(fixtureGames)).toEqual([activeFree]);
  });

  it('routes expired offers to the expired vault only', () => {
    expect(filterExpiredGames(fixtureGames)).toEqual([expiredFree]);
  });

  it('isolates legendary and partial discount deals in iconic feed', () => {
    expect(filterIconicDeals(fixtureGames)).toEqual([legendaryDeal, partialDiscount]);
  });

  it('validates partitions have no overlap', () => {
    const result = validateGamePartitions(fixtureGames);

    expect(result.isValid).toBe(true);
    expect(result.activeCount).toBe(1);
    expect(result.expiredCount).toBe(1);
    expect(result.iconicCount).toBe(2);
    expect(result.activeExpiredOverlap).toEqual([]);
    expect(result.activeIconicOverlap).toEqual([]);
    expect(result.expiredIconicOverlap).toEqual([]);
  });
});
