/**
 * Post-build static HTML generator for SEO / AI crawlers.
 * Runs after `vite build`. Writes crawlable HTML for routes while keeping the SPA bundle.
 */
import fs from 'node:fs/promises';
import path from 'node:path';
import { fileURLToPath, pathToFileURL } from 'node:url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const ROOT = path.resolve(__dirname, '..');
const DIST = path.join(ROOT, 'dist');
const SITE = 'https://stashbook.online';
const DEFAULT_IMAGE = `${SITE}/favicon.png`;

const { slugify } = await import(
  pathToFileURL(path.join(ROOT, 'src/utils/stringUtils.js')).href
);
const {
  filterActiveFreeGames,
  filterIconicDeals,
} = await import(
  pathToFileURL(path.join(ROOT, 'src/utils/gameFilters.js')).href
);

function escapeHtml(value = '') {
  return String(value)
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#39;');
}

function stripExistingSeoTags(html) {
  return html
    .replace(/<title>[\s\S]*?<\/title>/i, '')
    .replace(/<meta\s+name="description"[^>]*>/gi, '')
    .replace(/<meta\s+name="keywords"[^>]*>/gi, '')
    .replace(/<meta\s+name="twitter:[^"]+"[^>]*>/gi, '')
    .replace(/<link\s+rel="canonical"[^>]*>/gi, '')
    .replace(/<meta\s+property="og:[^"]+"[^>]*>/gi, '')
    .replace(
      /<script\s+type="application\/ld\+json">[\s\S]*?<\/script>/gi,
      ''
    );
}

function buildHead({ title, description, url, image, keywords, jsonLd }) {
  const safeTitle = escapeHtml(title);
  const safeDesc = escapeHtml(description);
  const safeImage = escapeHtml(image || DEFAULT_IMAGE);
  const safeUrl = escapeHtml(url);
  const safeKeywords = keywords ? escapeHtml(keywords) : '';

  return `
    <title>${safeTitle}</title>
    <meta name="description" content="${safeDesc}" />
    ${safeKeywords ? `<meta name="keywords" content="${safeKeywords}" />` : ''}
    <link rel="canonical" href="${safeUrl}" />
    <meta property="og:title" content="${safeTitle}" />
    <meta property="og:description" content="${safeDesc}" />
    <meta property="og:type" content="website" />
    <meta property="og:url" content="${safeUrl}" />
    <meta property="og:image" content="${safeImage}" />
    <meta name="twitter:card" content="summary_large_image" />
    <meta name="twitter:title" content="${safeTitle}" />
    <meta name="twitter:description" content="${safeDesc}" />
    <meta name="twitter:image" content="${safeImage}" />
    <script type="application/ld+json">${JSON.stringify(jsonLd)}</script>
  `;
}

function injectPage(templateHtml, { headInner, bodyInner }) {
  let html = stripExistingSeoTags(templateHtml);
  html = html.replace('</head>', `${headInner}\n  </head>`);
  html = html.replace(
    /<div id="root"><\/div>/,
    `<div id="root">${bodyInner}</div>`
  );
  return html;
}

async function writeRoute(relDir, html) {
  const outDir = relDir ? path.join(DIST, relDir) : DIST;
  await fs.mkdir(outDir, { recursive: true });
  await fs.writeFile(path.join(outDir, 'index.html'), html, 'utf8');
}

function gameStatusLabel(game) {
  if (game.is_expired) return 'Expired';
  const isIconic = game.platform?.toLowerCase().includes('iconic');
  if (isIconic) return 'Iconic discount / deal';
  return 'Free / Active';
}

function buildGameJsonLd(game, slug, url, description) {
  const offer = {
    '@type': 'Offer',
    url: game.game_url || url,
    priceCurrency: 'USD',
    price: game.is_expired ? undefined : 0,
    availability: game.is_expired
      ? 'https://schema.org/OutOfStock'
      : 'https://schema.org/InStock',
    seller: {
      '@type': 'Organization',
      name: game.platform || 'Unknown',
    },
  };

  if (
    !game.is_expired &&
    game.end_date &&
    game.end_date !== 'Unknown' &&
    game.end_date !== 'Limited Time'
  ) {
    offer.priceValidUntil = game.end_date;
  }

  // Drop undefined price for expired offers (cleaner JSON-LD)
  if (offer.price === undefined) {
    delete offer.price;
  }

  return {
    '@context': 'https://schema.org',
    '@graph': [
      {
        '@type': 'WebPage',
        '@id': `${url}#webpage`,
        url,
        name: `${game.title} - Free Download & Details | StashBook`,
        description,
        isPartOf: { '@id': `${SITE}/#website` },
      },
      {
        '@type': 'Product',
        name: game.title,
        description,
        image: game.image || undefined,
        category: game.platform || undefined,
        offers: offer,
      },
    ],
  };
}

function renderGameBody(game) {
  const status = gameStatusLabel(game);
  const endDate =
    !game.is_expired && game.end_date
      ? `<p>Valid until: ${escapeHtml(game.end_date)}</p>`
      : '';
  const price =
    game.original_price != null
      ? `<p>Original price: $${escapeHtml(game.original_price)}</p>`
      : '';
  const link = game.game_url
    ? `<p><a href="${escapeHtml(game.game_url)}">Get it now</a></p>`
    : '';

  return `
    <article>
      <h1>${escapeHtml(game.title)}</h1>
      <p>Platform: ${escapeHtml(game.platform || 'Unknown')}</p>
      <p>Status: ${escapeHtml(status)}</p>
      ${endDate}
      <p>${escapeHtml(game.description || 'No description available for this game.')}</p>
      ${price}
      ${link}
      <p><a href="/">Back to StashBook</a></p>
    </article>
  `;
}

function renderHomeBody(games) {
  const active = filterActiveFreeGames(games);
  const listItems = active
    .map((game) => {
      const slug = slugify(game.title);
      if (!slug) return '';
      return `<li><a href="/game/${slug}">${escapeHtml(game.title)}</a> — ${escapeHtml(game.platform || '')}</li>`;
    })
    .filter(Boolean)
    .join('\n');

  return `
    <header>
      <h1>StashBook</h1>
      <p>Browse live free game offers from Steam, Epic Games, GOG, Prime Gaming, and itch.io. StashBook scans storefronts daily so you never miss a giveaway.</p>
    </header>
    <nav>
      <ul>
        <li><a href="/">Games</a></li>
        <li><a href="/books">Books</a></li>
        <li><a href="/iconic">Iconic Items</a></li>
        <li><a href="/about">About</a></li>
        <li><a href="/contact">Contact</a></li>
      </ul>
    </nav>
    <section>
      <h2>Active free games</h2>
      <ul>
        ${listItems || '<li>No free games found right now.</li>'}
      </ul>
    </section>
  `;
}

function renderBooksBody(books) {
  const listItems = (books || [])
    .map(
      (book) =>
        `<li><strong>${escapeHtml(book.title)}</strong>${book.author ? ` — ${escapeHtml(book.author)}` : ''}${book.description ? `<p>${escapeHtml(book.description)}</p>` : ''}${book.book_url ? `<p><a href="${escapeHtml(book.book_url)}">Download</a></p>` : ''}</li>`
    )
    .join('\n');

  return `
    <header>
      <h1>Tech Archives – Free Developer Books</h1>
      <p>Discover free programming and developer books curated by StashBook.</p>
    </header>
    <ul>
      ${listItems || '<li>No books available right now.</li>'}
    </ul>
    <p><a href="/">Back to StashBook</a></p>
  `;
}

function renderIconicBody(games) {
  const iconic = filterIconicDeals(games);
  const listItems = iconic
    .map((game) => {
      const slug = slugify(game.title);
      if (!slug) return '';
      return `<li><a href="/game/${slug}">${escapeHtml(game.title)}</a> — ${escapeHtml(game.platform || '')}${game.description ? `: ${escapeHtml(game.description)}` : ''}</li>`;
    })
    .filter(Boolean)
    .join('\n');

  return `
    <header>
      <h1>Iconic Drops – Legendary Free Games</h1>
      <p>Track iconic and legendary free game drops and discounts on StashBook.</p>
    </header>
    <ul>
      ${listItems || '<li>No iconic drops secured yet.</li>'}
    </ul>
    <p><a href="/">Back to StashBook</a></p>
  `;
}

async function main() {
  const templatePath = path.join(DIST, 'index.html');
  const gamesPath = path.join(DIST, 'games.json');
  const booksPath = path.join(DIST, 'books.json');

  const template = await fs.readFile(templatePath, 'utf8');
  const games = JSON.parse(await fs.readFile(gamesPath, 'utf8'));

  let books = [];
  try {
    books = JSON.parse(await fs.readFile(booksPath, 'utf8'));
  } catch {
    // books.json is optional for prerender
  }

  const homeDescription =
    'Browse live free game offers from Steam, Epic Games, GOG, Prime Gaming, and itch.io. StashBook scans storefronts daily so you never miss a giveaway.';

  const homeHtml = injectPage(template, {
    headInner: buildHead({
      title: 'StashBook – Free Game Deals & Giveaways',
      description: homeDescription,
      url: `${SITE}/`,
      jsonLd: {
        '@context': 'https://schema.org',
        '@type': 'WebSite',
        '@id': `${SITE}/#website`,
        name: 'StashBook',
        url: SITE,
        description: homeDescription,
      },
    }),
    bodyInner: renderHomeBody(games),
  });
  await writeRoute('', homeHtml);

  const staticPages = [
    {
      dir: 'books',
      title: 'Tech Archives – Free Developer Books | StashBook',
      description:
        'Discover free programming and developer books curated by StashBook. Download PDFs on software engineering, game dev, and tech topics at no cost.',
      body: renderBooksBody(books),
      jsonLd: {
        '@context': 'https://schema.org',
        '@type': 'CollectionPage',
        name: 'Tech Archives – Free Developer Books',
        url: `${SITE}/books`,
        isPartOf: { '@id': `${SITE}/#website` },
      },
    },
    {
      dir: 'iconic',
      title: 'Iconic Drops – Legendary Free Games | StashBook',
      description:
        'Track iconic and legendary free game drops on StashBook. Limited-time offers on acclaimed titles from Steam, Epic, and other platforms.',
      body: renderIconicBody(games),
      jsonLd: {
        '@context': 'https://schema.org',
        '@type': 'CollectionPage',
        name: 'Iconic Drops – Legendary Free Games',
        url: `${SITE}/iconic`,
        isPartOf: { '@id': `${SITE}/#website` },
      },
    },
    {
      dir: 'about',
      title: 'About StashBook – How We Hunt Free Loot',
      description:
        'Learn how StashBook automates free game and digital asset discovery. Independent tool, not affiliated with Steam, Epic, or other storefronts.',
      body: `
        <article>
          <h1>About StashBook</h1>
          <p>StashBook is an automated aggregation engine designed to hunt down 100% free digital assets across the web. We track limited-time giveaways and digital loot so you don't have to manually check dozens of storefronts.</p>
          <p>StashBook is NOT affiliated with, endorsed by, or connected to Epic Games, Valve (Steam), Amazon (Prime Gaming), GOG, itch.io, or any other corporation mentioned on this platform.</p>
          <p><a href="/">Back to StashBook</a></p>
        </article>
      `,
      jsonLd: {
        '@context': 'https://schema.org',
        '@type': 'AboutPage',
        name: 'About StashBook',
        url: `${SITE}/about`,
        isPartOf: { '@id': `${SITE}/#website` },
      },
    },
    {
      dir: 'contact',
      title: 'Contact Us - StashBook',
      description:
        'Send feedback, report a broken link, or get in touch with the StashBook team.',
      body: `
        <article>
          <h1>Contact Us</h1>
          <p>Send feedback, report a broken link, or get in touch with the StashBook team.</p>
          <p><a href="/">Back to StashBook</a></p>
        </article>
      `,
      jsonLd: {
        '@context': 'https://schema.org',
        '@type': 'ContactPage',
        name: 'Contact Us - StashBook',
        url: `${SITE}/contact`,
        isPartOf: { '@id': `${SITE}/#website` },
      },
    },
  ];

  for (const page of staticPages) {
    const url = `${SITE}/${page.dir}`;
    const html = injectPage(template, {
      headInner: buildHead({
        title: page.title,
        description: page.description,
        url,
        jsonLd: page.jsonLd,
      }),
      bodyInner: page.body,
    });
    await writeRoute(page.dir, html);
  }

  let gamePages = 0;
  const seenSlugs = new Set();

  for (const game of games) {
    const slug = slugify(game.title);
    if (!slug || seenSlugs.has(slug)) continue;
    seenSlugs.add(slug);

    const url = `${SITE}/game/${slug}`;
    const description =
      game.description?.trim() ||
      `${game.title} free game details on StashBook.`;
    const keywords =
      game.seo_keywords?.trim() ||
      'free games, game deals, stashbook, free steam games, epic games free';

    const html = injectPage(template, {
      headInner: buildHead({
        title: `${game.title} - Free Download & Details | StashBook`,
        description,
        url,
        image: game.image || DEFAULT_IMAGE,
        keywords,
        jsonLd: buildGameJsonLd(game, slug, url, description),
      }),
      bodyInner: renderGameBody(game),
    });

    await writeRoute(path.join('game', slug), html);
    gamePages += 1;
  }

  console.log(
    `Prerender complete: homepage, ${staticPages.length} static pages, ${gamePages} game pages`
  );
}

main().catch((err) => {
  console.error('Prerender failed:', err);
  process.exit(1);
});
