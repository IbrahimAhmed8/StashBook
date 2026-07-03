import React from 'react';
import { Helmet } from 'react-helmet-async';

export default function Contact() {
  const handleSubmit = (e) => {
    e.preventDefault();
  };

  const inputClass =
    'w-full bg-black border-2 border-retro-gray text-retro-white text-sm md:text-base px-4 py-3 tracking-wider placeholder:text-retro-gray/50 focus:outline-none focus:border-retro-orange focus:shadow-[0_0_10px_rgba(255,140,0,0.3)] transition-all duration-300';

  return (
    <div className="w-full max-w-4xl mx-auto px-6 py-16 flex-grow flex flex-col items-center">
      <Helmet>
        <title>Contact Us - StashBook</title>
        <meta name="description" content="Get in touch with StashBook for inquiries, feedback, or support." />
      </Helmet>

      <h2 className="text-2xl md:text-4xl lg:text-5xl text-retro-white tracking-[0.1em] md:tracking-[0.2em] uppercase drop-shadow-[3px_3px_0_var(--color-retro-orange)] mb-10 text-center whitespace-nowrap">
        Contact Us
      </h2>

      <div className="w-full border-2 border-retro-gray p-6 md:p-10 bg-black shadow-[8px_8px_0_0_rgba(156,163,175,0.1)] flex flex-col gap-8 relative overflow-hidden">
        <div className="absolute inset-0 bg-[linear-gradient(rgba(255,255,255,0.02)_1px,transparent_1px),linear-gradient(90deg,rgba(255,255,255,0.02)_1px,transparent_1px)] bg-[size:20px_20px] pointer-events-none" />

        <form onSubmit={handleSubmit} className="relative z-10 flex flex-col gap-6">
          <div>
            <label htmlFor="name" className="block text-retro-orange uppercase tracking-widest text-xs md:text-sm mb-2">
              Name
            </label>
            <input
              type="text"
              id="name"
              name="name"
              className={inputClass}
              placeholder="Enter your callsign..."
              autoComplete="name"
            />
          </div>

          <div>
            <label htmlFor="email" className="block text-retro-orange uppercase tracking-widest text-xs md:text-sm mb-2">
              Email
            </label>
            <input
              type="email"
              id="email"
              name="email"
              className={inputClass}
              placeholder="your@email.com"
              autoComplete="email"
            />
          </div>

          <div>
            <label htmlFor="message" className="block text-retro-orange uppercase tracking-widest text-xs md:text-sm mb-2">
              Message
            </label>
            <textarea
              id="message"
              name="message"
              rows={6}
              className={`${inputClass} resize-y min-h-[120px]`}
              placeholder="Compose your transmission..."
            />
          </div>

          <button
            type="submit"
            className="mt-2 border-2 border-retro-orange text-retro-orange px-8 py-3 text-xs md:text-sm tracking-widest uppercase hover:bg-retro-orange hover:text-black transition-all duration-300 shadow-[0_0_15px_rgba(255,107,0,0.2)] hover:shadow-[0_0_25px_rgba(255,107,0,0.6)] self-center"
          >
            [ Submit ]
          </button>
        </form>

        <div className="relative z-10 border-t border-retro-gray/30 pt-8 text-center">
          <h3 className="text-retro-white uppercase tracking-widest text-sm md:text-base mb-4">
            Or email us directly:
          </h3>
          <a
            href="mailto:contact@stashbook.online"
            className="text-retro-orange uppercase tracking-widest text-sm md:text-base hover:text-retro-white transition-colors duration-300 drop-shadow-[0_0_8px_rgba(255,107,0,0.5)]"
          >
            contact@stashbook.online
          </a>
        </div>
      </div>
    </div>
  );
}
