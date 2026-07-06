import React, { useState, useEffect } from 'react';
import { Helmet } from 'react-helmet-async';
import GameCard from './GameCard';
import ServiceUnavailable from './ServiceUnavailable';

const DEFAULT_DESCRIPTION = 'StashBook is your daily loot drop for free game deals, developer books, and curated tech gear—save money and level up your setup without the grind.';
const DEFAULT_KEYWORDS = 'free games, game deals, stashbook, iconic drops, steam sales';

export default function IconicItems() {
  const [games, setGames] = useState([]);
  const [loading, setLoading] = useState(true);
  const [fetchError, setFetchError] = useState(false);
  const [selectedGame, setSelectedGame] = useState(null);

  useEffect(() => {
    fetch('/games.json')
      .then((res) => {
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        return res.json();
      })
      .then((data) => {
        const iconic = (data || []).filter(game => game.id.startsWith('legendary'));
        setGames(iconic);
        setLoading(false);
      })
      .catch((err) => {
        console.error("Error fetching data:", err);
        setFetchError(true);
        setLoading(false);
      });
  }, []);

  const closeModal = () => setSelectedGame(null);

  return (
    <div className="w-full flex-grow flex flex-col items-center pt-6 md:pt-10 pb-10 md:pb-16 px-2 md:px-6 overflow-x-hidden">
      <Helmet>
        <title>Iconic Drops – Legendary Free Games | StashBook</title>
        <meta name="description" content="Track iconic and legendary free game drops on StashBook. Limited-time offers on acclaimed titles from Steam, Epic, and other platforms." />
      </Helmet>
      <h2 className="text-lg sm:text-2xl md:text-3xl lg:text-5xl text-retro-orange tracking-[0.1em] md:tracking-[0.3em] uppercase drop-shadow-[3px_3px_0_rgba(255,107,0,0.5)] mb-8 md:mb-12 text-center break-words px-2">
        [ Iconic Drops ]
      </h2>

      {loading ? (
        <div className="text-center text-[10px] md:text-sm text-retro-gray uppercase animate-pulse">Scanning Vault...</div>
      ) : fetchError ? (
        <ServiceUnavailable />
      ) : games.length === 0 ? (
        <div className="text-center text-[10px] md:text-sm text-retro-gray uppercase mt-8 opacity-50">No iconic drops secured yet.</div>
      ) : (
        <div className="flex flex-wrap justify-center gap-4 w-full max-w-5xl mx-auto">
          {games.map(game => (
            <GameCard
              key={game.id}
              title={game.title}
              platform={game.platform}
              image={game.image}
              isExpired={game.is_expired}
              price={game.original_price} 
              onClick={() => setSelectedGame(game)}
            />
          ))}
        </div>
      )}

      {selectedGame && (
        <>
        <Helmet>
          <title>{selectedGame.title} - Stashbook Free Deals</title>
          <meta name="description" content={selectedGame.description?.trim() || DEFAULT_DESCRIPTION} />
          <meta name="keywords" content={selectedGame.seo_keywords?.trim() || DEFAULT_KEYWORDS} />
        </Helmet>
        <div className="fixed inset-0 z-[100] flex items-center justify-center p-2 md:p-4 bg-black/80 backdrop-blur-sm overflow-x-hidden">
          <div className="absolute inset-0 cursor-pointer" onClick={closeModal}></div>
          <div className="relative bg-black border-2 border-retro-orange p-4 md:p-8 w-full max-w-xl shadow-[8px_8px_0_0_var(--color-retro-orange)] animate-[fadeIn_0.2s_ease-out] overflow-hidden">
            <button onClick={closeModal} className="absolute top-2 right-3 md:top-3 md:right-4 text-retro-gray hover:text-retro-orange text-sm md:text-xl font-bold transition-colors z-10">[X]</button>
            <div className="flex flex-col gap-3 md:gap-5 mt-4">
              <div className="flex flex-col gap-2">
                <h2 className="text-sm md:text-2xl lg:text-3xl text-retro-white uppercase tracking-wider leading-tight break-words line-clamp-3 pr-6">{selectedGame.title}</h2>
                <div className="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-1 sm:gap-2 border-b border-retro-gray/30 pb-3">
                  <span className="text-[10px] md:text-sm text-retro-gray uppercase tracking-wider truncate">{selectedGame.platform}</span>
                  <span className="text-[10px] md:text-sm text-retro-orange animate-pulse font-bold tracking-widest shrink-0">
                    {selectedGame.is_expired ? 'OFFER EXPIRED' : (selectedGame.original_price > 0 ? `$${selectedGame.original_price}` : 'FREE')}
                  </span>
                </div>
              </div>
              <div className="py-2">
                <h4 className="text-[8px] md:text-[10px] text-retro-gray uppercase tracking-[0.2em] mb-2">Description</h4>
                <p className="text-[10px] md:text-xs lg:text-sm text-retro-white leading-relaxed max-h-48 overflow-y-auto pr-2 md:pr-3 custom-scrollbar break-words">
                  {selectedGame.description || "No description available for this drop."}
                </p>
              </div>
              <div className="flex flex-col gap-3 md:gap-4 mt-2 md:mt-4 pt-3 md:pt-4 border-t border-retro-gray/30">
                {!selectedGame.is_expired && (
                  <div className="text-[10px] md:text-xs text-retro-white uppercase tracking-wider text-center break-words">
                    Valid Until: <span className="text-retro-orange font-bold">{selectedGame.end_date || "Limited Time"}</span>
                  </div>
                )}
                <a href={selectedGame.game_url || "#"} target="_blank" rel="noopener noreferrer" className={`block text-center py-2 md:py-3 border-2 uppercase text-[10px] md:text-xs lg:text-sm tracking-[0.2em] md:tracking-[0.3em] font-bold transition-all duration-300 ${selectedGame.is_expired ? 'border-retro-gray/50 text-retro-gray/50 cursor-not-allowed pointer-events-none' : 'border-retro-orange text-retro-orange hover:bg-retro-orange hover:text-black hover:shadow-[0_0_15px_var(--color-retro-orange)]'}`}>
                  {selectedGame.is_expired ? 'Offer Ended' : 'Claim Drop'}
                </a>
              </div>
            </div>
          </div>
        </div>
        </>
      )}
    </div>
  );
}