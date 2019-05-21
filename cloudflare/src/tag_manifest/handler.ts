import { TagManifestParams } from "./";
import * as parseDuration from "parse-duration";
import * as moment from "moment";

export function tryParsePutTagRequest(method: string, url: URL): TagManifestParams | void {
  const split = url.pathname.split("/");

  // v2/ns/image/manifests/tag <- namespaced images
  // v2image/manifests/tag <- library images
  if (method !== "PUT") {
    return;
  }

  if (split.length === 6 && split[4] === "manifests") {
    const tagManifestParams: TagManifestParams = {
      namespace: split[2],
      name: split[3],
      tag: split[5],
    };

    return tagManifestParams;
  }

  if (split.length === 5 && split[3] === "manifests") {
    const tagManifestParams: TagManifestParams = {
      namespace: "",
      name: split[2],
      tag: split[4],
    };

    return tagManifestParams;
  }

  return;
}

export async function handleTagManifestRequest(r: Request, params: TagManifestParams): Promise<Response> {
  r.headers["Accept"] = "application/vnd.docker.distribution.manifest.v2+json";

  const upstreamResponse = await fetch(r);
  if (upstreamResponse.status != 201) {
    return upstreamResponse;
  }

  let imageName;
  if (params.namespace === "") {
    imageName = params.name;
  } else {
    imageName = `${params.namespace}/${params.name}`;
  }

  const teapotResponse = {
    "errors": [
      {
        "message": `\r              \nimage ttl.sh/${imageName} is available now and will be automatically deleted ${expirationFromTag(params.tag)}\n\ttl.sh is contributed by Replicated (www.replicated.com)`
      }
    ]
  };

  return new Response(JSON.stringify(teapotResponse), {
    status: 418, // This is what works /shrug
    headers: {
      "Content-Type": "application/json",
    },
  });
}

function expirationFromTag(tag: string): string {
  const parsed = parseDuration(tag);
  const now = new Date();
  const then = moment(now.getTime() + parsed);
  return then.fromNow();
}
