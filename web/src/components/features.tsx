"use client";

import { Shield, Zap, Globe, Lock, Clock, Workflow } from "lucide-react";

const features = [
  {
    icon: Shield,
    title: "Zero Authentication",
    description: "No tokens. No passwords. No API keys to rotate. Push and pull without any credentials whatsoever.",
    gradient: "from-green-500 to-emerald-500",
  },
  {
    icon: Clock,
    title: "Auto-Expiring",
    description: "Set your TTL from minutes to 24 hours. When time's up, the image vanishes. No cleanup, no bloat.",
    gradient: "from-blue-500 to-cyan-500",
  },
  {
    icon: Zap,
    title: "Blazing Fast",
    description: "Global CDN-backed infrastructure. Pull your images lightning fast from anywhere in the world.",
    gradient: "from-yellow-500 to-orange-500",
  },
  {
    icon: Globe,
    title: "OCI Compatible",
    description: "Works with Docker, Podman, Helm charts, and any OCI-compliant tooling. Standards-based, zero lock-in.",
    gradient: "from-purple-500 to-pink-500",
  },
  {
    icon: Lock,
    title: "Privacy First",
    description: "We don't track. We don't log. We don't even know what you're pushing. It's none of our business.",
    gradient: "from-red-500 to-rose-500",
  },
  {
    icon: Workflow,
    title: "CI/CD Native",
    description: "Built for pipelines. No secrets to configure, no credentials to leak. Just works in any CI system.",
    gradient: "from-indigo-500 to-violet-500",
  },
];

export function Features() {
  return (
    <section id="features" className="relative py-32 px-4">
      <div className="absolute inset-0 bg-gradient-to-b from-background via-muted/30 to-background" />
      
      <div className="relative z-10 max-w-6xl mx-auto">
        <div className="text-center mb-16">
          <h2 className="text-4xl md:text-5xl font-bold mb-4">
            Why developers{" "}
            <span className="bg-gradient-to-r from-accent to-accent/70 bg-clip-text text-transparent">
              love it
            </span>
          </h2>
          <p className="text-xl text-muted-foreground max-w-2xl mx-auto">
            We removed everything annoying about container registries.
            <br />
            What's left is pure simplicity.
          </p>
        </div>

        <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
          {features.map((feature, i) => (
            <div
              key={i}
              className="group relative p-8 rounded-2xl border border-border/50 bg-card/50 backdrop-blur-sm hover:border-accent/30 transition-all duration-300 hover:shadow-xl hover:shadow-accent/5"
            >
              <div className={`inline-flex p-3 rounded-xl bg-gradient-to-br ${feature.gradient} mb-4`}>
                <feature.icon className="w-6 h-6 text-white" />
              </div>
              <h3 className="text-xl font-semibold mb-2 group-hover:text-accent transition-colors">
                {feature.title}
              </h3>
              <p className="text-muted-foreground leading-relaxed">
                {feature.description}
              </p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
