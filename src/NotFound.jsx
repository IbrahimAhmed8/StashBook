import { Helmet } from 'react-helmet-async';
import { Link } from 'react-router-dom';

export default function NotFound() {
  return (
    <main className="w-full flex-grow flex flex-col items-center justify-center min-h-[50vh] px-2 md:px-6">
      <Helmet>
        <title>404 – Page Not Found | StashBook</title>
        <meta name="robots" content="noindex" />
      </Helmet>
      <div className="text-center">
        <h1 className="text-lg sm:text-2xl md:text-4xl text-retro-orange uppercase tracking-widest drop-shadow-[2px_2px_0_var(--color-retro-orange)] mb-4">
          ERROR 404: PAGE NOT FOUND
        </h1>
        <p className="text-[10px] md:text-sm text-retro-gray uppercase tracking-wider mb-8">
          This page does not exist or was moved.
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
