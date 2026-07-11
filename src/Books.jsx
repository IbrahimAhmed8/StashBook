import React, { useState } from 'react';
import { Helmet } from 'react-helmet-async';
import GameCard from './GameCard';
import ServiceUnavailable from './ServiceUnavailable';
import { useStashData } from './hooks/useStashData';

const DEFAULT_DESCRIPTION = 'StashBook is your daily loot drop for free game deals, developer books, and curated tech gear—save money and level up your setup without the grind.';
const DEFAULT_KEYWORDS = 'free books, developer books, tech books, stashbook, free pdf';

export default function Books() {
  const { data, loading, error: fetchError } = useStashData('/books.json');
  const books = data || [];
  const [selectedBook, setSelectedBook] = useState(null);

  const closeModal = () => setSelectedBook(null);

  return (
    <div className="w-full flex-grow flex flex-col items-center pt-6 md:pt-10 pb-10 md:pb-16 px-2 md:px-6 overflow-x-hidden">
      <Helmet>
        <title>Tech Archives – Free Developer Books | StashBook</title>
        <meta name="description" content="Discover free programming and developer books curated by StashBook. Download PDFs on software engineering, game dev, and tech topics at no cost." />
      </Helmet>
      <h1 className="text-lg sm:text-2xl md:text-3xl lg:text-5xl text-retro-orange tracking-[0.1em] md:tracking-[0.3em] uppercase drop-shadow-[3px_3px_0_rgba(255,107,0,0.5)] mb-8 md:mb-12 text-center break-words px-2">
        [ Tech Archives ]
      </h1>

      {loading ? (
        <div className="text-center text-[10px] md:text-sm text-retro-gray uppercase animate-pulse">Decrypting Library...</div>
      ) : fetchError ? (
        <ServiceUnavailable />
      ) : books.length === 0 ? (
        <div className="text-center text-[10px] md:text-sm text-retro-gray uppercase mt-8 opacity-50">No digital assets found.</div>
      ) : (
        <div className="flex flex-wrap justify-center gap-4 w-full max-w-5xl mx-auto">
          {books.map(book => (
            <GameCard
              key={book.id}
              title={book.title}
              platform={book.author}
              image={book.image}
              isExpired={false}
              price={0}
              titleTag="h2"
              onClick={() => setSelectedBook(book)}
            />
          ))}
        </div>
      )}

      {selectedBook && (
        <>
        <Helmet>
          <title>{selectedBook.title} - Stashbook Free Deals</title>
          <meta name="description" content={selectedBook.description?.trim() || DEFAULT_DESCRIPTION} />
          <meta name="keywords" content={selectedBook.seo_keywords?.trim() || DEFAULT_KEYWORDS} />
        </Helmet>
        <div className="fixed inset-0 z-[100] flex items-center justify-center p-2 md:p-4 bg-black/80 backdrop-blur-sm overflow-x-hidden">
          <div className="absolute inset-0 cursor-pointer" onClick={closeModal}></div>
          <div className="relative bg-black border-2 border-retro-orange p-4 md:p-8 w-full max-w-xl shadow-[8px_8px_0_0_var(--color-retro-orange)] animate-[fadeIn_0.2s_ease-out] overflow-hidden">
            <button onClick={closeModal} className="absolute top-2 right-3 md:top-3 md:right-4 text-retro-gray hover:text-retro-orange text-sm md:text-xl font-bold transition-colors z-10">[X]</button>
            <div className="flex flex-col gap-3 md:gap-5 mt-4">
              <div className="flex flex-col gap-2">
                <h2 className="text-sm md:text-2xl lg:text-3xl text-retro-white uppercase tracking-wider leading-tight break-words line-clamp-3 pr-6">{selectedBook.title}</h2>
                <div className="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-1 sm:gap-2 border-b border-retro-gray/30 pb-3">
                  <span className="text-[10px] md:text-sm text-retro-gray uppercase tracking-wider truncate">{selectedBook.author}</span>
                  <span className="text-[10px] md:text-sm text-retro-orange animate-pulse font-bold tracking-widest shrink-0">FREE PDF</span>
                </div>
              </div>
              <div className="py-2">
                <h3 className="text-[8px] md:text-[10px] text-retro-gray uppercase tracking-[0.2em] mb-2">Synopsis</h3>
                <p className="text-[10px] md:text-xs lg:text-sm text-retro-white leading-relaxed max-h-48 overflow-y-auto pr-2 md:pr-3 custom-scrollbar break-words">
                  {selectedBook.description || "No metadata available."}
                </p>
              </div>
              <div className="flex flex-col gap-3 md:gap-4 mt-2 md:mt-4 pt-3 md:pt-4 border-t border-retro-gray/30">
                <a href={selectedBook.book_url || "#"} target="_blank" rel="noopener noreferrer" className="block text-center py-2 md:py-3 border-2 border-retro-orange text-retro-orange uppercase text-[10px] md:text-xs lg:text-sm tracking-[0.2em] md:tracking-[0.3em] font-bold hover:bg-retro-orange hover:text-black hover:shadow-[0_0_15px_var(--color-retro-orange)] transition-all duration-300">
                  Download Asset
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