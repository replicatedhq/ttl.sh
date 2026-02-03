"use client";

import { useState } from "react";
import { Check, Copy } from "lucide-react";

const steps = [
  {
    number: "01",
    title: "Tag your image",
    description: "Use ttl.sh as the registry, add a unique name, and specify the TTL as the tag.",
    code: "docker build -t ttl.sh/my-image:1h .",
  },
  {
    number: "02",
    title: "Push it",
    description: "Push your image to ttl.sh. No authentication required.",
    code: "docker push ttl.sh/my-image:1h",
  },
  {
    number: "03",
    title: "Pull before it expires",
    description: "Pull your image from any machine before the TTL expires.",
    code: "docker pull ttl.sh/my-image:1h",
  },
];

export function HowTo() {
  const [copiedIndex, setCopiedIndex] = useState<number | null>(null);

  const copyToClipboard = async (text: string, index: number) => {
    await navigator.clipboard.writeText(text);
    setCopiedIndex(index);
    setTimeout(() => setCopiedIndex(null), 2000);
  };

  return (
    <section className="relative py-24 px-4 border-t border-border">
      <div className="max-w-4xl mx-auto">
        <div className="text-center mb-16">
          <h2 className="text-3xl md:text-4xl font-bold mb-4">
            Three simple steps
          </h2>
          <p className="text-muted-foreground text-lg">
            Get started in seconds. No account needed.
          </p>
        </div>

        <div className="space-y-8">
          {steps.map((step, index) => (
            <div
              key={step.number}
              className="flex flex-col md:flex-row gap-6 items-start"
            >
              <div className="flex-shrink-0 w-12 h-12 rounded-full bg-accent/10 flex items-center justify-center text-accent font-mono font-bold">
                {step.number}
              </div>
              <div className="flex-1">
                <h3 className="text-xl font-semibold mb-2">{step.title}</h3>
                <p className="text-muted-foreground mb-4">{step.description}</p>
                <div className="bg-card border border-border rounded-lg p-4 font-mono text-sm flex items-center justify-between group">
                  <code>
                    <span className="text-accent">$</span> {step.code}
                  </code>
                  <button
                    onClick={() => copyToClipboard(step.code, index)}
                    className="opacity-0 group-hover:opacity-100 transition-opacity p-1 hover:bg-secondary rounded ml-4"
                    aria-label="Copy command"
                  >
                    {copiedIndex === index ? (
                      <Check className="w-4 h-4 text-accent" />
                    ) : (
                      <Copy className="w-4 h-4 text-muted-foreground" />
                    )}
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
