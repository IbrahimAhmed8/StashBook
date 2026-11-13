import { Helmet } from 'react-helmet-async';
import { Link } from 'react-router-dom';

export default function Privacy() {
  return (
    <div className="w-full max-w-4xl mx-auto px-6 py-16 flex-grow flex flex-col items-center">
      <Helmet>
        <title>Privacy Policy | StashBook</title>
        <meta name="description" content="StashBook does not require an account, runs no analytics or advertising trackers, and sets no cookies. Read the full privacy policy." />
      </Helmet>

      <h1 className="text-2xl md:text-4xl lg:text-5xl text-retro-white tracking-[0.1em] md:tracking-[0.2em] uppercase drop-shadow-[3px_3px_0_var(--color-retro-orange)] mb-10 text-center whitespace-nowrap">
        Privacy
      </h1>

      <div className="w-full border-2 border-retro-gray p-6 md:p-10 bg-black shadow-[8px_8px_0_0_rgba(156,163,175,0.1)] flex flex-col gap-8 relative overflow-hidden">

        <div className="absolute inset-0 bg-[linear-gradient(rgba(255,255,255,0.02)_1px,transparent_1px),linear-gradient(90deg,rgba(255,255,255,0.02)_1px,transparent_1px)] bg-[size:20px_20px] pointer-events-none"></div>

        <div className="relative z-10">
          <h3 className="text-retro-orange uppercase tracking-widest border-b border-retro-gray/30 pb-2 mb-4 text-lg md:text-xl">
            [ The Short Version ]
          </h3>
          <p className="text-sm md:text-base text-retro-gray leading-relaxed tracking-wider">
            StashBook is a read-only catalogue of free game offers and books. You can browse the
            whole site without an account, and we do not ask you for personal information to use it.
          </p>
        </div>

        <div className="relative z-10">
          <h3 className="text-retro-white uppercase tracking-widest border-b border-retro-gray/30 pb-2 mb-4 text-lg md:text-xl">
            [ What We Collect ]
          </h3>
          <p className="text-sm md:text-base text-retro-gray leading-relaxed tracking-wider">
            Browsing StashBook does not send us any personal data. The site runs no analytics and
            no advertising trackers, and it does not set cookies or store anything in your browser.
            Like most websites, our hosting provider may process standard technical data (such as
            IP addresses in server logs) to operate and secure the site.
          </p>
        </div>

        <div className="relative z-10">
          <h3 className="text-retro-white uppercase tracking-widest border-b border-retro-gray/30 pb-2 mb-4 text-lg md:text-xl">
            [ Contact Form ]
          </h3>
          <p className="text-sm md:text-base text-retro-gray leading-relaxed tracking-wider">
            If you use the contact form, the name, email address and message you enter are sent to
            us through Web3Forms so we can read and reply to your message. We use that information
            only to respond to you, and we do not sell or share it.
          </p>
        </div>

        <div className="relative z-10">
          <h3 className="text-retro-white uppercase tracking-widest border-b border-retro-gray/30 pb-2 mb-4 text-lg md:text-xl">
            [ Outbound Links ]
          </h3>
          <p className="text-sm md:text-base text-retro-gray leading-relaxed tracking-wider">
            StashBook links out to third-party storefronts (Steam, Epic Games, GOG, itch.io and
            others) and loads game artwork from their content networks. Once you follow one of
            those links, that service&apos;s own privacy policy applies.
          </p>
        </div>

        <div className="relative z-10">
          <h3 className="text-retro-white uppercase tracking-widest border-b border-retro-gray/30 pb-2 mb-4 text-lg md:text-xl">
            [ Children ]
          </h3>
          <p className="text-sm md:text-base text-retro-gray leading-relaxed tracking-wider">
            StashBook is not directed at children under 13, and we do not knowingly collect
            personal information from children.
          </p>
        </div>

        <div className="relative z-10">
          <h3 className="text-retro-white uppercase tracking-widest border-b border-retro-gray/30 pb-2 mb-4 text-lg md:text-xl">
            [ Changes &amp; Questions ]
          </h3>
          <p className="text-sm md:text-base text-retro-gray leading-relaxed tracking-wider">
            If this policy changes, the updated version will be posted on this page. Questions
            about privacy can be sent through the <Link to="/contact" className="text-retro-orange hover:text-retro-white transition-colors duration-300">contact form</Link> or
            to <a href="mailto:contact@stashbook.online" className="text-retro-orange hover:text-retro-white transition-colors duration-300">contact@stashbook.online</a>.
          </p>
          <p className="text-[10px] text-retro-gray/60 uppercase tracking-wider mt-6">
            Last updated: November 2026
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
