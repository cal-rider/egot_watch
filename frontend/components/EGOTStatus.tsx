"use client";

import { EGOTStatus as EGOTStatusType } from "@/lib/api";

interface Props {
  status: EGOTStatusType;
  size?: "sm" | "md" | "lg";
}

const awards = [
  { key: "emmy", label: "E", full: "Emmy" },
  { key: "grammy", label: "G", full: "Grammy" },
  { key: "oscar", label: "O", full: "Oscar" },
  { key: "tony", label: "T", full: "Tony" },
] as const;

export default function EGOTStatus({ status, size = "md" }: Props) {
  const sizeClasses = {
    sm: "text-2xl gap-1",
    md: "text-4xl gap-2",
    lg: "text-6xl gap-3",
  };

  const starSizes = {
    sm: "w-8 h-8",
    md: "w-12 h-12",
    lg: "w-16 h-16",
  };

  return (
    <div className="flex flex-col items-center gap-2">
      <div className={`flex items-center ${sizeClasses[size]}`}>
        {awards.map(({ key, label, full }) => {
          const won = status[key as keyof typeof status];
          return (
            <div
              key={key}
              className="flex flex-col items-center group relative"
              title={full}
            >
              <div
                className={`
                  ${starSizes[size]}
                  flex items-center justify-center
                  rounded-full
                  font-display font-bold
                  transition-all duration-300
                  ${
                    won
                      ? "bg-gold-500 text-hollywood-black shadow-lg shadow-gold-500/30"
                      : "bg-hollywood-dark text-gray-600 border border-gray-700"
                  }
                `}
              >
                {label}
              </div>
              <span
                className={`
                  text-xs mt-1 font-body
                  ${won ? "text-gold-500" : "text-gray-600"}
                `}
              >
                {full}
              </span>
            </div>
          );
        })}
      </div>
      {status.isEGOT && (
        <div className="gold-shimmer text-lg font-display font-bold tracking-widest mt-2">
          ★ EGOT WINNER ★
        </div>
      )}
    </div>
  );
}
