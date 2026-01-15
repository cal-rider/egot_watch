"use client";

import { useState, useEffect, useRef, FormEvent } from "react";
import { useRouter } from "next/navigation";
import { autocompleteCelebrity, AutocompleteSuggestion } from "@/lib/api";

export default function SearchBar() {
  const [query, setQuery] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [suggestions, setSuggestions] = useState<AutocompleteSuggestion[]>([]);
  const [showSuggestions, setShowSuggestions] = useState(false);
  const [selectedIndex, setSelectedIndex] = useState(-1);
  const router = useRouter();
  const inputRef = useRef<HTMLInputElement>(null);
  const suggestionsRef = useRef<HTMLDivElement>(null);

  // Debounced autocomplete fetch
  useEffect(() => {
    const timer = setTimeout(async () => {
      if (query.trim().length >= 2) {
        const results = await autocompleteCelebrity(query);
        setSuggestions(results);
        setShowSuggestions(results.length > 0);
        setSelectedIndex(-1);
      } else {
        setSuggestions([]);
        setShowSuggestions(false);
      }
    }, 200);

    return () => clearTimeout(timer);
  }, [query]);

  // Close suggestions when clicking outside
  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (
        suggestionsRef.current &&
        !suggestionsRef.current.contains(e.target as Node) &&
        inputRef.current &&
        !inputRef.current.contains(e.target as Node)
      ) {
        setShowSuggestions(false);
      }
    };

    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    if (!query.trim()) return;

    setIsLoading(true);
    setShowSuggestions(false);
    router.push(`/celebrity/${encodeURIComponent(query.trim())}`);
  };

  const handleSelectSuggestion = (suggestion: AutocompleteSuggestion) => {
    setQuery(suggestion.name);
    setShowSuggestions(false);
    setIsLoading(true);
    router.push(`/celebrity/${encodeURIComponent(suggestion.name)}`);
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (!showSuggestions || suggestions.length === 0) return;

    if (e.key === "ArrowDown") {
      e.preventDefault();
      setSelectedIndex((prev) =>
        prev < suggestions.length - 1 ? prev + 1 : prev
      );
    } else if (e.key === "ArrowUp") {
      e.preventDefault();
      setSelectedIndex((prev) => (prev > 0 ? prev - 1 : -1));
    } else if (e.key === "Enter" && selectedIndex >= 0) {
      e.preventDefault();
      handleSelectSuggestion(suggestions[selectedIndex]);
    } else if (e.key === "Escape") {
      setShowSuggestions(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="w-full max-w-2xl mx-auto">
      <div className="relative">
        <input
          ref={inputRef}
          type="text"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          onKeyDown={handleKeyDown}
          onFocus={() => suggestions.length > 0 && setShowSuggestions(true)}
          placeholder="Search for a celebrity..."
          className="
            w-full px-6 py-4
            text-lg text-hollywood-cream
            search-input rounded-lg
            font-body placeholder-gray-500
          "
          disabled={isLoading}
          autoComplete="off"
        />
        <button
          type="submit"
          disabled={isLoading || !query.trim()}
          className="
            absolute right-2 top-1/2 -translate-y-1/2
            px-6 py-2
            bg-gold-500 text-hollywood-black
            font-display font-semibold
            rounded-md
            hover:bg-gold-400
            disabled:opacity-50 disabled:cursor-not-allowed
            transition-all duration-300
            gold-glow
          "
        >
          {isLoading ? "..." : "Search"}
        </button>

        {/* Autocomplete dropdown */}
        {showSuggestions && suggestions.length > 0 && (
          <div
            ref={suggestionsRef}
            className="
              absolute z-50 w-full mt-2
              bg-hollywood-black/95 border border-gold-500/30
              rounded-lg shadow-2xl
              backdrop-blur-sm
              overflow-hidden
            "
          >
            {suggestions.map((suggestion, index) => (
              <button
                key={suggestion.id}
                type="button"
                onClick={() => handleSelectSuggestion(suggestion)}
                className={`
                  w-full px-6 py-3 text-left
                  flex items-center gap-4
                  transition-all duration-200
                  ${
                    index === selectedIndex
                      ? "bg-gold-500/20 text-gold-400"
                      : "text-hollywood-cream hover:bg-gold-500/10"
                  }
                  ${index !== suggestions.length - 1 ? "border-b border-gold-500/10" : ""}
                `}
              >
                {suggestion.photo_url ? (
                  <img
                    src={suggestion.photo_url}
                    alt=""
                    className="w-10 h-10 rounded-full object-cover border border-gold-500/30"
                  />
                ) : (
                  <div className="w-10 h-10 rounded-full bg-gold-500/20 flex items-center justify-center">
                    <span className="text-gold-500 text-lg">★</span>
                  </div>
                )}
                <span className="font-body text-lg">{suggestion.name}</span>
              </button>
            ))}
          </div>
        )}
      </div>
    </form>
  );
}
