"use client";

const timeLimits = [
  { tag: ":5m", duration: "5 minutes", description: "Quick tests" },
  { tag: ":30m", duration: "30 minutes", description: "CI jobs" },
  { tag: ":1h", duration: "1 hour", description: "Most common" },
  { tag: ":6h", duration: "6 hours", description: "Long pipelines" },
  { tag: ":24h", duration: "24 hours", description: "Maximum" },
];

export function TimeLimits() {
  return (
    <section className="relative py-32 px-4">
      <div className="max-w-4xl mx-auto">
        <div className="text-center mb-16">
          <h2 className="text-4xl md:text-5xl font-bold mb-4">
            Set it and{" "}
            <span className="bg-gradient-to-r from-accent to-accent/70 bg-clip-text text-transparent">
              forget it
            </span>
          </h2>
          <p className="text-xl text-muted-foreground max-w-2xl mx-auto">
            Choose your TTL. Default is 1 hour, max is 24 hours.
            <br />
            When time's up, it's gone forever.
          </p>
        </div>

        <div className="flex flex-wrap justify-center gap-4">
          {timeLimits.map((limit, i) => (
            <div
              key={i}
              className="group relative px-6 py-5 rounded-2xl border border-border/50 bg-card/50 hover:border-accent/50 hover:bg-card transition-all duration-300 text-center min-w-[140px]"
            >
              <code className="text-2xl font-bold text-accent font-mono">{limit.tag}</code>
              <p className="text-sm text-foreground font-medium mt-2">{limit.duration}</p>
              <p className="text-xs text-muted-foreground mt-1">{limit.description}</p>
            </div>
          ))}
        </div>

        <div className="mt-12 text-center">
          <p className="text-muted-foreground">
            No TTL specified? Defaults to <code className="text-accent font-mono">:1h</code>
          </p>
        </div>
      </div>
    </section>
  );
}
