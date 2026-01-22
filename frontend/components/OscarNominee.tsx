"use client";

import { OscarNominee as OscarNomineeType } from "@/lib/api";
import { useState, useEffect } from "react";

interface OscarNomineeProps {
  nominee: OscarNomineeType;
  hasWinnerInCategory: boolean;
}

export default function OscarNominee({ nominee, hasWinnerInCategory }: OscarNomineeProps) {
  const [showConfetti, setShowConfetti] = useState(false);
  const isWinner = nominee.is_winner;
  const isDimmed = hasWinnerInCategory && !isWinner;

  useEffect(() => {
    if (isWinner) {
      setShowConfetti(true);
      const timer = setTimeout(() => setShowConfetti(false), 3000);
      return () => clearTimeout(timer);
    }
  }, [isWinner]);

  return (
    <div
      className={`
        relative rounded-lg overflow-visible transition-all duration-300
        ${isWinner
          ? "scale-110 z-20"
          : isDimmed
            ? "opacity-40 grayscale scale-95"
            : "hover:scale-102 hover:shadow-lg"
        }
      `}
    >
      {/* Winner effects - BIG GOLD BORDER */}
      {isWinner && (
        <>
          {/* Animated gold border */}
          <div className="absolute -inset-1 rounded-xl winner-border-glow z-0" />

          {/* Trophy decorations */}
          <div className="absolute -top-4 -left-2 z-30 text-3xl animate-bounce drop-shadow-[0_0_10px_rgba(212,175,55,0.8)]">
            🏆
          </div>
          <div className="absolute -top-4 -right-2 z-30 text-3xl animate-bounce drop-shadow-[0_0_10px_rgba(212,175,55,0.8)]" style={{ animationDelay: "0.2s" }}>
            🏆
          </div>
          <div className="absolute -bottom-3 left-1/2 -translate-x-1/2 z-30 text-2xl animate-pulse drop-shadow-[0_0_10px_rgba(212,175,55,0.8)]">
            🏆
          </div>

          {/* Winner banner */}
          <div className="absolute -top-6 left-1/2 -translate-x-1/2 z-30 whitespace-nowrap">
            <div className="bg-gradient-to-r from-yellow-600 via-yellow-400 to-yellow-600 text-black text-xs font-black px-4 py-1 rounded-full shadow-lg tracking-widest animate-pulse">
              ⭐ WINNER ⭐
            </div>
          </div>

          {/* Shimmer overlay */}
          <div className="absolute inset-0 z-10 pointer-events-none rounded-lg oscar-winner-shimmer" />

          {/* Confetti */}
          {showConfetti && (
            <div className="absolute inset-0 z-40 pointer-events-none overflow-visible">
              {[...Array(30)].map((_, i) => (
                <div
                  key={i}
                  className="confetti-piece"
                  style={{
                    left: `${Math.random() * 100}%`,
                    backgroundColor: ["#d4af37", "#ffd700", "#f5f0e6", "#ffec8b", "#b8960f"][Math.floor(Math.random() * 5)],
                    animationDelay: `${Math.random() * 0.5}s`,
                  }}
                />
              ))}
            </div>
          )}
        </>
      )}

      {/* Card content */}
      <div className={`
        relative rounded-lg overflow-hidden
        ${isWinner
          ? "ring-4 ring-yellow-400 shadow-[0_0_40px_rgba(255,215,0,0.6),0_0_80px_rgba(212,175,55,0.4)]"
          : ""
        }
      `}>
        {/* Card background */}
        <div className={`
          p-4
          ${isWinner
            ? "bg-gradient-to-b from-[#3d2a10] via-[#2a1a08] to-[#1a0d05]"
            : "bg-gradient-to-b from-[#1a1515] to-[#0d0a0a]"
          }
        `}>
          {/* Photo */}
          <div className={`
            relative aspect-[3/4] mb-3 rounded overflow-hidden
            ${isWinner
              ? "ring-4 ring-yellow-500 shadow-[0_0_20px_rgba(255,215,0,0.5)]"
              : "ring-1 ring-gold-500/20"
            }
          `}>
            {nominee.photo_url ? (
              <img
                src={nominee.photo_url}
                alt={nominee.name}
                className="absolute inset-0 w-full h-full object-cover"
              />
            ) : (
              <div className="w-full h-full bg-gradient-to-br from-[#2a2020] to-[#1a1515] flex items-center justify-center">
                <span className="text-5xl opacity-40">
                  {nominee.celebrity_id ? "👤" : "🎬"}
                </span>
              </div>
            )}
          </div>

          {/* Name */}
          <h3 className={`
            font-display text-sm md:text-base font-semibold text-center leading-tight mb-1
            ${isWinner
              ? "text-yellow-400 text-shadow-gold font-bold"
              : "text-hollywood-cream"
            }
          `}>
            {nominee.name}
          </h3>

          {/* Work title */}
          {nominee.work_title && nominee.work_title !== nominee.name && (
            <p className={`text-xs text-center italic truncate ${isWinner ? "text-yellow-200/70" : "text-hollywood-cream/50"}`}>
              {nominee.work_title}
            </p>
          )}
        </div>
      </div>

      <style jsx>{`
        .winner-border-glow {
          background: linear-gradient(
            90deg,
            #ffd700,
            #ffec8b,
            #d4af37,
            #ffec8b,
            #ffd700
          );
          background-size: 200% 200%;
          animation: border-shimmer 2s linear infinite;
          filter: blur(4px);
        }

        @keyframes border-shimmer {
          0% { background-position: 0% 50%; }
          50% { background-position: 100% 50%; }
          100% { background-position: 0% 50%; }
        }

        .oscar-winner-shimmer {
          background: linear-gradient(
            110deg,
            transparent 0%,
            transparent 40%,
            rgba(255, 215, 0, 0.2) 50%,
            transparent 60%,
            transparent 100%
          );
          background-size: 200% 100%;
          animation: winner-shimmer 1.5s linear infinite;
        }

        @keyframes winner-shimmer {
          0% { background-position: 200% 0; }
          100% { background-position: -200% 0; }
        }

        .confetti-piece {
          position: absolute;
          top: -20px;
          width: 10px;
          height: 10px;
          border-radius: 2px;
          animation: confetti-fall 2s ease-out forwards;
        }

        @keyframes confetti-fall {
          0% {
            transform: translateY(0) rotate(0deg) scale(1);
            opacity: 1;
          }
          100% {
            transform: translateY(200px) rotate(720deg) scale(0.5);
            opacity: 0;
          }
        }

        .text-shadow-gold {
          text-shadow: 0 0 10px rgba(255, 215, 0, 0.5), 0 0 20px rgba(212, 175, 55, 0.3);
        }
      `}</style>
    </div>
  );
}
