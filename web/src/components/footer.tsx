"use client";

import { Heart } from "lucide-react";

export function Footer() {
  return (
    <footer className="relative py-16 px-4 border-t border-border/50">
      <div className="max-w-6xl mx-auto">
        <div className="flex flex-col md:flex-row items-center justify-between gap-8">
          <div className="text-center md:text-left">
            <div className="flex items-center justify-center md:justify-start gap-2 mb-2">
              <span className="text-2xl font-bold font-mono">ttl.sh</span>
            </div>
            <p className="text-muted-foreground text-sm">
              The ephemeral container registry.
            </p>
          </div>

          <div className="flex flex-wrap items-center justify-center gap-6 text-sm">
            <a
              href="https://github.com/replicatedhq/ttl.sh"
              target="_blank"
              rel="noopener noreferrer"
              className="text-muted-foreground hover:text-foreground transition-colors"
            >
              GitHub
            </a>
            <a
              href="https://github.com/replicatedhq/ttl.sh/issues"
              target="_blank"
              rel="noopener noreferrer"
              className="text-muted-foreground hover:text-foreground transition-colors"
            >
              Issues
            </a>
            <a
              href="https://www.replicated.com"
              target="_blank"
              rel="noopener noreferrer"
              className="text-muted-foreground hover:text-foreground transition-colors"
            >
              Replicated
            </a>
          </div>
        </div>

        <div className="mt-12 pt-8 border-t border-border/30 text-center">
          <p className="text-sm text-muted-foreground flex items-center justify-center gap-1">
            Made with <Heart className="w-4 h-4 text-red-500 fill-red-500" /> by{" "}
            <a
              href="https://www.replicated.com"
              target="_blank"
              rel="noopener noreferrer"
              className="text-foreground hover:text-accent transition-colors font-medium"
            >
              Replicated
            </a>
          </p>
          <p className="text-xs text-muted-foreground/60 mt-2">
            © {new Date().getFullYear()} Replicated, Inc. All rights reserved.
          </p>
        </div>
      </div>
    </footer>
  );
}
