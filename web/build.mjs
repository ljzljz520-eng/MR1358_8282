import { mkdir, readFile, writeFile } from "node:fs/promises";
import { join } from "node:path";

const root = new URL(".", import.meta.url);
const source = await readFile(new URL("public/index.html", root), "utf8");
const dist = new URL("dist/", root);
await mkdir(dist, { recursive: true });
const generated = source.replace("__BUILD_LABEL__", "NightGuide route desk");
await writeFile(new URL("index.html", dist), generated);
await writeFile(new URL("build-info.json", dist), JSON.stringify({ node: process.version, entry: "nightguide" }, null, 2));
console.log(`built ${join("dist", "index.html")}`);
