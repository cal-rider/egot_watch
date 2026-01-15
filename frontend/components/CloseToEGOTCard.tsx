"use client";

import Link from "next/link";
import { CelebrityWithProgress } from "@/lib/api";

interface Props {
  celebrity: CelebrityWithProgress;
}

const awardTypes = ["Emmy", "Grammy", "Oscar", "Tony"] as const;
const awardLetters: Record<string, string> = {
  Emmy: "E",
  Grammy: "G",
  Oscar: "O",
  Tony: "T",
};

export default function CloseToEGOTCard({ celebrity }: Props) {
  const missingAward = awardTypes.find((a) => !celebrity.won_awards.includes(a));

  return (
    <Link href={`/celebrity/${encodeURIComponent(celebrity.name)}`}>
      <div className="award-card rounded-lg p-4 hover:border-gold-500/50 transition-all cursor-pointer">
        <div className="flex items-center gap-4">
          {/* Photo or placeholder */}
          {celebrity.photo_url ? (
            <img
              src={celebrity.photo_url}
              alt=""
              className="w-14 h-14 rounded-full object-cover border-2 border-gold-500/30"
            />
          ) : (
            <div className="w-14 h-14 rounded-full bg-gold-500/20 flex items-center justify-center border-2 border-gold-500/30">
              <span className="text-gold-500 text-xl">★</span>
            </div>
          )}

          <div className="flex-1 min-w-0">
            <h3 className="font-display text-lg text-hollywood-cream truncate">
              {celebrity.name}
            </h3>

            {/* EGOT progress indicator */}
            <div className="flex gap-1.5 mt-2">
              {awardTypes.map((type) => {
                const won = celebrity.won_awards.includes(type);
                return (
                  <span
                    key={type}
                    className={`w-7 h-7 rounded-full flex items-center justify-center text-xs font-bold transition-all
                      ${
                        won
                          ? "bg-gold-500 text-hollywood-black shadow-sm"
                          : "bg-hollywood-dark text-gray-600 border border-gray-700"
                      }`}
                    title={type}
                  >
                    {awardLetters[type]}
                  </span>
                );
              })}
            </div>

            <p className="text-gray-400 text-sm mt-2">
              Needs: <span className="text-gold-400 font-semibold">{missingAward}</span>
            </p>
          </div>
        </div>
      </div>
    </Link>
  );
}
