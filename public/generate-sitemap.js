const fs = require('fs');

const gamesData = require('./src/data/games.json'); 
const booksData = require('./src/data/books.json'); 

const DOMAIN = 'https://stashbook.online';

function generateSitemap() {
  let xml = `<?xml version="1.0" encoding="UTF-8"?>\n`;
  xml += `<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">\n`;

  const staticPages = [
    { path: '', freq: 'daily', priority: '1.0' },               
    { path: '/books', freq: 'daily', priority: '0.9' },         
    { path: '/iconic-items', freq: 'daily', priority: '0.9' }, 
    { path: '/about', freq: 'monthly', priority: '0.5' },     
    { path: '/contact', freq: 'monthly', priority: '0.5' }  
  ];

  staticPages.forEach(page => {
    xml += `  <url>\n`;
    xml += `    <loc>${DOMAIN}${page.path}</loc>\n`;
    xml += `    <changefreq>${page.freq}</changefreq>\n`;
    xml += `    <priority>${page.priority}</priority>\n`;
    xml += `  </url>\n`;
  });

  gamesData.forEach(game => {
    const encodedId = encodeURIComponent(game.id);
    xml += `  <url>\n`;
    xml += `    <loc>${DOMAIN}/game/${encodedId}</loc>\n`;
    xml += `    <changefreq>weekly</changefreq>\n`;
    xml += `    <priority>0.7</priority>\n`;
    xml += `  </url>\n`;
  });

  booksData.forEach(book => {
    const encodedId = encodeURIComponent(book.id);
    xml += `  <url>\n`;
    xml += `    <loc>${DOMAIN}/book/${encodedId}</loc>\n`;
    xml += `    <changefreq>monthly</changefreq>\n`;
    xml += `    <priority>0.6</priority>\n`;
    xml += `  </url>\n`;
  });

  xml += `</urlset>`;

  fs.writeFileSync('./public/sitemap.xml', xml);
  console.log('✅ Sitemap completely generated for stashbook.online!');
}

generateSitemap();
