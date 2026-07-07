import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Helmet } from 'react-helmet-async';
import GameCard from './GameCard';
import ServiceUnavailable from './ServiceUnavailable';
import { filterIconicDeals } from './utils/gameFilters';
import { slugify } from './utils/stringUtils';

let cachedIconicGames = null;

export default function IconicItems() {
  const navigate = useNavigate();
  const [games, setGames] = useState([]);
  const [loading, setLoading] = useState(true);
  const [fetchError, setFetchError] = useState(false);

  useEffect(() => {
    if (cachedIconicGames) {
      setGames(cachedIconicGames);
      setLoading(false);
    }

    fetch('/games.json')
      .then((res) => {
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        return res.json();
      })
      .then((data) => {
        const filtered = filterIconicDeals(data);
        cachedIconicGames = filtered;
        setGames(filtered);
        setLoading(false);
      })
      .catch((err) => {
        console.error("Error fetching data:", err);
        if (!cachedIconicGames) {
          setFetchError(true);
          setLoading(false);
        }
      });
  }, []);

  return (
    <div className="w-full flex-grow flex flex-col items-center pt-6 md:pt-10 pb-10 md:pb-16 px-2 md:px-6 overflow-x-hidden">
      <Helmet>
        <title>Iconic Drops – Legendary Free Games | StashBook</title>
        <meta name="description" content="Track iconic and legendary free game drops on StashBook. Limited-time offers on acclaimed titles from Steam, Epic, and other platforms." />
      </Helmet>
      <h1 className="text-lg sm:text-2xl md:text-3xl lg:text-5xl text-retro-orange tracking-[0.1em] md:tracking-[0.3em] uppercase drop-shadow-[3px_3px_0_rgba(255,107,0,0.5)] mb-8 md:mb-12 text-center break-words px-2">
        [ Iconic Drops ]
      </h1>

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
              description={game.description}
              titleTag="h2"
              onClick={() => navigate(`/game/${slugify(game.title)}`)}
            />
          ))}
        </div>
      )}
    </div>
  );
}