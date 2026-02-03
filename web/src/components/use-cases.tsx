import { GitBranch, Users, TestTube } from "lucide-react";

const useCases = [
  {
    icon: GitBranch,
    title: "CI/CD Workflows",
    description:
      "Build images in CI, test them in parallel jobs, and discard them after. No need to manage credentials or clean up old images.",
  },
  {
    icon: Users,
    title: "Open Source Contributions",
    description:
      "Allow contributors to test their changes without requiring them to configure their own registry credentials.",
  },
  {
    icon: TestTube,
    title: "Quick Testing",
    description:
      "Share images with teammates for quick testing sessions. The images disappear automatically when you're done.",
  },
];

export function UseCases() {
  return (
    <section className="relative py-24 px-4 border-t border-border">
      <div className="max-w-6xl mx-auto">
        <div className="text-center mb-16">
          <h2 className="text-3xl md:text-4xl font-bold mb-4">
            Perfect for CI workflows
          </h2>
          <p className="text-muted-foreground text-lg max-w-2xl mx-auto">
            ttl.sh solves the problem of sharing intermediate build artifacts 
            without managing credentials or cleanup.
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          {useCases.map((useCase) => (
            <div
              key={useCase.title}
              className="p-6 bg-card border border-border rounded-lg"
            >
              <div className="w-12 h-12 rounded-lg bg-secondary flex items-center justify-center mb-4">
                <useCase.icon className="w-6 h-6 text-foreground" />
              </div>
              <h3 className="text-xl font-semibold mb-3">{useCase.title}</h3>
              <p className="text-muted-foreground leading-relaxed">
                {useCase.description}
              </p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
