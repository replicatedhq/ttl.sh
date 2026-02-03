export function Footer() {
  return (
    <footer className="border-t border-border py-12 px-4">
      <div className="max-w-6xl mx-auto">
        <div className="flex flex-col md:flex-row items-center justify-between gap-6">
          <div className="flex items-center gap-2">
            <span className="text-xl font-bold font-mono">ttl.sh</span>
            <span className="text-muted-foreground text-sm">
              — Anonymous & Ephemeral OCI Registry
            </span>
          </div>
          
          <div className="flex items-center gap-6 text-sm text-muted-foreground">
            <a
              href="https://github.com/replicatedhq/ttl.sh"
              target="_blank"
              rel="noopener noreferrer"
              className="hover:text-foreground transition-colors"
            >
              GitHub
            </a>
            <a
              href="https://www.replicated.com"
              target="_blank"
              rel="noopener noreferrer"
              className="hover:text-foreground transition-colors"
            >
              Replicated
            </a>
          </div>
        </div>
        
        <div className="mt-8 pt-8 border-t border-border text-center text-sm text-muted-foreground">
          <p>
            Contributed by{" "}
            <a
              href="https://www.replicated.com"
              target="_blank"
              rel="noopener noreferrer"
              className="text-foreground hover:text-accent transition-colors"
            >
              Replicated
            </a>
          </p>
          <p className="mt-2">Open source. Free to use.</p>
        </div>
      </div>
    </footer>
  );
}
