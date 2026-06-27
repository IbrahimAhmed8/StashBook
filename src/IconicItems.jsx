import React, { useState, useEffect } from 'react';
import GameCard from './GameCard';
import ServiceUnavailable from './ServiceUnavailable';

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
    <div className="w-full flex-grow flex flex-col items-center pt-10 pb-16 px-6">
      <h2 className="text-3xl md:text-5xl text-retro-orange tracking-[0.3em] uppercase drop-shadow-[3px_3px_0_rgba(255,107,0,0.5)] mb-12 text-center whitespace-nowrap">
        [ Iconic Drops ]
      </h2>

      {loading ? (
        <div className="text-center text-retro-gray uppercase animate-pulse">Scanning Vault...</div>
      ) : fetchError ? (
        <ServiceUnavailable />
      ) : games.length === 0 ? (
        <div className="text-center text-retro-gray uppercase mt-8 opacity-50">No iconic drops secured yet.</div>
      ) : (
        <div className="flex flex-wrap justify-center gap-4 md:gap-6 w-full max-w-5xl mx-auto">
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
        <div className="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-black/80 backdrop-blur-sm">
          <div className="absolute inset-0 cursor-pointer" onClick={closeModal}></div>
          <div className="relative bg-black border-2 border-retro-orange p-8 w-full max-w-xl shadow-[8px_8px_0_0_var(--color-retro-orange)] animate-[fadeIn_0.2s_ease-out]">
            <button onClick={closeModal} className="absolute top-3 right-4 text-retro-gray hover:text-retro-orange text-xl font-bold transition-colors z-10">[X]</button>
            <div className="flex flex-col gap-5 mt-4">
              <div className="flex flex-col gap-2">
                <h2 className="text-2xl md:text-3xl text-retro-white uppercase tracking-wider leading-tight">{selectedGame.title}</h2>
                <div className="flex justify-between items-center border-b border-retro-gray/30 pb-3">
                  <span className="text-sm text-retro-gray uppercase tracking-wider">{selectedGame.platform}</span>
                  <span className="text-sm text-retro-orange animate-pulse font-bold tracking-widest">
                    {selectedGame.is_expired ? 'OFFER EXPIRED' : (selectedGame.original_price > 0 ? `$${selectedGame.original_price}` : 'FREE')}
                  </span>
                </div>
              </div>
              <div className="py-2">
                <h4 className="text-[10px] text-retro-gray uppercase tracking-[0.2em] mb-2">Description</h4>
                <p className="text-xs md:text-sm text-retro-white leading-relaxed max-h-48 overflow-y-auto pr-3 custom-scrollbar">
                  {selectedGame.description || "No description available for this drop."}
                </p>
              </div>
              <div className="flex flex-col gap-4 mt-4 pt-4 border-t border-retro-gray/30">
                {!selectedGame.is_expired && (
                  <div className="text-xs text-retro-white uppercase tracking-wider text-center">
                    Valid Until: <span className="text-retro-orange font-bold">{selectedGame.end_date || "Limited Time"}</span>
                  </div>
                )}
                <a href={selectedGame.game_url || "#"} target="_blank" rel="noopener noreferrer" className={`block text-center py-3 border-2 uppercase text-xs md:text-sm tracking-[0.3em] font-bold transition-all duration-300 ${selectedGame.is_expired ? 'border-retro-gray/50 text-retro-gray/50 cursor-not-allowed pointer-events-none' : 'border-retro-orange text-retro-orange hover:bg-retro-orange hover:text-black hover:shadow-[0_0_15px_var(--color-retro-orange)]'}`}>
                  {selectedGame.is_expired ? 'Offer Ended' : 'Claim Drop'}
                </a>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}