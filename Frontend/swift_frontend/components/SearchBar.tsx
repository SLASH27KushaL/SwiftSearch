"use client";

import { useState } from "react";
import { Search, Loader2, ArrowRight } from "lucide-react";

interface SearchBarProps {
  onSearch: (query: string) => void;
  loading: boolean;
  initialQuery?: string;
  isCompact?: boolean;
}

export default function SearchBar({ onSearch, loading, initialQuery = "", isCompact = false }: SearchBarProps) {
  const [query, setQuery] = useState(initialQuery);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (query.trim()) {
      onSearch(query.trim());
    }
  };

  return (
    <form onSubmit={handleSubmit} className={`w-full transition-all duration-500 ${isCompact ? "max-w-2xl" : "max-w-xl"}`}>
      <div className={`relative flex items-center bg-zinc-900 border border-zinc-800 focus-within:border-zinc-600 transition-all shadow-xl ${isCompact ? "rounded-xl py-2 px-4" : "rounded-2xl py-3 px-5"}`}>
        <Search className="text-zinc-400 mr-3 shrink-0" size={isCompact ? 18 : 20} />
        <input
          type="text"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Search documentation, topics, or web index..."
          className="w-full bg-transparent text-zinc-100 placeholder-zinc-500 focus:outline-none text-sm md:text-base font-medium"
        />
        <button
          type="submit"
          disabled={loading || !query.trim()}
          className="ml-2 bg-zinc-100 hover:bg-white text-zinc-950 px-4 py-2 rounded-lg font-medium text-sm transition-all disabled:opacity-40 flex items-center justify-center shrink-0 cursor-pointer"
        >
          {loading ? <Loader2 className="animate-spin text-zinc-950" size={16} /> : <ArrowRight size={16} />}
        </button>
      </div>
    </form>
  );
}