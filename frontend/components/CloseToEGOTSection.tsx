"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { CelebrityWithProgress, getCloseToEGOT } from "@/lib/api";
import CloseToEGOTCard from "./CloseToEGOTCard";

export default function CloseToEGOTSection() {
  const [celebrities, setCelebrities] = useState<CelebrityWithProgress[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function fetchData() {
      try {
        const data = await getCloseToEGOT(6); // Show top 6 on home page
        setCelebrities(data);
      } catch (err) {
        setError("Failed to load data");
        console.error("Failed to fetch close to EGOT:", err);
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
          <div className="text-gold-500 animate-pulse text-2xl">★ ★ ★</div>
          <p className="text-gray-400 mt-2">Loading...</p>
        </div>
      </section>
    );
  }

  if (error || celebrities.length === 0) {
    return null; // Don't show section if no data
  }

  return (
    <section className="w-full max-w-4xl mx-auto mt-16 px-4">
      <div className="art-deco-line mb-6" />

      <div className="flex items-center justify-between mb-6">
        <h2 className="font-display text-xl md:text-2xl text-gold-500 tracking-wider">
          ★ CLOSE TO EGOT ★
        </h2>
        <Link
          href="/close-to-egot"
          className="text-sm text-gray-400 hover:text-gold-500 transition-colors font-display"
        >
          View All →
        </Link>
      </div>

      <p className="text-gray-400 text-sm mb-6">
        These legends are just one award away from achieving EGOT status
      </p>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {celebrities.map((celeb) => (
          <CloseToEGOTCard key={celeb.id} celebrity={celeb} />
        ))}
      </div>
    </section>
  );
}
