import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App.tsx";
import "./index.css";
[[ if eq .Auth "clerk" ]]
import { ClerkProvider } from "@clerk/clerk-react";
[[ end ]]
[[ if eq .API "trpc" ]]
import { TRPCReactProvider } from "./components/trpc-provider";
[[ else if eq .API "orpc" ]]
import { ORPCReactProvider } from "./components/orpc-provider";
[[ end ]]

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
[[ if eq .Auth "clerk" ]]
    <ClerkProvider publishableKey={import.meta.env.VITE_CLERK_PUBLISHABLE_KEY || ""}>
[[ end ]]
[[ if eq .API "trpc" ]]
      <TRPCReactProvider>
        <App />
      </TRPCReactProvider>
[[ else if eq .API "orpc" ]]
      <ORPCReactProvider>
        <App />
      </ORPCReactProvider>
[[ else ]]
      <App />
[[ end ]]
[[ if eq .Auth "clerk" ]]
    </ClerkProvider>
[[ end ]]
  </React.StrictMode>
);
