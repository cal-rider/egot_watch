import SearchBar from "@/components/SearchBar";
import CloseToEGOTSection from "@/components/CloseToEGOTSection";

export default function Home() {
  return (
    <main className="min-h-screen flex flex-col items-center px-4 py-16">
      {/* Art deco decorative top */}
      <div className="absolute top-0 left-0 right-0 h-1 bg-gradient-to-r from-transparent via-gold-500 to-transparent opacity-50" />

      {/* Main content */}
      <div className="text-center mb-12 mt-8">
        {/* Decorative stars */}
        <div className="text-gold-500 text-2xl mb-4 tracking-[1em]">★ ★ ★ ★</div>

        {/* Title */}
        <h1 className="font-display text-6xl md:text-8xl font-bold mb-2">
          <span className="gold-shimmer">EGOT</span>
        </h1>
        <h2 className="font-display text-2xl md:text-3xl text-hollywood-cream tracking-[0.3em] mb-2">
          TRACKER
        </h2>

        {/* Decorative line */}
        <div className="art-deco-line w-64 mx-auto my-6" />

        {/* Subtitle */}
        <p className="text-gray-400 font-body text-lg max-w-md mx-auto">
          Discover who has achieved the prestigious Emmy, Grammy, Oscar &amp; Tony grand slam
        </p>
      </div>

      {/* Search */}
      <SearchBar />

      {/* Close to EGOT Featured Section */}
      <CloseToEGOTSection />

      {/* Footer decoration */}
      <div className="mt-16 text-center">
        <div className="text-gold-500/30 text-sm tracking-[0.5em] font-display">
          ★ HOLLYWOOD ★
        </div>
      </div>

      {/* Corner decorations */}
      <div className="fixed top-4 left-4 w-16 h-16 border-l-2 border-t-2 border-gold-500/30" />
      <div className="fixed top-4 right-4 w-16 h-16 border-r-2 border-t-2 border-gold-500/30" />
      <div className="fixed bottom-4 left-4 w-16 h-16 border-l-2 border-b-2 border-gold-500/30" />
      <div className="fixed bottom-4 right-4 w-16 h-16 border-r-2 border-b-2 border-gold-500/30" />
    </main>
  );
}
