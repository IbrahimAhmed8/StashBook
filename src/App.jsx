import React, { Suspense, lazy } from 'react';
import { BrowserRouter as Router, Routes, Route, Link, useLocation } from 'react-router-dom';
import Home from './Home';
import ErrorBoundary from './ErrorBoundary';

const About = lazy(() => import('./About'));
const IconicItems = lazy(() => import('./IconicItems'));
const Books = lazy(() => import('./Books'));

function RouteFallback() {
  return (
    <div className="w-full flex justify-center py-16">
      <p className="text-retro-gray uppercase text-xs tracking-widest animate-pulse">Loading...</p>
    </div>
  );
}

function Navigation() {
  const location = useLocation();

  const getLinkClass = (path) => {
    const isActive = location.pathname === path;
    return isActive
      ? "text-[8px] md:text-xs lg:text-sm tracking-widest uppercase text-retro-orange drop-shadow-[0_0_8px_rgba(255,107,0,0.5)]"
      : "text-[8px] md:text-xs lg:text-sm tracking-widest uppercase text-retro-gray hover:text-retro-white transition-colors duration-300";
  };

  return (
    <nav className="w-full border-b border-retro-gray/30 bg-black/90 backdrop-blur-md sticky top-0 z-40 flex flex-col md:flex-row justify-center md:justify-between items-center gap-2 md:gap-0 py-2 md:py-4 px-2 md:px-6 lg:px-12">
      <Link to="/" className="text-sm md:text-xl lg:text-2xl font-bold tracking-widest uppercase text-retro-white drop-shadow-[2px_2px_0_var(--color-retro-orange)] cursor-pointer">
        SB
      </Link>
      <div className="flex flex-wrap justify-center gap-3 md:gap-6 lg:gap-10 items-center">
        <Link to="/" className={getLinkClass('/')}>Games</Link>
        <Link to="/books" className={getLinkClass('/books')}>Books</Link>
        <Link to="/iconic" className={getLinkClass('/iconic')}>Iconic Items</Link>
      </div>
    </nav>
  );
}

export default function App() {
  return (
    <Router>
      <div className="min-h-screen flex flex-col items-center relative bg-black overflow-x-hidden w-full max-w-[100vw]">
        <Navigation />

        <ErrorBoundary>
          <Suspense fallback={<RouteFallback />}>
            <Routes>
              <Route path="/" element={<Home />} />
              <Route path="/about" element={<About />} />
              <Route path="/iconic" element={<IconicItems />} />
              <Route path="/books" element={<Books />} />
            </Routes>
          </Suspense>
        </ErrorBoundary>

        <footer className="w-full border-t border-retro-gray/30 bg-black pt-4 md:pt-8 pb-4 md:pb-6 flex flex-col items-center gap-3 md:gap-6 mt-auto px-2">
          <div className="flex flex-wrap justify-center gap-3 md:gap-6 lg:gap-10 text-[8px] md:text-[10px] lg:text-xs tracking-widest uppercase text-retro-gray">
            <Link to="/about" className="hover:text-retro-orange transition-colors duration-300">About</Link>
            <Link to="/" className="hover:text-retro-orange transition-colors duration-300">Contact Us</Link>
          </div>
          <p className="text-[8px] md:text-[9px] tracking-wider text-retro-gray/40 uppercase">
            © {new Date().getFullYear()} StashBook
          </p>
        </footer>
      </div>
    </Router>
  );
}
