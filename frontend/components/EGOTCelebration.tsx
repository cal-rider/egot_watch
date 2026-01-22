"use client";

import { useEffect, useState } from "react";

export default function EGOTCelebration() {
  const [confetti, setConfetti] = useState<Array<{ id: number; left: number; delay: number; duration: number; color: string }>>([]);
  const [showBanner, setShowBanner] = useState(false);

  useEffect(() => {
    // Generate confetti pieces
    const pieces = Array.from({ length: 100 }, (_, i) => ({
      id: i,
      left: Math.random() * 100,
      delay: Math.random() * 3,
      duration: 3 + Math.random() * 2,
      color: ["#d4af37", "#ffd700", "#f5f5dc", "#ffffff", "#c0c0c0"][Math.floor(Math.random() * 5)],
    }));
    setConfetti(pieces);

    // Show banner after a brief delay
    setTimeout(() => setShowBanner(true), 500);
  }, []);

  return (
    <>
      {/* Searchlights */}
      <div className="fixed inset-0 pointer-events-none overflow-hidden z-40">
        <div className="searchlight searchlight-left" />
        <div className="searchlight searchlight-right" />
      </div>

      {/* Confetti */}
      <div className="fixed inset-0 pointer-events-none overflow-hidden z-50">
        {confetti.map((piece) => (
          <div
            key={piece.id}
            className="confetti"
            style={{
              left: `${piece.left}%`,
              animationDelay: `${piece.delay}s`,
              animationDuration: `${piece.duration}s`,
              backgroundColor: piece.color,
            }}
          />
        ))}
      </div>

      {/* Flying Celebration Banner */}
      {showBanner && (
        <div className="fixed z-50 flying-banner pointer-events-none">
          <div className="bg-gradient-to-r from-gold-600 via-gold-500 to-gold-600 px-6 py-3 rounded-lg shadow-2xl border-2 border-gold-400">
            <div className="text-center whitespace-nowrap">
              <span className="text-xl mr-2">&#127942;</span>
              <span className="font-display text-xl text-hollywood-black font-bold tracking-wider">
                EGOT WINNER!
              </span>
              <span className="text-xl ml-2">&#127942;</span>
            </div>
          </div>
        </div>
      )}

      {/* Golden glow overlay */}
      <div className="fixed inset-0 pointer-events-none z-30 bg-gradient-radial from-gold-500/10 via-transparent to-transparent" />

      {/* Sparkles */}
      <div className="fixed inset-0 pointer-events-none z-45">
        {Array.from({ length: 20 }).map((_, i) => (
          <div
            key={i}
            className="sparkle"
            style={{
              left: `${Math.random() * 100}%`,
              top: `${Math.random() * 100}%`,
              animationDelay: `${Math.random() * 2}s`,
            }}
          />
        ))}
      </div>

      <style jsx>{`
        .searchlight {
          position: absolute;
          width: 200px;
          height: 800px;
          background: linear-gradient(
            to top,
            transparent,
            rgba(212, 175, 55, 0.1) 20%,
            rgba(212, 175, 55, 0.2) 50%,
            rgba(212, 175, 55, 0.1) 80%,
            transparent
          );
          transform-origin: bottom center;
        }

        .searchlight-left {
          bottom: 0;
          left: 10%;
          animation: searchlight-sweep-left 4s ease-in-out infinite;
        }

        .searchlight-right {
          bottom: 0;
          right: 10%;
          animation: searchlight-sweep-right 4s ease-in-out infinite;
        }

        @keyframes searchlight-sweep-left {
          0%, 100% { transform: rotate(-30deg); }
          50% { transform: rotate(20deg); }
        }

        @keyframes searchlight-sweep-right {
          0%, 100% { transform: rotate(30deg); }
          50% { transform: rotate(-20deg); }
        }

        .confetti {
          position: absolute;
          top: -10px;
          width: 10px;
          height: 10px;
          animation: confetti-fall linear forwards;
        }

        @keyframes confetti-fall {
          0% {
            transform: translateY(0) rotate(0deg);
            opacity: 1;
          }
          100% {
            transform: translateY(100vh) rotate(720deg);
            opacity: 0;
          }
        }

        .sparkle {
          position: absolute;
          width: 4px;
          height: 4px;
          background: #ffd700;
          border-radius: 50%;
          animation: sparkle 1.5s ease-in-out infinite;
        }

        @keyframes sparkle {
          0%, 100% { opacity: 0; transform: scale(0); }
          50% { opacity: 1; transform: scale(1); }
        }

        .flying-banner {
          animation: fly-around 8s ease-in-out infinite;
        }

        @keyframes fly-around {
          0% {
            top: 10%;
            left: -20%;
            transform: rotate(-5deg);
          }
          25% {
            top: 60%;
            left: 80%;
            transform: rotate(5deg);
          }
          50% {
            top: 20%;
            left: 70%;
            transform: rotate(-3deg);
          }
          75% {
            top: 70%;
            left: 10%;
            transform: rotate(4deg);
          }
          100% {
            top: 10%;
            left: -20%;
            transform: rotate(-5deg);
          }
        }
      `}</style>
    </>
  );
}
