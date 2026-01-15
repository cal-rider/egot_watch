"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { CelebrityWithProgress, getCloseToEGOT } from "@/lib/api";
import CloseToEGOTCard from "@/components/CloseToEGOTCard";

export default function CloseToEGOTPage() {
  const [celebrities, setCelebrities] = useState<CelebrityWithProgress[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function fetchData() {
      try {
        const data = await getCloseToEGOT(100);
        setCelebrities(data);
      } catch (err) {
        setError("Failed to load data");
        console.error(err);
      } finally {
        setLoading(false);
      }
    }
    fetchData();
  }, []);

  return (
    <main className="min-h-screen px-4 py-12">
      <div className="max-w-4xl mx-auto">
        {/* Back link */}
        <Link
          href="/"
          className="text-gold-500 hover:text-gold-400 font-display transition-colors"
        >
          ← Back to Search
        </Link>

        {/* Header */}
        <div className="text-center my-12">
          <div className="text-gold-500 text-2xl mb-4 tracking-[0.5em]">★ ★ ★</div>
          <h1 className="font-display text-4xl md:text-5xl font-bold mb-4">
            <span className="gold-shimmer">Close to EGOT</span>
          </h1>
          <div className="art-deco-line w-64 mx-auto my-6" />
          <p className="text-gray-400 max-w-lg mx-auto">
            These legendary artists have achieved 3 of 4 EGOT awards.
            They stand on the precipice of entertainment history.
          </p>
        </div>

        {/* Loading state */}
        {loading && (
          <div className="text-center py-12">
            <div className="flex justify-center gap-2 mb-4">
              {[0, 1, 2, 3].map((i) => (
                <div
                  key={i}
                  className="text-gold-500 text-xl animate-pulse"
                  style={{ animationDelay: `${i * 0.2}s` }}
                >
                  ★
                </div>
              ))}
            </div>
            <p className="text-gray-400">Loading legends...</p>
          </div>
        )}

        {/* Error state */}
        {error && (
          <div className="text-center py-12">
            <div className="text-red-500 text-4xl mb-4">✗</div>
            <p className="text-hollywood-cream text-xl font-display mb-2">{error}</p>
            <p className="text-gray-400">Please try again later.</p>
          </div>
        )}

        {/* Content */}
        {!loading && !error && (
          <>
            {celebrities.length === 0 ? (
              <div className="text-center py-12">
                <p className="text-gray-400">
                  No celebrities found who are close to EGOT status.
                </p>
                <p className="text-gray-500 text-sm mt-2">
                  Search for celebrities to add them to the database.
                </p>
              </div>
            ) : (
              <>
                <p className="text-center text-gray-500 mb-8">
                  {celebrities.length} artist{celebrities.length !== 1 ? "s" : ""} one award away from EGOT
                </p>
                <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                  {celebrities.map((celeb) => (
                    <CloseToEGOTCard key={celeb.id} celebrity={celeb} />
                  ))}
                </div>
              </>
            )}
          </>
        )}
      </div>

      {/* Corner decorations */}
      <div className="fixed top-4 left-4 w-12 h-12 border-l-2 border-t-2 border-gold-500/20" />
      <div className="fixed top-4 right-4 w-12 h-12 border-r-2 border-t-2 border-gold-500/20" />
      <div className="fixed bottom-4 left-4 w-12 h-12 border-l-2 border-b-2 border-gold-500/20" />
      <div className="fixed bottom-4 right-4 w-12 h-12 border-r-2 border-b-2 border-gold-500/20" />
    </main>
  );
}
