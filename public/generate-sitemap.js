const fs = require('fs');
const path = require('path');

const gamesData = require(path.join(__dirname, 'public', 'games.json'));
const booksData = require(path.join(__dirname, 'public', 'books.json'));

const DOMAIN = 'https://stashbook.online';

function generateSitemap() {
  let xml = `<?xml version="1.0" encoding="UTF-8"?>\n`;
  xml += `<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">\n`;

  const staticPages = [
    { path: '/', freq: 'daily', priority: '1.0' },
    { path: '/books', freq: 'daily', priority: '0.9' },
    { path: '/iconic-items', freq: 'daily', priority: '0.9' },
    { path: '/about', freq: 'monthly', priority: '0.5' },
    { path: '/contact', freq: 'monthly', priority: '0.5' }
  ];

  staticPages.forEach(page => {
    xml += `  <url>\n`;
    xml += `    <loc>${DOMAIN}${page.path === '/' ? '' : page.path}</loc>\n`;
    xml += `    <changefreq>${page.freq}</changefreq>\n`;
    xml += `    <priority>${page.priority}</priority>\n`;
    xml += `  </url>\n`;
  });

  gamesData.forEach(game => {
    xml += `  <url>\n`;
    xml += `    <loc>${DOMAIN}/game/${encodeURIComponent(game.id)}</loc>\n`;
    xml += `    <changefreq>weekly</changefreq>\n`;
    xml += `    <priority>0.7</priority>\n`;
    xml += `  </url>\n`;
  });

  booksData.forEach(book => {
    xml += `  <url>\n`;
    xml += `    <loc>${DOMAIN}/book/${encodeURIComponent(book.id)}</loc>\n`;
    xml += `    <changefreq>monthly</changefreq>\n`;
    xml += `    <priority>0.6</priority>\n`;
    xml += `  </url>\n`;
  });

  xml += `</urlset>`;

  const outputPath = path.join(__dirname, 'public', 'sitemap.xml');
  fs.writeFileSync(outputPath, xml);
  console.log(`✅ Sitemap successfully generated at: ${outputPath}`);
}

generateSitemap();
