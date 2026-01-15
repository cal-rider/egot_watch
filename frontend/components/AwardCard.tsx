import { Award } from "@/lib/api";

interface Props {
  award: Award;
}

const awardColors = {
  Emmy: "border-l-pink-500",
  Grammy: "border-l-yellow-500",
  Oscar: "border-l-amber-400",
  Tony: "border-l-red-500",
};

const awardIcons = {
  Emmy: "📺",
  Grammy: "🎵",
  Oscar: "🎬",
  Tony: "🎭",
};

export default function AwardCard({ award }: Props) {
  return (
    <div
      className={`
        award-card rounded-lg p-4
        border-l-4 ${awardColors[award.type]}
      `}
    >
      <div className="flex items-start justify-between gap-4">
        <div className="flex-1">
          <div className="flex items-center gap-2 mb-1">
            <span className="text-xl">{awardIcons[award.type]}</span>
            <span className="text-gold-500 font-display font-semibold">
              {award.type}
            </span>
            {award.year > 0 && (
              <span className="text-gray-400 text-sm">({award.year})</span>
            )}
          </div>
          <p className="text-hollywood-cream text-sm mb-1">{award.category}</p>
          {award.work && (
            <p className="text-gray-400 text-sm italic">
              for &ldquo;{award.work}&rdquo;
            </p>
          )}
        </div>
        {award.is_winner && (
          <span className="text-gold-500 text-2xl" title="Winner">
            ★
          </span>
        )}
      </div>
    </div>
  );
}
