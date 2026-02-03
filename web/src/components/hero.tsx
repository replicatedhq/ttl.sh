"use client";

import { useState } from "react";
import { Check, Copy } from "lucide-react";

export function Hero() {
  const [copiedIndex, setCopiedIndex] = useState<number | null>(null);

  const commands = [
    "IMAGE_NAME=$(uuidgen)",
    "docker build -t ttl.sh/${IMAGE_NAME}:1h .",
    "docker push ttl.sh/${IMAGE_NAME}:1h",
  ];

  const copyToClipboard = async (text: string, index: number) => {
    await navigator.clipboard.writeText(text);
    setCopiedIndex(index);
    setTimeout(() => setCopiedIndex(null), 2000);
  };

  return (
    <section className="relative min-h-screen flex flex-col items-center justify-center px-4 py-20">
      {/* Grid background */}
      <div className="absolute inset-0 bg-[linear-gradient(to_right,rgba(0,0,0,0.05)_1px,transparent_1px),linear-gradient(to_bottom,rgba(0,0,0,0.05)_1px,transparent_1px)] dark:bg-[linear-gradient(to_right,rgba(255,255,255,0.03)_1px,transparent_1px),linear-gradient(to_bottom,rgba(255,255,255,0.03)_1px,transparent_1px)] bg-[size:4rem_4rem]" />
      
      <div className="relative z-10 max-w-4xl mx-auto text-center">
        <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full border border-border bg-card/50 text-sm text-muted-foreground mb-8">
          <span className="w-2 h-2 rounded-full bg-accent animate-pulse" />
          Free to use. No sign-up required.
        </div>
        
        <h1 className="text-4xl sm:text-5xl md:text-7xl font-bold tracking-tight mb-6 text-balance">
          Anonymous & Ephemeral
          <br />
          <span className="text-muted-foreground">OCI Registry</span>
        </h1>
        
        <p className="text-lg md:text-xl text-muted-foreground max-w-2xl mx-auto mb-12 text-pretty">
          Push your OCI images with a time limit. Pull them before they expire. 
          No authentication, no configuration, just works.
        </p>

        {/* Terminal */}
        <div className="bg-card border border-border rounded-lg overflow-hidden text-left max-w-2xl mx-auto">
          <div className="flex items-center gap-2 px-4 py-3 border-b border-border bg-secondary/30">
            <div className="w-3 h-3 rounded-full bg-red-500/80" />
            <div className="w-3 h-3 rounded-full bg-yellow-500/80" />
            <div className="w-3 h-3 rounded-full bg-green-500/80" />
            <span className="ml-2 text-xs text-muted-foreground font-mono">terminal</span>
          </div>
          <div className="p-4 font-mono text-sm space-y-2">
            {commands.map((cmd, i) => (
              <div key={i} className="flex items-center gap-3 group">
                <span className="text-accent select-none">$</span>
                <code className="flex-1 text-foreground/90">{cmd}</code>
                <button
                  onClick={() => copyToClipboard(cmd, i)}
                  className="opacity-0 group-hover:opacity-100 transition-opacity p-1 hover:bg-secondary rounded"
                  aria-label="Copy command"
                >
                  {copiedIndex === i ? (
                    <Check className="w-4 h-4 text-accent" />
                  ) : (
                    <Copy className="w-4 h-4 text-muted-foreground" />
                  )}
                </button>
              </div>
            ))}
            <div className="pt-2 text-muted-foreground text-xs">
              {"# "}Your image is now available for 1 hour at ttl.sh/${"{"}IMAGE_NAME{"}"}:1h
            </div>
          </div>
        </div>

        {/* Quick links */}
        <div className="flex flex-wrap items-center justify-center gap-4 mt-8">
          <a
            href="https://github.com/replicatedhq/ttl.sh"
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center gap-2 px-5 py-2.5 bg-primary text-primary-foreground rounded-lg font-medium hover:bg-primary/90 transition-colors"
          >
            <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
              <path fillRule="evenodd" clipRule="evenodd" d="M12 2C6.477 2 2 6.477 2 12c0 4.42 2.865 8.17 6.839 9.49.5.092.682-.217.682-.482 0-.237-.008-.866-.013-1.7-2.782.604-3.369-1.34-3.369-1.34-.454-1.156-1.11-1.464-1.11-1.464-.908-.62.069-.608.069-.608 1.003.07 1.531 1.03 1.531 1.03.892 1.529 2.341 1.087 2.91.831.092-.646.35-1.086.636-1.336-2.22-.253-4.555-1.11-4.555-4.943 0-1.091.39-1.984 1.029-2.683-.103-.253-.446-1.27.098-2.647 0 0 .84-.269 2.75 1.025A9.578 9.578 0 0112 6.836c.85.004 1.705.114 2.504.336 1.909-1.294 2.747-1.025 2.747-1.025.546 1.377.203 2.394.1 2.647.64.699 1.028 1.592 1.028 2.683 0 3.842-2.339 4.687-4.566 4.935.359.309.678.919.678 1.852 0 1.336-.012 2.415-.012 2.743 0 .267.18.578.688.48C19.138 20.167 22 16.418 22 12c0-5.523-4.477-10-10-10z" />
            </svg>
            View on GitHub
          </a>
          <a
            href="#how-it-works"
            className="inline-flex items-center gap-2 px-5 py-2.5 border border-border rounded-lg font-medium hover:bg-secondary transition-colors"
          >
            Learn more
          </a>
        </div>
      </div>
    </section>
  );
}
