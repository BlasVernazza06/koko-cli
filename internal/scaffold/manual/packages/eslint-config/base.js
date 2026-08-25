/** @type {import("eslint").Linter.Config} */
module.exports = {
  extends: ["turbo", "prettier"],
  rules: {
    "@next/next/no-html-link-for-pages": "off",
  },
  env: {
    node: true,
    es2022: true,
  },
};
