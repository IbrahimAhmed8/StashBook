import React, { useState, useMemo } from 'react';
import { useNavigate } from 'react-router-dom';
import { Helmet } from 'react-helmet-async';
import GameCard from './GameCard';
import ServiceUnavailable from './ServiceUnavailable';
import { filterActiveFreeGames, filterExpiredGames, filterMainFeedGames } from './utils/gameFilters';
import { slugify } from './utils/stringUtils';
import { useStashData } from './hooks/useStashData';

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
  const navigate = useNavigate();
  const { data, loading, error: fetchError } = useStashData('/games.json');
  const games = useMemo(() => (data ? filterMainFeedGames(data) : []), [data]);
  const [view, setView] = useState('active');

  const activeGames = filterActiveFreeGames(games);
  const expiredGames = filterExpiredGames(games);

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
          <h2 className="text-[10px] sm:text-sm md:text-xl lg:text-2xl text-retro-white tracking-[0.1em] md:tracking-[0.3em] uppercase drop-shadow-[2px_2px_0_rgba(255,107,0,0.5)] break-words line-clamp-2">
            {platformName}
          </h2>
        </div>
        
        <div className="flex flex-wrap justify-center gap-4 w-full max-w-5xl mx-auto px-2 md:px-6">
          {platformGames.map((game) => (
            <GameCard 
              key={game.id}
              title={game.title}
              platform={platformName} 
              image={game.image}
              isExpired={isExpiredView}
              description={game.description}
              onClick={() => navigate(`/game/${slugify(game.title)}`)}
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
      <Helmet>
        <title>StashBook – Free Game Deals &amp; Giveaways</title>
        <meta name="description" content="Browse live free game offers from Steam, Epic Games, GOG, Prime Gaming, and itch.io. StashBook scans storefronts daily so you never miss a giveaway." />
      </Helmet>
      <section className="w-full min-h-[40vh] md:min-h-[60vh] flex flex-col items-center justify-center text-center px-2 md:px-6 border-b border-retro-gray/20 mb-8 md:mb-16 relative overflow-hidden">
        <div className="absolute inset-0 bg-[linear-gradient(rgba(255,255,255,0.03)_1px,transparent_1px),linear-gradient(90deg,rgba(255,255,255,0.03)_1px,transparent_1px)] bg-[size:40px_40px] pointer-events-none"></div>
        <h1 className="text-2xl sm:text-3xl md:text-5xl lg:text-7xl xl:text-8xl tracking-[0.05em] sm:tracking-[0.1em] md:tracking-[0.2em] uppercase text-retro-white drop-shadow-[2px_2px_0_var(--color-retro-orange)] md:drop-shadow-[4px_4px_0_var(--color-retro-orange)] z-10 mb-4 md:mb-6 break-words max-w-full px-2">
          StashBook
        </h1>
        <button onClick={scrollToLoot} className="z-10 border-2 border-retro-orange text-retro-orange px-4 md:px-8 py-2 md:py-3 text-[10px] md:text-xs lg:text-sm tracking-widest uppercase hover:bg-retro-orange hover:text-black transition-all duration-300 shadow-[0_0_15px_rgba(255,107,0,0.2)] hover:shadow-[0_0_25px_rgba(255,107,0,0.6)]">
           Browse Offers 
        </button>
      </section>

      <div id="loot-section" className="flex flex-col w-full items-center pt-2 md:pt-4 flex-grow mb-6 md:mb-10 overflow-x-hidden">
        <div className="flex flex-col sm:flex-row gap-3 sm:gap-8 mb-8 md:mb-12 border-b border-retro-gray/50 pb-2 px-2 md:px-6 items-center justify-center">
          <button onClick={() => setView('active')} className={`text-[10px] md:text-sm tracking-widest uppercase transition-colors duration-300 ${view === 'active' ? 'text-retro-orange drop-shadow-[0_0_8px_rgba(255,107,0,0.5)]' : 'text-retro-gray hover:text-retro-white'}`}>
            Active Offers
          </button>
          <button onClick={() => setView('expired')} className={`text-[10px] md:text-sm tracking-widest uppercase transition-colors duration-300 ${view === 'expired' ? 'text-retro-white drop-shadow-[0_0_8px_rgba(255,255,255,0.5)]' : 'text-retro-gray hover:text-retro-white'}`}>
            Expired
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
    </div>
  );
}