import { tryParsePutTagRequest } from "./handler";
import * as chai from "chai";

const testTryParsePutRequest = async () => {
  let result = await tryParsePutTagRequest("PUT", new URL("v2/random/more-random/manifests/2h"));
  expect(result).to.be.defined();
  expect(result.namespace).to.equal("random");
  expect(result.name).to.equal("more-random");
  expect(result.tag).to.equal("2h");

  result = await tryParsePutTagRequest("PUT", new URL("v2/image-name/manifests/2h"));
  expect(result).to.be.defined();
  expect(result.namespace).to.equal("");
  expect(result.name).to.equal("image-name");
  expect(result.tag).to.equal("2h");

  
}