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
      ? "text-xs md:text-sm tracking-widest uppercase text-retro-orange drop-shadow-[0_0_8px_rgba(255,107,0,0.5)]"
      : "text-xs md:text-sm tracking-widest uppercase text-retro-gray hover:text-retro-white transition-colors duration-300";
  };

  return (
    <nav className="w-full border-b border-retro-gray/30 bg-black/90 backdrop-blur-md sticky top-0 z-40 flex justify-between items-center py-4 px-6 md:px-12">
      <Link to="/" className="text-xl md:text-2xl font-bold tracking-widest uppercase text-retro-white drop-shadow-[2px_2px_0_var(--color-retro-orange)] cursor-pointer">
        SB
      </Link>
      <div className="flex gap-6 md:gap-10 items-center">
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
      <div className="min-h-screen flex flex-col items-center relative bg-black">
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

        <footer className="w-full border-t border-retro-gray/30 bg-black pt-8 pb-6 flex flex-col items-center gap-6 mt-auto">
          <div className="flex flex-wrap justify-center gap-6 md:gap-10 text-[10px] md:text-xs tracking-widest uppercase text-retro-gray">
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
