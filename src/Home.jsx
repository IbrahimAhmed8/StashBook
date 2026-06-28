import React, { useState, useEffect } from 'react';
import { Helmet } from 'react-helmet-async';
import GameCard from './GameCard';
import ServiceUnavailable from './ServiceUnavailable';

const DEFAULT_DESCRIPTION = 'StashBook is your daily loot drop for free game deals, developer books, and curated tech gear—save money and level up your setup without the grind.';
const DEFAULT_KEYWORDS = 'free games, game deals, stashbook, free steam games, epic games free';

const groupGamesByPlatform = (gamesList) => {
  const grouped = {};
  gamesList.forEach(game => {
    const pName = game.platform.toLowerCase();
    let groupName = "Other Stores";
    
    if (pName.includes("epic")) groupName = "Epic Games";
    else if (pName.includes("steam")) groupName = "Steam";
    else if (pName.includes("prime") || pName.includes("amazon")) groupName = "Prime Gaming";
    else if (pName.includes("gog")) groupName = "GOG";
    else if (pName.includes("itch")) groupName = "itch.io";
    else if (pName.includes("vr")) groupName = "VR Games";
    else if (pName.includes("android") || pName.includes("ios")) groupName = "Mobile";
    else if (pName.includes("drm-free") || pName.includes("indiegala")) groupName = "PC / DRM-Free";
    else groupName = "PC"; 

    if (!grouped[groupName]) { grouped[groupName] = []; }
    grouped[groupName].push(game);
  });
  return grouped;
};

export default function Home() {
  const [games, setGames] = useState([]);
  const [loading, setLoading] = useState(true);
  const [fetchError, setFetchError] = useState(false);
  const [view, setView] = useState('active'); 
  const [selectedGame, setSelectedGame] = useState(null);

  useEffect(() => {
    fetch('/games.json')
      .then((res) => {
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        return res.json();
      })
      .then((data) => {
        const freeOnly = (data || []).filter(game => !game.id.startsWith('legendary'));
        setGames(freeOnly);
        setLoading(false);
      })
      .catch((err) => {
        console.error("Error fetching games:", err);
        setFetchError(true);
        setLoading(false);
      });
  }, []);

  const activeGames = games.filter(game => !game.is_expired);
  const expiredGames = games.filter(game => game.is_expired); 
  const closeModal = () => setSelectedGame(null);

  const renderGroupedSections = (gamesList, isExpiredView) => {
    if (gamesList.length === 0) {
      return (
        <div className="text-center text-[10px] md:text-sm text-retro-gray uppercase mt-8 opacity-50 min-h-[30vh] px-2 break-words">
          {isExpiredView ? 'Vault is empty. No expired games tracked yet.' : 'No free games found right now.'}
        </div>
      );
    }

    const groupedGames = groupGamesByPlatform(gamesList);

    return Object.entries(groupedGames).map(([platformName, platformGames]) => (
      <div key={platformName} className="mb-10 md:mb-16 w-full flex flex-col items-center">
        <div className="text-center w-full flex justify-center items-center mb-6 md:mb-8 px-2 md:px-6">
          <h3 className="text-[10px] sm:text-sm md:text-xl lg:text-2xl text-retro-white tracking-[0.1em] md:tracking-[0.3em] uppercase drop-shadow-[2px_2px_0_rgba(255,107,0,0.5)] break-words line-clamp-2">
            {platformName}
          </h3>
        </div>
        
        <div className="flex flex-wrap justify-center gap-4 w-full max-w-5xl mx-auto px-2 md:px-6">
          {platformGames.map((game) => (
            <GameCard 
              key={game.id}
              title={game.title}
              platform={platformName} 
              image={game.image}
              isExpired={isExpiredView}
              onClick={() => setSelectedGame(game)} 
            />
          ))}
        </div>
      </div>
    ));
  };

  const scrollToLoot = () => {
    const element = document.getElementById('loot-section');
    if (element) {
      element.scrollIntoView({ behavior: 'smooth' });
    }
  };

  return (
    <div className="w-full overflow-x-hidden">
      <section className="w-full min-h-[40vh] md:min-h-[60vh] flex flex-col items-center justify-center text-center px-2 md:px-6 border-b border-retro-gray/20 mb-8 md:mb-16 relative overflow-hidden">
        <div className="absolute inset-0 bg-[linear-gradient(rgba(255,255,255,0.03)_1px,transparent_1px),linear-gradient(90deg,rgba(255,255,255,0.03)_1px,transparent_1px)] bg-[size:40px_40px] pointer-events-none"></div>
        <h1 className="text-2xl sm:text-3xl md:text-5xl lg:text-7xl xl:text-8xl tracking-[0.05em] sm:tracking-[0.1em] md:tracking-[0.2em] uppercase text-retro-white drop-shadow-[2px_2px_0_var(--color-retro-orange)] md:drop-shadow-[4px_4px_0_var(--color-retro-orange)] z-10 mb-4 md:mb-6 break-words max-w-full px-2">
          StashBook
        </h1>
        <button onClick={scrollToLoot} className="z-10 border-2 border-retro-orange text-retro-orange px-4 md:px-8 py-2 md:py-3 text-[10px] md:text-xs lg:text-sm tracking-widest uppercase hover:bg-retro-orange hover:text-black transition-all duration-300 shadow-[0_0_15px_rgba(255,107,0,0.2)] hover:shadow-[0_0_25px_rgba(255,107,0,0.6)]">
          [ Access Loot ]
        </button>
      </section>

      <div id="loot-section" className="flex flex-col w-full items-center pt-2 md:pt-4 flex-grow mb-6 md:mb-10 overflow-x-hidden">
        <div className="flex flex-col sm:flex-row gap-3 sm:gap-8 mb-8 md:mb-12 border-b border-retro-gray/50 pb-2 px-2 md:px-6 items-center justify-center">
          <button onClick={() => setView('active')} className={`text-[10px] md:text-sm tracking-widest uppercase transition-colors duration-300 ${view === 'active' ? 'text-retro-orange drop-shadow-[0_0_8px_rgba(255,107,0,0.5)]' : 'text-retro-gray hover:text-retro-white'}`}>
            Active Offers
          </button>
          <button onClick={() => setView('expired')} className={`text-[10px] md:text-sm tracking-widest uppercase transition-colors duration-300 ${view === 'expired' ? 'text-retro-white drop-shadow-[0_0_8px_rgba(255,255,255,0.5)]' : 'text-retro-gray hover:text-retro-white'}`}>
            Expired Vault
          </button>
        </div>

        <section className="w-full flex flex-col items-center">
          {loading ? (
            <div className="text-center text-retro-gray uppercase animate-pulse">Scanning Servers...</div>
          ) : fetchError ? (
            <ServiceUnavailable />
          ) : view === 'active' ? (
            renderGroupedSections(activeGames, false)
          ) : (
            renderGroupedSections(expiredGames, true)
          )}
        </section>
      </div>

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
                  <span className="text-[10px] md:text-sm text-retro-orange animate-pulse font-bold tracking-widest shrink-0">{selectedGame.is_expired ? 'OFFER EXPIRED' : 'FREE (100% OFF)'}</span>
                </div>
              </div>
              <div className="py-2">
                <h4 className="text-[8px] md:text-[10px] text-retro-gray uppercase tracking-[0.2em] mb-2">Description</h4>
                <p className="text-[10px] md:text-xs lg:text-sm text-retro-white leading-relaxed max-h-48 overflow-y-auto pr-2 md:pr-3 custom-scrollbar break-words">{selectedGame.description || "No description available for this game."}</p>
              </div>
              <div className="flex flex-col gap-3 md:gap-4 mt-2 md:mt-4 pt-3 md:pt-4 border-t border-retro-gray/30">
                {!selectedGame.is_expired && (
                  <div className="text-[10px] md:text-xs text-retro-white uppercase tracking-wider text-center break-words">
                    Valid Until: <span className="text-retro-orange font-bold">{selectedGame.end_date || "Limited Time"}</span>
                  </div>
                )}
                <a href={selectedGame.game_url || "#"} target="_blank" rel="noopener noreferrer" className={`block text-center py-2 md:py-3 border-2 uppercase text-[10px] md:text-xs lg:text-sm tracking-[0.2em] md:tracking-[0.3em] font-bold transition-all duration-300 ${selectedGame.is_expired ? 'border-retro-gray/50 text-retro-gray/50 cursor-not-allowed pointer-events-none' : 'border-retro-orange text-retro-orange hover:bg-retro-orange hover:text-black hover:shadow-[0_0_15px_var(--color-retro-orange)]'}`}>
                  {selectedGame.is_expired ? 'Offer Ended' : 'Get It Now'}
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