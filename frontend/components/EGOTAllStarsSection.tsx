"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { CelebrityWithProgress, getEGOTWinners } from "@/lib/api";

const awardLetters: Record<string, string> = {
  Emmy: "E",
  Grammy: "G",
  Oscar: "O",
  Tony: "T",
};

function EGOTCard({ celebrity }: { celebrity: CelebrityWithProgress }) {
  return (
    <Link href={`/celebrity/${encodeURIComponent(celebrity.name)}`}>
      <div className="award-card rounded-lg p-4 hover:border-gold-500/50 transition-all cursor-pointer bg-gradient-to-br from-gold-500/10 to-transparent">
        <div className="flex items-center gap-4">
          {celebrity.photo_url ? (
            <img
              src={celebrity.photo_url}
              alt=""
              className="w-14 h-14 rounded-full object-cover border-2 border-gold-500"
            />
          ) : (
            <div className="w-14 h-14 rounded-full bg-gold-500/20 flex items-center justify-center border-2 border-gold-500">
              <span className="text-gold-500 text-xl">★</span>
            </div>
          )}

          <div className="flex-1 min-w-0">
            <h3 className="font-display text-lg text-hollywood-cream truncate">
              {celebrity.name}
            </h3>

            <div className="flex gap-1.5 mt-2">
              {["Emmy", "Grammy", "Oscar", "Tony"].map((type) => (
                <span
                  key={type}
                  className="w-7 h-7 rounded-full flex items-center justify-center text-xs font-bold bg-gold-500 text-hollywood-black shadow-sm"
                  title={type}
                >
                  {awardLetters[type]}
                </span>
              ))}
            </div>

            <p className="text-gold-400 text-sm mt-2 font-semibold">
              EGOT Winner
            </p>
          </div>
        </div>
      </div>
    </Link>
  );
}

export default function EGOTAllStarsSection() {
  const [celebrities, setCelebrities] = useState<CelebrityWithProgress[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function fetchData() {
      try {
        const data = await getEGOTWinners(6);
        setCelebrities(data);
      } catch (err) {
        setError("Failed to load data");
        console.error("Failed to fetch EGOT winners:", err);
      } finally {
        setLoading(false);
      }
    }
    fetchData();
  }, []);

  if (loading) {
    return (
      <section className="w-full max-w-4xl mx-auto mt-16 px-4">
        <div className="text-center py-8">
          <div className="text-gold-500 animate-pulse text-2xl">★ ★ ★ ★</div>
          <p className="text-gray-400 mt-2">Loading...</p>
        </div>
      </section>
    );
  }

  if (error || celebrities.length === 0) {
    return null;
  }

  return (
    <section className="w-full max-w-4xl mx-auto mt-16 px-4">
      <div className="art-deco-line mb-6" />

      <div className="flex items-center justify-between mb-6">
        <h2 className="font-display text-xl md:text-2xl text-gold-500 tracking-wider">
          ★ EGOT ALL STARS ★
        </h2>
        <span className="text-sm text-gray-400 font-display">
          {celebrities.length} Legends
        </span>
      </div>

      <p className="text-gray-400 text-sm mb-6">
        The elite few who have conquered Emmy, Grammy, Oscar, and Tony
      </p>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {celebrities.map((celeb) => (
          <EGOTCard key={celeb.id} celebrity={celeb} />
        ))}
      </div>
    </section>
  );
}
