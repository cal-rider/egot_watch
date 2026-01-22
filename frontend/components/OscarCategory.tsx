"use client";

import { OscarCategory as OscarCategoryType } from "@/lib/api";
import OscarNominee from "./OscarNominee";

interface OscarCategoryProps {
  category: OscarCategoryType;
}

export default function OscarCategory({ category }: OscarCategoryProps) {
  const hasWinner = category.nominees.some((n) => n.is_winner);

  return (
    <div className="relative">
      {/* Art deco divider */}
      <div className="flex items-center gap-4 mb-8">
        <div className="flex-1 h-px bg-gradient-to-r from-transparent via-gold-500/30 to-gold-500/50" />
        <div className="flex items-center gap-3">
          <span className="text-gold-500/50">◆</span>
          <h2 className="font-display text-2xl md:text-3xl font-bold text-center tracking-wide">
            <span className="text-gold-500">{category.name.toUpperCase()}</span>
          </h2>
          <span className="text-gold-500/50">◆</span>
        </div>
        <div className="flex-1 h-px bg-gradient-to-l from-transparent via-gold-500/30 to-gold-500/50" />
      </div>

      {/* Winner announced badge */}
      {hasWinner && (
        <div className="absolute -top-2 left-1/2 -translate-x-1/2 z-10">
          <div className="bg-gradient-to-r from-[#8b0000] via-[#b22222] to-[#8b0000] text-white text-xs px-3 py-1 rounded-full font-display tracking-wider">
            WINNER ANNOUNCED
          </div>
        </div>
      )}

      {/* Nominees grid */}
      <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-4 md:gap-6">
        {category.nominees.map((nominee) => (
          <OscarNominee
            key={nominee.id}
            nominee={nominee}
            hasWinnerInCategory={hasWinner}
          />
        ))}
      </div>
    </div>
  );
}
