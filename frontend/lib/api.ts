const API_BASE = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export interface Award {
  id: string;
  celebrity_id: string;
  type: "Emmy" | "Grammy" | "Oscar" | "Tony";
  year: number;
  work: string;
  category: string;
  is_winner: boolean;
  ceremony_date?: string;
  is_upcoming?: boolean;
}

export interface Celebrity {
  id: string;
  name: string;
  slug: string;
  photo_url: string | null;
  summary: string | null;
  last_updated: string;
  awards: Award[];
}

export interface EGOTStatus {
  emmy: boolean;
  grammy: boolean;
  oscar: boolean;
  tony: boolean;
  isEGOT: boolean;
  count: number;
}

export function getEGOTStatus(awards: Award[]): EGOTStatus {
  const status = {
    emmy: awards.some((a) => a.type === "Emmy" && a.is_winner),
    grammy: awards.some((a) => a.type === "Grammy" && a.is_winner),
    oscar: awards.some((a) => a.type === "Oscar" && a.is_winner),
    tony: awards.some((a) => a.type === "Tony" && a.is_winner),
    isEGOT: false,
    count: 0,
  };
  status.count = [status.emmy, status.grammy, status.oscar, status.tony].filter(Boolean).length;
  status.isEGOT = status.count === 4;
  return status;
}

export async function searchCelebrity(name: string): Promise<Celebrity> {
  // Use AbortController with 30s timeout for Wikidata fetches
  const controller = new AbortController();
  const timeoutId = setTimeout(() => controller.abort(), 30000);

  try {
    const response = await fetch(
      `${API_BASE}/api/celebrity/search?q=${encodeURIComponent(name)}`,
      { signal: controller.signal }
    );

    clearTimeout(timeoutId);

    if (!response.ok) {
      if (response.status === 404) {
        throw new Error("Celebrity not found");
      }
      throw new Error("Failed to fetch celebrity");
    }

    return response.json();
  } catch (err) {
    clearTimeout(timeoutId);
    if (err instanceof Error && err.name === "AbortError") {
      throw new Error("Request timed out - please try again");
    }
    throw err;
  }
}

export interface AutocompleteSuggestion {
  id: string;
  name: string;
  slug: string;
  photo_url: string | null;
}

export async function autocompleteCelebrity(query: string): Promise<AutocompleteSuggestion[]> {
  if (!query.trim()) return [];

  const response = await fetch(
    `${API_BASE}/api/celebrity/autocomplete?q=${encodeURIComponent(query)}`
  );

  if (!response.ok) {
    return [];
  }

  return response.json();
}

// Close to EGOT types and functions
export interface CelebrityWithProgress {
  id: string;
  name: string;
  slug: string;
  photo_url: string | null;
  last_updated: string;
  egot_win_count: number;
  won_awards: string[]; // ["Emmy", "Grammy", "Oscar"]
}

export async function getCloseToEGOT(limit?: number): Promise<CelebrityWithProgress[]> {
  const url = limit
    ? `${API_BASE}/api/celebrity/close-to-egot?limit=${limit}`
    : `${API_BASE}/api/celebrity/close-to-egot`;

  const response = await fetch(url);

  if (!response.ok) {
    throw new Error("Failed to fetch close to EGOT celebrities");
  }

  return response.json();
}

export async function getEGOTWinners(limit?: number): Promise<CelebrityWithProgress[]> {
  const url = limit
    ? `${API_BASE}/api/celebrity/egot-winners?limit=${limit}`
    : `${API_BASE}/api/celebrity/egot-winners`;

  const response = await fetch(url);

  if (!response.ok) {
    throw new Error("Failed to fetch EGOT winners");
  }

  return response.json();
}

export interface CelebrityBasic {
  id: string;
  name: string;
  slug: string;
  photo_url: string | null;
  summary: string | null;
  last_updated: string;
}

export async function getNoAwards(limit?: number): Promise<CelebrityBasic[]> {
  const url = limit
    ? `${API_BASE}/api/celebrity/no-awards?limit=${limit}`
    : `${API_BASE}/api/celebrity/no-awards`;

  const response = await fetch(url);

  if (!response.ok) {
    throw new Error("Failed to fetch celebrities with no awards");
  }

  return response.json();
}

// Oscar Race types and functions
export interface OscarNominee {
  id: string;
  category_id: string;
  celebrity_id?: string;
  name: string;
  photo_url: string | null;
  work_title: string | null;
  is_winner: boolean;
  display_order: number;
}

export interface OscarCategory {
  id: string;
  ceremony_id: string;
  name: string;
  display_order: number;
  winner_announced: boolean;
  nominees: OscarNominee[];
}

export interface OscarCeremony {
  id: string;
  year: number;
  ceremony_name: string | null;
  ceremony_date: string | null;
  is_complete: boolean;
  created_at: string;
  categories: OscarCategory[];
}

export async function getOscarCeremony(year: number): Promise<OscarCeremony> {
  const response = await fetch(`${API_BASE}/api/oscar-race/${year}`);

  if (!response.ok) {
    if (response.status === 404) {
      throw new Error("Ceremony not found");
    }
    throw new Error("Failed to fetch Oscar ceremony");
  }

  return response.json();
}

export async function getOscarYears(): Promise<number[]> {
  const response = await fetch(`${API_BASE}/api/oscar-race/years`);

  if (!response.ok) {
    throw new Error("Failed to fetch Oscar years");
  }

  return response.json();
}

export async function setOscarWinner(
  year: number,
  categoryId: string,
  nomineeId: string
): Promise<void> {
  const response = await fetch(
    `${API_BASE}/api/oscar-race/${year}/category/${categoryId}/winner/${nomineeId}`,
    { method: "PUT" }
  );

  if (!response.ok) {
    throw new Error("Failed to set winner");
  }
}
