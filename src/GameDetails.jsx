import React, { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import { Helmet } from 'react-helmet-async';
import { slugify } from './utils/stringUtils';

const DEFAULT_DESCRIPTION = 'StashBook is your daily loot drop for free game deals, developer books, and curated tech gear—save money and level up your setup without the grind.';
const DEFAULT_KEYWORDS = 'free games, game deals, stashbook, free steam games, epic games free';

export default function GameDetails() {
  const { slug } = useParams();
  const [game, setGame] = useState(null);
  const [loading, setLoading] = useState(true);
  const [notFound, setNotFound] = useState(false);

  useEffect(() => {
    setLoading(true);
    setNotFound(false);
    setGame(null);

    fetch('/games.json')
      .then((res) => {
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        return res.json();
      })
      .then((data) => {
        const match = data.find((g) => slugify(g.title) === slug);
        if (match) {
          setGame(match);
        } else {
          setNotFound(true);
        }
        setLoading(false);
      })
      .catch((err) => {
        console.error('Error fetching game:', err);
        setNotFound(true);
        setLoading(false);
      });
  }, [slug]);

  if (loading) {
    return (
      <main className="w-full flex-grow flex flex-col items-center justify-center min-h-[50vh] px-2 md:px-6">
        <div className="text-center text-retro-gray uppercase animate-pulse text-[10px] md:text-sm tracking-widest">
          Loading Loot...
        </div>
      </main>
    );
  }

  if (notFound || !game) {
    return (
      <main className="w-full flex-grow flex flex-col items-center justify-center min-h-[50vh] px-2 md:px-6">
        <Helmet>
          <title>404 – Loot Not Found | StashBook</title>
          <meta name="robots" content="noindex" />
        </Helmet>
        <div className="text-center">
          <h1 className="text-lg sm:text-2xl md:text-4xl text-retro-orange uppercase tracking-widest drop-shadow-[2px_2px_0_var(--color-retro-orange)] mb-4">
            ERROR 404: LOOT NOT FOUND
          </h1>
          <p className="text-[10px] md:text-sm text-retro-gray uppercase tracking-wider mb-8">
            This game deal may have expired or never existed.
          </p>
          <Link
            to="/"
            className="inline-block border-2 border-retro-orange text-retro-orange px-4 md:px-8 py-2 md:py-3 text-[10px] md:text-xs tracking-widest uppercase hover:bg-retro-orange hover:text-black transition-all duration-300"
          >
            [ Back to Loot ]
          </Link>
        </div>
      </main>
    );
  }

  const description = game.description?.trim() || DEFAULT_DESCRIPTION;
  const keywords = game.seo_keywords?.trim() || DEFAULT_KEYWORDS;
  const canonicalUrl = `https://stashbook.online/game/${slug}`;

  return (
    <main className="w-full flex-grow flex flex-col items-center justify-center py-8 md:py-16 px-2 md:px-6 overflow-x-hidden">
      <Helmet>
        <title>{game.title} - Free Download &amp; Details | StashBook</title>
        <meta name="description" content={description} />
        <meta name="keywords" content={keywords} />
        <link rel="canonical" href={canonicalUrl} />
      </Helmet>

      <div className="relative bg-black border-2 border-retro-orange p-4 md:p-8 w-full max-w-xl shadow-[8px_8px_0_0_var(--color-retro-orange)] animate-[fadeIn_0.2s_ease-out] overflow-hidden">
        <Link
          to="/"
          className="absolute top-2 right-3 md:top-3 md:right-4 text-retro-gray hover:text-retro-orange text-sm md:text-xl font-bold transition-colors z-10"
        >
          [X]
        </Link>
        <div className="flex flex-col gap-3 md:gap-5 mt-4">
          <div className="flex flex-col gap-2">
            <h1 className="text-sm md:text-2xl lg:text-3xl text-retro-white uppercase tracking-wider leading-tight break-words line-clamp-3 pr-6">
              {game.title}
            </h1>
            <div className="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-1 sm:gap-2 border-b border-retro-gray/30 pb-3">
              <span className="text-[10px] md:text-sm text-retro-gray uppercase tracking-wider truncate">
                {game.platform}
              </span>
              <span className="text-[10px] md:text-sm text-retro-orange animate-pulse font-bold tracking-widest shrink-0">
                {game.is_expired ? 'OFFER EXPIRED' : 'FREE (100% OFF)'}
              </span>
            </div>
          </div>
          <div className="py-2">
            <h2 className="text-[8px] md:text-[10px] text-retro-gray uppercase tracking-[0.2em] mb-2">
              Description
            </h2>
            <p className="text-[10px] md:text-xs lg:text-sm text-retro-white leading-relaxed max-h-48 overflow-y-auto pr-2 md:pr-3 custom-scrollbar break-words">
              {game.description || 'No description available for this game.'}
            </p>
          </div>
          <div className="flex flex-col gap-3 md:gap-4 mt-2 md:mt-4 pt-3 md:pt-4 border-t border-retro-gray/30">
            {!game.is_expired && (
              <div className="text-[10px] md:text-xs text-retro-white uppercase tracking-wider text-center break-words">
                Valid Until:{' '}
                <span className="text-retro-orange font-bold">{game.end_date || 'Limited Time'}</span>
              </div>
            )}
            <a
              href={game.game_url || '#'}
              target="_blank"
              rel="noopener noreferrer"
              className={`block text-center py-2 md:py-3 border-2 uppercase text-[10px] md:text-xs lg:text-sm tracking-[0.2em] md:tracking-[0.3em] font-bold transition-all duration-300 ${
                game.is_expired
                  ? 'border-retro-gray/50 text-retro-gray/50 cursor-not-allowed pointer-events-none'
                  : 'border-retro-orange text-retro-orange hover:bg-retro-orange hover:text-black hover:shadow-[0_0_15px_var(--color-retro-orange)]'
              }`}
            >
              {game.is_expired ? 'Offer Ended' : 'Get It Now'}
            </a>
          </div>
        </div>
      </div>
    </main>
  );
}
