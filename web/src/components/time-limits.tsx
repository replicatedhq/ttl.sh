const timeLimits = [
  { tag: ":5m", description: "5 minutes" },
  { tag: ":30m", description: "30 minutes" },
  { tag: ":1h", description: "1 hour" },
  { tag: ":6h", description: "6 hours" },
  { tag: ":12h", description: "12 hours" },
  { tag: ":1d", description: "24 hours (max)" },
];

export function TimeLimits() {
  return (
    <section className="relative py-24 px-4 border-t border-border">
      <div className="max-w-4xl mx-auto">
        <div className="text-center mb-16">
          <h2 className="text-3xl md:text-4xl font-bold mb-4">
            Flexible time limits
          </h2>
          <p className="text-muted-foreground text-lg">
            Set your image expiry using the tag. Default and maximum is 24 hours.
          </p>
        </div>

        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-4">
          {timeLimits.map((limit) => (
            <div
              key={limit.tag}
              className="p-4 bg-card border border-border rounded-lg text-center hover:border-accent/50 transition-colors"
            >
              <code className="text-lg font-mono text-accent font-semibold">
                {limit.tag}
              </code>
              <p className="text-sm text-muted-foreground mt-1">
                {limit.description}
              </p>
            </div>
          ))}
        </div>

        <div className="mt-8 p-4 bg-secondary/50 border border-border rounded-lg">
          <p className="text-sm text-muted-foreground text-center">
            You can also use seconds format like <code className="text-accent">:1600s</code> for 
            more precise control. Any value exceeding 24 hours will be capped.
          </p>
        </div>
      </div>
    </section>
  );
}
