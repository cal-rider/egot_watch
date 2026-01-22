"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import { Celebrity, Award, searchCelebrity } from "@/lib/api";
import CelebrityHeader from "@/components/CelebrityHeader";
import AwardCard from "@/components/AwardCard";

const loadingMessages = [
  "Searching the archives...",
  "Checking our records...",
  "Fetching from Wikidata...",
  "Gathering award history...",
  "Almost there...",
];

export default function CelebrityPage() {
  const params = useParams();
  const router = useRouter();
  const [celebrity, setCelebrity] = useState<Celebrity | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [loadingMessageIndex, setLoadingMessageIndex] = useState(0);

  const slug = params.slug as string;

  // Cycle through loading messages while fetching
  useEffect(() => {
    if (!loading) return;

    const interval = setInterval(() => {
      setLoadingMessageIndex((prev) =>
        prev < loadingMessages.length - 1 ? prev + 1 : prev
      );
    }, 2000);

    return () => clearInterval(interval);
  }, [loading]);

  useEffect(() => {
    async function fetchCelebrity() {
      try {
        setLoading(true);
        setError(null);
        setLoadingMessageIndex(0);
        const data = await searchCelebrity(decodeURIComponent(slug));
        setCelebrity(data);
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to load celebrity");
      } finally {
        setLoading(false);
      }
    }

    if (slug) {
      fetchCelebrity();
    }
  }, [slug]);

  // Group awards by type
  const groupedAwards = celebrity?.awards.reduce(
    (acc, award) => {
      if (!acc[award.type]) {
        acc[award.type] = [];
      }
      acc[award.type].push(award);
      return acc;
    },
    {} as Record<string, Award[]>
  );

  const awardOrder = ["Oscar", "Emmy", "Grammy", "Tony"] as const;

  return (
    <main className="min-h-screen px-4 py-12">
      {/* Back link */}
      <div className="max-w-4xl mx-auto mb-8">
        <Link
          href="/"
          className="text-gold-500 hover:text-gold-400 font-display transition-colors"
        >
          ← Back to Search
        </Link>
      </div>

      <div className="max-w-4xl mx-auto">
        {loading && (
          <div className="text-center py-20">
            {/* Animated stars */}
            <div className="flex justify-center gap-2 mb-6">
              {[0, 1, 2, 3].map((i) => (
                <div
                  key={i}
                  className="text-gold-500 text-2xl animate-pulse"
                  style={{ animationDelay: `${i * 0.2}s` }}
                >
                  ★
                </div>
              ))}
            </div>

            {/* Loading message */}
            <p className="text-hollywood-cream text-xl font-display mb-2">
              {loadingMessages[loadingMessageIndex]}
            </p>

            {/* Progress hint for new celebrities */}
            {loadingMessageIndex >= 2 && (
              <p className="text-gray-500 text-sm mt-4 animate-pulse">
                New celebrity - fetching from Wikidata may take a moment...
              </p>
            )}

            {/* Loading bar */}
            <div className="w-64 mx-auto mt-8 h-1 bg-hollywood-black/50 rounded-full overflow-hidden">
              <div
                className="h-full bg-gold-500 rounded-full transition-all duration-1000 ease-out"
                style={{
                  width: `${Math.min(((loadingMessageIndex + 1) / loadingMessages.length) * 100, 90)}%`,
                }}
              />
            </div>
          </div>
        )}

        {error && (
          <div className="text-center py-20">
            <div className="text-red-500 text-4xl mb-4">✗</div>
            <p className="text-hollywood-cream text-xl font-display mb-2">
              {error}
            </p>
            <p className="text-gray-400 mb-8">
              {error.includes("not found")
                ? "This person may not have EGOT-eligible awards in Wikidata."
                : "Please check your connection and try again."}
            </p>
            <button
              onClick={() => router.push("/")}
              className="px-6 py-3 bg-gold-500 text-hollywood-black font-display font-semibold rounded-lg hover:bg-gold-400 transition-colors"
            >
              Try Another Search
            </button>
          </div>
        )}

        {celebrity && (
          <>
            <CelebrityHeader celebrity={celebrity} />

            {/* Biography Summary */}
            {celebrity.summary && (
              <div className="mt-8 p-6 bg-hollywood-black/30 border border-gold-500/20 rounded-lg">
                <p className="text-hollywood-cream/90 leading-relaxed">
                  {celebrity.summary}
                </p>
              </div>
            )}

            {/* Awards Section */}
            <div className="mt-12">
              <div className="art-deco-line mb-8" />
              <h2 className="font-display text-2xl text-gold-500 text-center mb-8 tracking-wider">
                ★ AWARDS ★
              </h2>

              {awardOrder.map((type) => {
                const awards = groupedAwards?.[type];
                if (!awards || awards.length === 0) return null;

                return (
                  <div key={type} className="mb-8">
                    <h3 className="font-display text-xl text-hollywood-cream mb-4">
                      {type} Awards ({awards.length})
                    </h3>
                    <div className="grid gap-4 md:grid-cols-2">
                      {awards
                        .sort((a, b) => b.year - a.year)
                        .map((award) => (
                          <AwardCard key={award.id} award={award} />
                        ))}
                    </div>
                  </div>
                );
              })}

              {celebrity.awards.length === 0 && (
                <p className="text-center text-gray-400 py-8">
                  No EGOT awards found for this celebrity.
                </p>
              )}
            </div>
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
