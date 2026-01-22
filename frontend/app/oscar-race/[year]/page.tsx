"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useParams } from "next/navigation";
import { OscarCeremony, getOscarCeremony } from "@/lib/api";
import OscarCategory from "@/components/OscarCategory";

export default function OscarRacePage() {
  const params = useParams();
  const year = parseInt(params.year as string, 10);

  const [ceremony, setCeremony] = useState<OscarCeremony | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function fetchData() {
      try {
        const data = await getOscarCeremony(year);
        setCeremony(data);
      } catch (err) {
        setError("Failed to load Oscar ceremony data");
        console.error(err);
      } finally {
        setLoading(false);
      }
    }
    if (year) {
      fetchData();
    }
  }, [year]);

  return (
    <main className="min-h-screen relative overflow-hidden">
      {/* Red carpet gradient background */}
      <div className="fixed inset-0 bg-gradient-to-b from-[#1a0505] via-[#2d0a0a] to-[#0a0505] -z-20" />

      {/* Sweeping spotlight effects */}
      <div className="fixed inset-0 -z-10 overflow-hidden pointer-events-none">
        <div className="spotlight-sweep spotlight-left" />
        <div className="spotlight-sweep spotlight-right" />
      </div>

      {/* Floating champagne bubbles */}
      <div className="fixed inset-0 -z-5 pointer-events-none overflow-hidden">
        {[...Array(15)].map((_, i) => (
          <div
            key={i}
            className="champagne-bubble"
            style={{
              left: `${5 + Math.random() * 90}%`,
              animationDelay: `${Math.random() * 5}s`,
              animationDuration: `${6 + Math.random() * 4}s`,
            }}
          />
        ))}
      </div>

      <div className="relative px-4 py-12">
        <div className="max-w-6xl mx-auto">
          {/* Back link */}
          <Link
            href="/"
            className="text-gold-500 hover:text-gold-400 font-display transition-colors inline-flex items-center gap-2"
          >
            <span className="text-xl">←</span> Back to EGOT Tracker
          </Link>

          {/* Header */}
          <div className="text-center my-12">
            {/* Oscar statue decorations */}
            <div className="flex justify-center items-center gap-6 mb-6">
              <span className="text-4xl opacity-60">🏆</span>
              <div className="text-gold-500 text-2xl tracking-[0.3em]">THE ACADEMY AWARDS</div>
              <span className="text-4xl opacity-60">🏆</span>
            </div>

            {loading ? (
              <h1 className="font-display text-5xl md:text-7xl font-bold mb-4">
                <span className="oscar-shimmer">Loading...</span>
              </h1>
            ) : ceremony ? (
              <>
                <h1 className="font-display text-5xl md:text-7xl font-bold mb-4">
                  <span className="oscar-shimmer">{ceremony.ceremony_name || `${year} Oscars`}</span>
                </h1>
                {ceremony.ceremony_date && (
                  <p className="text-hollywood-cream/70 text-lg">
                    {new Date(ceremony.ceremony_date).toLocaleDateString("en-US", {
                      weekday: "long",
                      year: "numeric",
                      month: "long",
                      day: "numeric",
                    })}
                  </p>
                )}
              </>
            ) : (
              <h1 className="font-display text-5xl md:text-7xl font-bold mb-4">
                <span className="oscar-shimmer">{year} Oscar Race</span>
              </h1>
            )}

            <div className="oscar-divider w-96 mx-auto my-8" />
          </div>

          {/* Loading state */}
          {loading && (
            <div className="text-center py-16">
              <div className="flex justify-center gap-4 mb-6">
                {[0, 1, 2, 3, 4].map((i) => (
                  <div
                    key={i}
                    className="text-3xl animate-pulse"
                    style={{ animationDelay: `${i * 0.15}s` }}
                  >
                    🏆
                  </div>
                ))}
              </div>
              <p className="text-hollywood-cream/60 text-lg font-display">
                Rolling out the red carpet...
              </p>
            </div>
          )}

          {/* Error state */}
          {error && (
            <div className="text-center py-16">
              <div className="text-red-400 text-5xl mb-6">✗</div>
              <p className="text-hollywood-cream text-xl font-display mb-2">{error}</p>
              <p className="text-hollywood-cream/50">
                No Oscar ceremony found for {year}. The ceremony may not have been set up yet.
              </p>
            </div>
          )}

          {/* Categories */}
          {!loading && !error && ceremony && (
            <div className="space-y-12">
              {ceremony.categories.length === 0 ? (
                <div className="text-center py-16">
                  <p className="text-hollywood-cream/60">
                    No categories found for this ceremony.
                  </p>
                </div>
              ) : (
                ceremony.categories.map((category) => (
                  <OscarCategory key={category.id} category={category} />
                ))
              )}
            </div>
          )}
        </div>
      </div>

      {/* Red carpet border effect */}
      <div className="fixed left-0 top-0 bottom-0 w-4 bg-gradient-to-r from-[#8b0000] to-transparent opacity-30" />
      <div className="fixed right-0 top-0 bottom-0 w-4 bg-gradient-to-l from-[#8b0000] to-transparent opacity-30" />

      {/* Corner Oscar decorations */}
      <div className="fixed top-6 left-6 text-2xl text-gold-500/20">🏆</div>
      <div className="fixed top-6 right-6 text-2xl text-gold-500/20">🏆</div>
      <div className="fixed bottom-6 left-6 text-2xl text-gold-500/20">🏆</div>
      <div className="fixed bottom-6 right-6 text-2xl text-gold-500/20">🏆</div>

      <style jsx>{`
        .oscar-shimmer {
          background: linear-gradient(
            90deg,
            #b8960f 0%,
            #f9d054 25%,
            #d4af37 50%,
            #f9d054 75%,
            #b8960f 100%
          );
          background-size: 200% auto;
          -webkit-background-clip: text;
          background-clip: text;
          -webkit-text-fill-color: transparent;
          animation: shimmer 3s linear infinite;
        }

        @keyframes shimmer {
          0% { background-position: -200% center; }
          100% { background-position: 200% center; }
        }

        .oscar-divider {
          height: 2px;
          background: linear-gradient(
            90deg,
            transparent 0%,
            #8b0000 20%,
            #d4af37 50%,
            #8b0000 80%,
            transparent 100%
          );
        }

        .spotlight-sweep {
          position: absolute;
          width: 200px;
          height: 100%;
          background: linear-gradient(
            90deg,
            transparent,
            rgba(212, 175, 55, 0.03),
            transparent
          );
          animation: sweep 8s ease-in-out infinite;
        }

        .spotlight-left {
          left: -200px;
          animation-delay: 0s;
        }

        .spotlight-right {
          right: -200px;
          animation-direction: reverse;
          animation-delay: 4s;
        }

        @keyframes sweep {
          0%, 100% { transform: translateX(0); }
          50% { transform: translateX(calc(100vw + 200px)); }
        }

        .champagne-bubble {
          position: absolute;
          bottom: -20px;
          width: 6px;
          height: 6px;
          background: radial-gradient(circle, rgba(212, 175, 55, 0.6) 0%, transparent 70%);
          border-radius: 50%;
          animation: rise linear infinite;
        }

        @keyframes rise {
          0% {
            transform: translateY(0) scale(1);
            opacity: 0.6;
          }
          100% {
            transform: translateY(-100vh) scale(0.5);
            opacity: 0;
          }
        }
      `}</style>
    </main>
  );
}
