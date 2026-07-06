const PARTIAL_DISCOUNT_PATTERN = /^\[\s*\d+(?:\.\d+)?% OFF \]/i;

export function unescapeHtmlEntities(text) {
  if (typeof text !== 'string' || text === '') {
    return text ?? '';
  }

  return text
    .replace(/&amp;/g, '&')
    .replace(/&quot;/g, '"')
    .replace(/&#39;/g, "'")
    .replace(/&apos;/g, "'")
    .replace(/&lt;/g, '<')
    .replace(/&gt;/g, '>');
}

export function isLegendaryId(game) {
  return Boolean(game?.id?.startsWith('legendary'));
}

export function isPartialDiscount(game) {
  const description = game?.description?.trim() ?? '';
  return PARTIAL_DISCOUNT_PATTERN.test(description);
}

export function isIconicDeal(game) {
  return isLegendaryId(game) || isPartialDiscount(game);
}

export function isFullyFreeDeal(game) {
  return !isIconicDeal(game);
}

export function filterMainFeedGames(games) {
  return (games ?? []).filter(isFullyFreeDeal);
}

export function filterActiveFreeGames(games) {
  return filterMainFeedGames(games).filter((game) => game.is_expired === false);
}

export function filterExpiredGames(games) {
  return filterMainFeedGames(games).filter((game) => game.is_expired === true);
}

export function filterIconicDeals(games) {
  return (games ?? []).filter(isIconicDeal);
}

export function validateGamePartitions(games) {
  const active = filterActiveFreeGames(games);
  const expired = filterExpiredGames(games);
  const iconic = filterIconicDeals(games);

  const activeIds = new Set(active.map((game) => game.id));
  const expiredIds = new Set(expired.map((game) => game.id));
  const iconicIds = new Set(iconic.map((game) => game.id));

  const overlap = (left, right) => [...left].filter((id) => right.has(id));

  return {
    activeCount: active.length,
    expiredCount: expired.length,
    iconicCount: iconic.length,
    activeExpiredOverlap: overlap(activeIds, expiredIds),
    activeIconicOverlap: overlap(activeIds, iconicIds),
    expiredIconicOverlap: overlap(expiredIds, iconicIds),
    isValid:
      overlap(activeIds, expiredIds).length === 0 &&
      overlap(activeIds, iconicIds).length === 0 &&
      overlap(expiredIds, iconicIds).length === 0,
  };
}
