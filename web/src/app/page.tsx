export default function Home() {
  return (
    <div className="min-h-screen bg-white flex flex-col items-center justify-center">
      <main className="text-center px-6">
        <h1 className="text-6xl font-bold text-red-500 mb-4">ttl.sh</h1>
        <p className="text-xl text-gray-600 mb-8">
          Anonymous &amp; ephemeral Docker image registry
        </p>
        <div className="bg-gray-900 text-gray-100 p-6 rounded-lg text-left font-mono text-sm max-w-xl">
          <p className="text-gray-400"># Push an image with a 1 hour TTL</p>
          <p>$ docker push ttl.sh/my-image:1h</p>
          <p className="mt-4 text-gray-400"># Pull it back (before it expires)</p>
          <p>$ docker pull ttl.sh/my-image:1h</p>
        </div>
        <p className="mt-8 text-gray-500">
          Free to use. No sign-up required. Max TTL: 24 hours.
        </p>
        <p className="mt-4 text-gray-400 text-sm">
          Contributed by{" "}
          <a
            href="https://www.replicated.com"
            className="underline hover:text-gray-600"
            target="_blank"
            rel="noopener noreferrer"
          >
            Replicated
          </a>
        </p>
      </main>
    </div>
  );
}
