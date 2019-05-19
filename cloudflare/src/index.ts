import { tryParsePutTagRequest, handleTagManifestRequest } from "./tag_manifest";

if (typeof addEventListener === 'function') {
  addEventListener('fetch', (e: Event): void => {
    // work around as strict typescript check doesn't allow e to be of type FetchEvent
    const fe = e as FetchEvent
    fe.respondWith(proxyRequest(fe.request))
  });
}

async function proxyRequest(r: Request): Promise<Response> {
  const url = new URL(r.url);

  if (!isRegistryRequest(r)) {
    return fetch(`https://friendly-goldstine-7897fc.netlify.com/${url.pathname}`);
  }    

  if (isCatalogRequest(r)) {
    return new Response(null, {
      status: 401,
    });
  }

  const tagManifestParams = await tryParsePutTagRequest(r.method, url);
  if (tagManifestParams) {
    return handleTagManifestRequest(r, tagManifestParams);
  }

  return fetch(r);
}

function isRegistryRequest(r: Request): boolean {
  const url = new URL(r.url);
  return url.pathname.startsWith(`/v2`);
}

function isCatalogRequest(r: Request): boolean {
  const url = new URL(r.url);
  return url.pathname.endsWith(`v2/_catalog`);
}

interface FetchEvent extends Event {
  request: Request;

  respondWith(r: Promise<Response> | Response): Promise<Response>;
}