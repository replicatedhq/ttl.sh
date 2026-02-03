import { Clock, Shield, Zap, Package } from "lucide-react";

const features = [
  {
    icon: Shield,
    title: "Anonymous",
    description:
      "No login required. Image names provide the initial secrecy for access. Add a UUID to reduce discoverability.",
  },
  {
    icon: Clock,
    title: "Ephemeral",
    description:
      "Image tags set the time limit. Default is 24 hours, max is 24 hours. Valid formats: :5m, :1600s, :4h, :1d",
  },
  {
    icon: Zap,
    title: "Fast",
    description:
      "Pulling images is really quick thanks to Cloudflare CDN. Works great even if you aren't near us-east-1.",
  },
  {
    icon: Package,
    title: "OCI Compatible",
    description:
      "Supports Helm charts and other OCI artifacts. Push with helm push my-chart.tgz oci://ttl.sh/your-name",
  },
];

export function Features() {
  return (
    <section id="how-it-works" className="relative py-24 px-4">
      <div className="max-w-6xl mx-auto">
        <div className="text-center mb-16">
          <h2 className="text-3xl md:text-4xl font-bold mb-4">
            Simple by design
          </h2>
          <p className="text-muted-foreground text-lg max-w-2xl mx-auto">
            No configuration, no authentication, no hassle. Just push and pull.
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {features.map((feature) => (
            <div
              key={feature.title}
              className="group p-6 bg-card border border-border rounded-lg hover:border-accent/50 transition-colors"
            >
              <div className="w-10 h-10 rounded-lg bg-accent/10 flex items-center justify-center mb-4 group-hover:bg-accent/20 transition-colors">
                <feature.icon className="w-5 h-5 text-accent" />
              </div>
              <h3 className="text-lg font-semibold mb-2">{feature.title}</h3>
              <p className="text-muted-foreground text-sm leading-relaxed">
                {feature.description}
              </p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
