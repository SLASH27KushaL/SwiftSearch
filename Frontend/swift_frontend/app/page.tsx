"use client";

import { useState } from "react";
import SearchBar from "@/components/SearchBar";
import SearchResults from "@/components/SearchResults";
import { fetchSearchResults, SearchResponse } from "@/lib/api";
import { Terminal } from "lucide-react";

export default function Home() {
  const [searchData, setSearchData] = useState<SearchResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [hasSearched, setHasSearched] = useState(false);
  const [currentQuery, setCurrentQuery] = useState("");

  const handleSearch = async (query: string) => {
    setLoading(true);
    setHasSearched(true);
    setCurrentQuery(query);
    try {
      const data = await fetchSearchResults(query);
      setSearchData(data);
    } catch (err) {
      console.error("Search request failed:", err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <main className="min-h-screen flex flex-col justify-between p-6 md:p-12">
      {/* Animated container shifting from center to top */}
      <div className={`w-full max-w-2xl mx-auto transition-all duration-700 ease-in-out ${hasSearched ? "pt-2" : "pt-[25vh] md:pt-[30vh]"}`}>
        
        {!hasSearched && (
          <div className="text-center mb-8 space-y-3">
            <div className="inline-flex items-center justify-center p-3 rounded-2xl bg-zinc-900 border border-zinc-800 shadow-inner">
              <Terminal className="text-zinc-200" size={28} />
            </div>
            <h1 className="text-3xl md:text-4xl font-bold tracking-tight text-zinc-100">
              Moogle Search
            </h1>
            <p className="text-sm text-zinc-500 max-w-xs mx-auto">
              Minimalist distributed retrieval engine powered by Go.
            </p>
          </div>
        )}

        {hasSearched && (
          <div className="flex items-center justify-between mb-6 pb-4 border-b border-zinc-900">
            <div className="flex items-center gap-2 cursor-pointer" onClick={() => setHasSearched(false)}>
              <Terminal className="text-zinc-200" size={18} />
              <span className="font-bold text-zinc-100 text-sm tracking-tight">Moogle</span>
            </div>
          </div>
        )}

        <div className="flex justify-center">
          <SearchBar 
            onSearch={handleSearch} 
            loading={loading} 
            initialQuery={currentQuery} 
            isCompact={hasSearched} 
          />
        </div>
      </div>

      <div className="flex-grow flex flex-col items-center">
        <SearchResults data={searchData} hasSearched={hasSearched} />
      </div>

      <footer className="text-center text-xs text-zinc-600 mt-16">
        SwiftSearch Engine &bull; Next.js & Go Microservices Architecture
      </footer>
    </main>
  );
}