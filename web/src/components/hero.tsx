"use client";

import { useState, useEffect } from "react";
import { Check, Copy, Sparkles } from "lucide-react";

export function Hero() {
  const [copiedIndex, setCopiedIndex] = useState<number | null>(null);
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
  }, []);

  const commands = [
    "docker build -t ttl.sh/my-image:1h .",
    "docker push ttl.sh/my-image:1h",
  ];

  const copyToClipboard = async (text: string, index: number) => {
    await navigator.clipboard.writeText(text);
    setCopiedIndex(index);
    setTimeout(() => setCopiedIndex(null), 2000);
  };

  return (
    <section className="relative min-h-screen flex flex-col items-center justify-center px-4 py-20 overflow-hidden">
      {/* Animated gradient background */}
      <div className="absolute inset-0 bg-gradient-to-br from-background via-background to-accent/5" />
      
      {/* Grid pattern */}
      <div className="absolute inset-0 bg-[linear-gradient(to_right,rgba(0,0,0,0.03)_1px,transparent_1px),linear-gradient(to_bottom,rgba(0,0,0,0.03)_1px,transparent_1px)] dark:bg-[linear-gradient(to_right,rgba(255,255,255,0.02)_1px,transparent_1px),linear-gradient(to_bottom,rgba(255,255,255,0.02)_1px,transparent_1px)] bg-[size:3rem_3rem]" />
      
      {/* Glow effects */}
      <div className="absolute top-1/4 left-1/4 w-96 h-96 bg-accent/20 rounded-full blur-[128px] animate-pulse" />
      <div className="absolute bottom-1/4 right-1/4 w-96 h-96 bg-accent/10 rounded-full blur-[128px] animate-pulse delay-1000" />
      
      <div className={`relative z-10 max-w-5xl mx-auto text-center transition-all duration-1000 ${mounted ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-4'}`}>
        {/* Badge */}
        <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full border border-accent/30 bg-accent/5 text-sm text-accent mb-8 backdrop-blur-sm">
          <Sparkles className="w-4 h-4" />
          <span className="font-medium">Free forever. Zero sign-up.</span>
        </div>
        
        {/* Main heading with gradient */}
        <h1 className="text-5xl sm:text-6xl md:text-8xl font-bold tracking-tight mb-6">
          <span className="bg-gradient-to-r from-foreground via-foreground to-foreground/70 bg-clip-text text-transparent">
            Push. Pull.
          </span>
          <br />
          <span className="bg-gradient-to-r from-accent via-accent to-accent/70 bg-clip-text text-transparent">
            Forget.
          </span>
        </h1>
        
        <p className="text-xl md:text-2xl text-muted-foreground max-w-2xl mx-auto mb-4 text-pretty font-light">
          The ephemeral container registry that just works.
        </p>
        
        <p className="text-lg text-muted-foreground/80 max-w-xl mx-auto mb-12 text-pretty">
          Push any OCI artifact with a TTL. No auth, no config, no account.
          <br className="hidden sm:block" />
          It expires automatically. Like magic.
        </p>

        {/* Terminal */}
        <div className="bg-card/80 backdrop-blur-xl border border-border/50 rounded-xl overflow-hidden text-left max-w-2xl mx-auto shadow-2xl shadow-accent/5">
          <div className="flex items-center gap-2 px-4 py-3 border-b border-border/50 bg-muted/30">
            <div className="w-3 h-3 rounded-full bg-red-500/80" />
            <div className="w-3 h-3 rounded-full bg-yellow-500/80" />
            <div className="w-3 h-3 rounded-full bg-green-500/80" />
            <span className="ml-2 text-xs text-muted-foreground font-mono">~ two commands. that&apos;s it.</span>
          </div>
          <div className="p-6 font-mono text-sm space-y-3">
            {commands.map((cmd, i) => (
              <div key={i} className="flex items-center gap-3 group">
                <span className="text-accent select-none font-bold">$</span>
                <code className="flex-1 text-foreground">{cmd}</code>
                <button
                  onClick={() => copyToClipboard(cmd, i)}
                  className="opacity-0 group-hover:opacity-100 transition-opacity p-1.5 hover:bg-accent/10 rounded-md"
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
            <div className="pt-3 border-t border-border/30 text-muted-foreground text-xs">
              <span className="text-accent/70">#</span> Available worldwide for 1 hour. Then it vanishes. ✨
            </div>
          </div>
        </div>

        {/* CTA buttons */}
        <div className="flex flex-wrap items-center justify-center gap-4 mt-10">
          <a
            href="https://github.com/replicatedhq/ttl.sh"
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center gap-2 px-6 py-3 bg-accent text-accent-foreground rounded-lg font-semibold hover:bg-accent/90 transition-all hover:scale-105 shadow-lg shadow-accent/20"
          >
            <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
              <path fillRule="evenodd" clipRule="evenodd" d="M12 2C6.477 2 2 6.477 2 12c0 4.42 2.865 8.17 6.839 9.49.5.092.682-.217.682-.482 0-.237-.008-.866-.013-1.7-2.782.604-3.369-1.34-3.369-1.34-.454-1.156-1.11-1.464-1.11-1.464-.908-.62.069-.608.069-.608 1.003.07 1.531 1.03 1.531 1.03.892 1.529 2.341 1.087 2.91.831.092-.646.35-1.086.636-1.336-2.22-.253-4.555-1.11-4.555-4.943 0-1.091.39-1.984 1.029-2.683-.103-.253-.446-1.27.098-2.647 0 0 .84-.269 2.75 1.025A9.578 9.578 0 0112 6.836c.85.004 1.705.114 2.504.336 1.909-1.294 2.747-1.025 2.747-1.025.546 1.377.203 2.394.1 2.647.64.699 1.028 1.592 1.028 2.683 0 3.842-2.339 4.687-4.566 4.935.359.309.678.919.678 1.852 0 1.336-.012 2.415-.012 2.743 0 .267.18.578.688.48C19.138 20.167 22 16.418 22 12c0-5.523-4.477-10-10-10z" />
            </svg>
            Star on GitHub
          </a>
          <a
            href="#features"
            className="inline-flex items-center gap-2 px-6 py-3 border border-border rounded-lg font-semibold hover:bg-muted/50 transition-all"
          >
            See how it works
          </a>
        </div>

        {/* Social proof */}
        <div className="mt-16 pt-8 border-t border-border/30">
          <p className="text-sm text-muted-foreground mb-6">Trusted by engineering teams worldwide</p>
          <div className="flex flex-wrap items-center justify-center gap-8 opacity-60">
            <span className="text-lg font-semibold text-muted-foreground">GitHub Actions</span>
            <span className="text-lg font-semibold text-muted-foreground">GitLab CI</span>
            <span className="text-lg font-semibold text-muted-foreground">CircleCI</span>
            <span className="text-lg font-semibold text-muted-foreground">Jenkins</span>
          </div>
        </div>
      </div>
    </section>
  );
}
