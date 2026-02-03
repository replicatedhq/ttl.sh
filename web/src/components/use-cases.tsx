"use client";

import { GitBranch, TestTube, Users, Boxes } from "lucide-react";

const useCases = [
  {
    icon: GitBranch,
    title: "CI/CD Pipelines",
    description: "Build once, test everywhere. Push your image in the build stage, pull it in parallel test jobs. No credentials to manage or rotate.",
    highlight: "Zero credential management",
  },
  {
    icon: TestTube,
    title: "Integration Testing",
    description: "Spin up ephemeral environments with throwaway images. Test your containers in isolation, then let them disappear.",
    highlight: "Clean slate, every time",
  },
  {
    icon: Users,
    title: "Open Source Projects",
    description: "Contributors can test their changes without needing registry access. No secrets, no onboarding friction, just push and test.",
    highlight: "Instant contributor experience",
  },
  {
    icon: Boxes,
    title: "Local Development",
    description: "Share images between your local machines, test on different architectures, or hand off to a teammate. Quick, easy, no setup.",
    highlight: "Share without friction",
  },
];

export function UseCases() {
  return (
    <section className="relative py-32 px-4 overflow-hidden">
      <div className="absolute inset-0 bg-gradient-to-b from-muted/30 via-background to-muted/30" />
      
      <div className="relative z-10 max-w-6xl mx-auto">
        <div className="text-center mb-16">
          <h2 className="text-4xl md:text-5xl font-bold mb-4">
            Built for{" "}
            <span className="bg-gradient-to-r from-accent to-accent/70 bg-clip-text text-transparent">
              real workflows
            </span>
          </h2>
          <p className="text-xl text-muted-foreground max-w-2xl mx-auto">
            From solo developers to massive CI/CD pipelines,
            <br />
            ttl.sh fits wherever you need temporary images.
          </p>
        </div>

        <div className="grid md:grid-cols-2 gap-8">
          {useCases.map((useCase, i) => (
            <div
              key={i}
              className="group relative p-8 rounded-2xl border border-border/50 bg-card/30 backdrop-blur-sm hover:bg-card/50 transition-all duration-300"
            >
              <div className="flex items-start gap-5">
                <div className="shrink-0 p-3 rounded-xl bg-accent/10 text-accent group-hover:bg-accent group-hover:text-accent-foreground transition-colors">
                  <useCase.icon className="w-6 h-6" />
                </div>
                <div>
                  <h3 className="text-xl font-semibold mb-2">{useCase.title}</h3>
                  <p className="text-muted-foreground mb-4 leading-relaxed">{useCase.description}</p>
                  <span className="inline-block px-3 py-1 rounded-full bg-accent/10 text-accent text-sm font-medium">
                    {useCase.highlight}
                  </span>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
