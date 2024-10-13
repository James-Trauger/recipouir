const esbuild = require("esbuild");
const { sassPlugin } = require("esbuild-sass-plugin");

esbuild
  .build({
    entryPoints: ["react/Application.tsx", "react/style.scss"],
    outdir: "assets",
    bundle: true,
    minify: true,
    plugins: [sassPlugin()],
  })
  .then(() => console.log("⚡ Build complete! ⚡"))
  .catch(() => process.exit(1));