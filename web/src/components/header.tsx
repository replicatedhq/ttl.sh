"use client";

import { useState, useEffect } from "react";
import { Menu, X, Github } from "lucide-react";
import { ThemeToggle } from "@/components/theme-toggle";

export function Header() {
  const [isScrolled, setIsScrolled] = useState(false);
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

  useEffect(() => {
    const handleScroll = () => {
      setIsScrolled(window.scrollY > 10);
    };
    window.addEventListener("scroll", handleScroll);
    return () => window.removeEventListener("scroll", handleScroll);
  }, []);

  return (
    <header
      className={`fixed top-0 left-0 right-0 z-50 transition-all duration-300 ${
        isScrolled
          ? "bg-background/80 backdrop-blur-xl border-b border-border/50 shadow-sm"
          : "bg-transparent"
      }`}
    >
      <div className="max-w-6xl mx-auto px-4 py-4 flex items-center justify-between">
        <a href="/" className="flex items-center gap-2 group">
          <span className="text-2xl font-bold font-mono group-hover:text-accent transition-colors">ttl.sh</span>
        </a>

        {/* Desktop nav */}
        <nav className="hidden md:flex items-center gap-1">
          <a
            href="#features"
            className="px-4 py-2 text-sm text-muted-foreground hover:text-foreground transition-colors rounded-lg hover:bg-muted/50"
          >
            Features
          </a>
          <a
            href="#how-it-works"
            className="px-4 py-2 text-sm text-muted-foreground hover:text-foreground transition-colors rounded-lg hover:bg-muted/50"
          >
            How it works
          </a>
          <a
            href="https://github.com/replicatedhq/ttl.sh"
            target="_blank"
            rel="noopener noreferrer"
            className="px-4 py-2 text-sm text-muted-foreground hover:text-foreground transition-colors rounded-lg hover:bg-muted/50 flex items-center gap-2"
          >
            <Github className="w-4 h-4" />
            GitHub
          </a>
          <div className="ml-2 pl-2 border-l border-border">
            <ThemeToggle />
          </div>
        </nav>

        {/* Mobile menu button */}
        <div className="flex items-center gap-2 md:hidden">
          <ThemeToggle />
          <button
            className="p-2 text-muted-foreground hover:text-foreground rounded-lg hover:bg-muted/50"
            onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
            aria-label="Toggle menu"
          >
            {isMobileMenuOpen ? <X className="w-5 h-5" /> : <Menu className="w-5 h-5" />}
          </button>
        </div>
      </div>

      {/* Mobile menu */}
      {isMobileMenuOpen && (
        <div className="md:hidden bg-background/95 backdrop-blur-xl border-b border-border/50">
          <nav className="flex flex-col px-4 py-4 gap-1">
            <a
              href="#features"
              className="px-4 py-3 text-sm text-muted-foreground hover:text-foreground transition-colors rounded-lg hover:bg-muted/50"
              onClick={() => setIsMobileMenuOpen(false)}
            >
              Features
            </a>
            <a
              href="#how-it-works"
              className="px-4 py-3 text-sm text-muted-foreground hover:text-foreground transition-colors rounded-lg hover:bg-muted/50"
              onClick={() => setIsMobileMenuOpen(false)}
            >
              How it works
            </a>
            <a
              href="https://github.com/replicatedhq/ttl.sh"
              target="_blank"
              rel="noopener noreferrer"
              className="px-4 py-3 text-sm text-muted-foreground hover:text-foreground transition-colors rounded-lg hover:bg-muted/50 flex items-center gap-2"
            >
              <Github className="w-4 h-4" />
              GitHub
            </a>
          </nav>
        </div>
      )}
    </header>
  );
}
