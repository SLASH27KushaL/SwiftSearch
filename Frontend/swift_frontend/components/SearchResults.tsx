import { SearchResponse } from "@/lib/api";
import { ExternalLink, Sparkles } from "lucide-react";

interface SearchResultsProps {
  data: SearchResponse | null;
  hasSearched: boolean;
}

export default function SearchResults({ data, hasSearched }: SearchResultsProps) {
  if (!hasSearched) return null;

  if (!data || data.results.length === 0) {
    return (
      <div className="w-full max-w-2xl mt-10 p-8 text-center rounded-2xl bg-zinc-900/40 border border-zinc-800/80">
        <p className="text-zinc-300 font-medium">No results found.</p>
        <p className="text-zinc-500 text-sm mt-1">Make sure your Go indexing pipeline has processed data for this query.</p>
      </div>
    );
  }

  return (
    <div className="w-full max-w-2xl mt-8 space-y-6">
      {/* Metrics Bar */}
      <div className="flex items-center justify-between text-xs text-zinc-500 px-1 pb-2 border-b border-zinc-900">
        <span>Found {data.total_results} results</span>
        <span>Resolved in {data.execution_time_ms}ms via Go & Redis</span>
      </div>

      {/* Result Cards Feed */}
      <div className="space-y-4">
        {data.results.map((result, idx) => (
          <article
            key={idx}
            className="group p-5 rounded-2xl bg-zinc-900/30 hover:bg-zinc-900/70 border border-zinc-800/60 hover:border-zinc-700 transition-all duration-200"
          >
            <div className="flex items-center justify-between gap-2 mb-1.5">
              <span className="text-xs font-mono text-zinc-500 truncate max-w-[80%]">
                {result.url}
              </span>
              <div className="flex items-center gap-1 text-[11px] font-mono text-cyan-400 bg-cyan-950/40 px-2 py-0.5 rounded-full border border-cyan-800/30">
                <Sparkles size={11} />
                <span>Score: {result.score.toFixed(4)}</span>
              </div>
            </div>

            <a
              href={result.url}
              target="_blank"
              rel="noopener noreferrer"
              className="text-lg font-semibold text-zinc-100 group-hover:text-cyan-400 transition-colors flex items-center gap-2"
            >
              <span>{result.title || "Untitled Document"}</span>
              <ExternalLink size={14} className="opacity-0 group-hover:opacity-100 transition-opacity text-cyan-400 shrink-0" />
            </a>
          </article>
        ))}
      </div>
    </div>
  );
}