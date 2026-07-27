export interface SearchResult {
  title: string;
  url: string;
  score: number;
}

export interface SearchResponse {
  query: string;
  total_results: number;
  execution_time_ms: number;
  results: SearchResult[];
}

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export async function fetchSearchResults(query: string): Promise<SearchResponse> {
  const res = await fetch(`${API_BASE_URL}/search?q=${encodeURIComponent(query)}`, {
    cache: "no-store",
  });

  if (!res.ok) {
    // 🚨 GRAB THE REAL ERROR FROM THE BACKEND
    const errorText = await res.text();
    console.error(`🚨 BACKEND CRASHED! Status: ${res.status} | Message: ${errorText}`);

    // Throw it so Next.js shows it on your screen instead of the generic message
    throw new Error(`Backend Error ${res.status}: ${errorText}`);
  }

  return res.json();
}