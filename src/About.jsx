import React from 'react';
import { Helmet } from 'react-helmet-async';
import { Link } from 'react-router-dom';

export default function About() {
  return (
    <div className="w-full max-w-4xl mx-auto px-6 py-16 flex-grow flex flex-col items-center">
      <Helmet>
        <title>About StashBook – How We Hunt Free Loot</title>
        <meta name="description" content="Learn how StashBook automates free game and digital asset discovery. Independent tool, not affiliated with Steam, Epic, or other storefronts." />
      </Helmet>
      
      <h1 className="text-2xl md:text-4xl lg:text-5xl text-retro-white tracking-[0.1em] md:tracking-[0.2em] uppercase drop-shadow-[3px_3px_0_var(--color-retro-orange)] mb-10 text-center whitespace-nowrap">
        ABOUT
      </h1>

      <div className="w-full border-2 border-retro-gray p-6 md:p-10 bg-black shadow-[8px_8px_0_0_rgba(156,163,175,0.1)] flex flex-col gap-8 relative overflow-hidden">
        
        <div className="absolute inset-0 bg-[linear-gradient(rgba(255,255,255,0.02)_1px,transparent_1px),linear-gradient(90deg,rgba(255,255,255,0.02)_1px,transparent_1px)] bg-[size:20px_20px] pointer-events-none"></div>

        <div className="relative z-10">
          <h3 className="text-retro-orange uppercase tracking-widest border-b border-retro-gray/30 pb-2 mb-4 text-lg md:text-xl">
            [ The Protocol ]
          </h3>
          <p className="text-sm md:text-base text-retro-gray leading-relaxed tracking-wider">
            StashBook is an automated aggregation engine designed with one core directive: to hunt down 100% free digital assets across the web. We deploy background scrapers and API monitors to track limited time giveaways, 'free to play' transitions, and digital loot so you don't have to manually check dozens of storefronts.
          </p>
        </div>

        <div className="relative z-10">
          <h3 className="text-retro-white uppercase tracking-widest border-b border-retro-gray/30 pb-2 mb-4 text-lg md:text-xl">
            [ Disclaimer ]
          </h3>
          <p className="text-sm md:text-base text-retro-gray leading-relaxed tracking-wider border-l-2 border-retro-orange pl-4 bg-retro-orange/5 p-3">
            This project is strictly an independent, community driven tool. <br/><br/>
            <span className="text-retro-white">StashBook is NOT affiliated with, endorsed by, or connected to Epic Games, Valve (Steam), Amazon (Prime Gaming), GOG, itch.io, or any other corporation mentioned on this platform.</span> All game titles, images, and logos belong to their respective copyright holders.
          </p>
        </div>

      </div>

      <Link 
        to="/" 
        className="mt-12 border-2 border-retro-orange text-retro-orange px-8 py-3 text-xs md:text-sm tracking-widest uppercase hover:bg-retro-orange hover:text-black transition-all duration-300 shadow-[0_0_15px_rgba(255,107,0,0.2)]"
      >
        [ Return to Root ]
      </Link>
      
    </div>
  );
}