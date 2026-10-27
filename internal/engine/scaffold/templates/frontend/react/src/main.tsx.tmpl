import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App.tsx";
import "./index.css";
[[ if hasAuth "clerk" ]]
import { ClerkProvider } from "@clerk/clerk-react";
[[ end ]]
[[ if isAPI "trpc" ]]
import { TRPCReactProvider } from "./components/trpc-provider";
[[ else if isAPI "orpc" ]]
import { ORPCReactProvider } from "./components/orpc-provider";
[[ end ]]

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
[[ if hasAuth "clerk" ]]
    <ClerkProvider publishableKey={import.meta.env.VITE_CLERK_PUBLISHABLE_KEY || ""}>
[[ end ]]
[[ if isAPI "trpc" ]]
      <TRPCReactProvider>
        <App />
      </TRPCReactProvider>
[[ else if isAPI "orpc" ]]
      <ORPCReactProvider>
        <App />
      </ORPCReactProvider>
[[ else ]]
      <App />
[[ end ]]
[[ if hasAuth "clerk" ]]
    </ClerkProvider>
[[ end ]]
  </React.StrictMode>
);
