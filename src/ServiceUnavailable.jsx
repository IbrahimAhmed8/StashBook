export default function ServiceUnavailable() {
  return (
    <div className="w-full flex flex-col items-center justify-center min-h-[30vh] px-6 py-16">
      <p className="text-sm md:text-base text-retro-orange uppercase tracking-widest text-center animate-pulse">
        Service Temporarily Unavailable
      </p>
      <p className="text-[10px] text-retro-gray uppercase tracking-wider mt-4 text-center opacity-60">
        Unable to reach the aggregation server. Please try again shortly.
      </p>
    </div>
  );
}
