import React, { useState } from 'react';

export default function GameCard({ title, platform, image, isExpired, price, onClick }) {
  const [imageLoaded, setImageLoaded] = useState(false);

  const words = title.split(' ');
  const shortTitle = words.slice(0, 3).join(' ') + (words.length > 3 ? '...' : '');

  let badgeText = "FREE";
  let badgeClass = "text-retro-orange animate-pulse";

  if (isExpired) {
    badgeText = "EXPIRED";
    badgeClass = "text-retro-gray/50";
  } else if (price > 0) {
    badgeText = `$${price}`;
    badgeClass = "text-retro-white drop-shadow-[0_0_5px_rgba(255,255,255,0.8)] font-bold";
  }

  return (
    <div
      onClick={onClick}
      className={`w-36 md:w-44 border-2 p-2.5 bg-black flex flex-col gap-3 group transition-all duration-200 relative
                    ${isExpired
                      ? 'border-retro-gray/30 opacity-50 cursor-not-allowed'
                      : 'border-retro-gray hover:border-retro-orange hover:-translate-y-1 hover:shadow-[4px_4px_0_0_var(--color-retro-orange)] cursor-pointer'}`}
    >
      <div
        className={`relative overflow-hidden border transition-colors duration-200 aspect-[460/215] w-full h-24 md:h-28 bg-retro-gray/10
                      ${isExpired ? 'border-retro-gray/30' : 'border-retro-gray group-hover:border-retro-orange'}`}
      >
        {!imageLoaded && (
          <div className="absolute inset-0 bg-retro-gray/20 animate-pulse" aria-hidden="true" />
        )}
        <img
          src={image}
          alt={title}
          width={460}
          height={215}
          loading="lazy"
          decoding="async"
          onLoad={() => setImageLoaded(true)}
          className={`w-full h-full object-cover transition-all duration-300
                      ${!imageLoaded ? 'opacity-0' : isExpired ? 'grayscale opacity-30' : 'grayscale opacity-70 group-hover:grayscale-0 group-hover:opacity-100'}`}
        />
      </div>

      <div className="flex flex-col gap-2 mt-1 relative">
        <h2 className={`text-[10px] md:text-xs uppercase leading-tight truncate
                        ${isExpired ? 'text-retro-gray' : 'text-retro-white'}`}>
          {title}
        </h2>

        <div className="absolute bottom-full left-1/2 -translate-x-1/2 mb-2 w-max max-w-[140px] md:max-w-[160px] z-50
                        opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all duration-200 pointer-events-none">
          <div className={`bg-black border p-1.5 md:p-1.5 text-center
                          ${isExpired
                            ? 'border-retro-gray/50 shadow-[2px_2px_0_0_rgba(156,163,175,0.5)]'
                            : 'border-retro-orange shadow-[2px_2px_0_0_var(--color-retro-orange)]'}`}>
            <p className={`text-[8px] uppercase whitespace-normal break-words leading-relaxed tracking-wider
                            ${isExpired ? 'text-retro-gray' : 'text-retro-white'}`}>
              {shortTitle}
            </p>
          </div>
        </div>

        <div className="flex justify-between items-center mt-1">
          <span className={`text-[8px] uppercase ${isExpired ? 'text-retro-gray/50' : 'text-retro-gray'}`}>
            {platform}
          </span>
          <span className={`text-[9px] uppercase ${badgeClass}`}>
            {badgeText}
          </span>
        </div>
      </div>
    </div>
  );
}
