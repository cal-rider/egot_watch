"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { CelebrityBasic, getNoAwards } from "@/lib/api";

function LoserCard({ celebrity }: { celebrity: CelebrityBasic }) {
  return (
    <Link href={`/celebrity/${encodeURIComponent(celebrity.name)}`}>
      <div className="award-card rounded-lg p-4 hover:border-gray-500/50 transition-all cursor-pointer bg-gradient-to-br from-gray-800/30 to-transparent">
        <div className="flex items-center gap-4">
          {celebrity.photo_url ? (
            <img
              src={celebrity.photo_url}
              alt=""
              className="w-14 h-14 rounded-full object-cover border-2 border-gray-600 grayscale opacity-80"
            />
          ) : (
            <div className="w-14 h-14 rounded-full bg-gray-800 flex items-center justify-center border-2 border-gray-600">
              <span className="text-gray-500 text-xl">?</span>
            </div>
          )}

          <div className="flex-1 min-w-0">
            <h3 className="font-display text-lg text-gray-300 truncate">
              {celebrity.name}
            </h3>

            <div className="flex gap-1.5 mt-2">
              {["E", "G", "O", "T"].map((letter) => (
                <span
                  key={letter}
                  className="w-7 h-7 rounded-full flex items-center justify-center text-xs font-bold bg-hollywood-dark text-gray-600 border border-gray-700"
                >
                  {letter}
                </span>
              ))}
            </div>

            <p className="text-gray-500 text-sm mt-2">
              No EGOT awards yet
            </p>
          </div>
        </div>
      </div>
    </Link>
  );
}

export default function BigLosersSection() {
  const [celebrities, setCelebrities] = useState<CelebrityBasic[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function fetchData() {
      try {
        const data = await getNoAwards(6);
        setCelebrities(data);
      } catch (err) {
        setError("Failed to load data");
        console.error("Failed to fetch no awards:", err);
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
          <div className="text-gray-500 animate-pulse text-2xl">. . .</div>
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
      <div className="w-full h-px bg-gradient-to-r from-transparent via-gray-600 to-transparent mb-6" />

      <div className="flex items-center justify-between mb-6">
        <h2 className="font-display text-xl md:text-2xl text-gray-400 tracking-wider">
          ★ BIG LOSERS ★
        </h2>
        <span className="text-sm text-gray-500 font-display">
          {celebrities.length} Hopefuls
        </span>
      </div>

      <p className="text-gray-500 text-sm mb-6">
        Famous faces still waiting for their first EGOT win
      </p>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {celebrities.map((celeb) => (
          <LoserCard key={celeb.id} celebrity={celeb} />
        ))}
      </div>
    </section>
  );
}
